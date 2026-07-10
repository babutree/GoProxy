package storage

import (
	"path/filepath"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// TestRiskColumnsMigrationAndDefault 验证 ipapiis_score / ipapi_flags 两列存在且默认值正确
// （未探测：ipapiis_score = -1，ipapi_flags = ""）。
func TestRiskColumnsMigrationAndDefault(t *testing.T) {
	store := newTestStorage(t)
	assertColumnExists(t, store, "ipapiis_score")
	assertColumnExists(t, store, "ipapi_flags")

	insertProxy(t, store, "risk-default:8080", "http", "us", SourceManual, 100, "active", 0)
	p, err := store.GetProxyByAddress("risk-default:8080")
	if err != nil {
		t.Fatalf("GetProxyByAddress() error = %v", err)
	}
	// 未经风险探测的新节点，DB DEFAULT：ipapiis_score=-1、ipapi_flags=""。
	if p.IPAPIIsScore != -1 {
		t.Fatalf("default ipapiis_score = %v, want -1", p.IPAPIIsScore)
	}
	if p.IPAPIFlags != "" {
		t.Fatalf("default ipapi_flags = %q, want empty", p.IPAPIFlags)
	}
}

// TestRiskColumnsMigrationIdempotentOnLegacyDB 在旧库（无风险列）上重复 New()，
// 迁移应幂等：第一次补两列，第二次不报错、值不变。
func TestRiskColumnsMigrationIdempotentOnLegacyDB(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "proxy.db")
	seedLegacyDB(t, dbPath)

	store, err := New(dbPath)
	if err != nil {
		t.Fatalf("New() first run error = %v", err)
	}
	assertColumnExists(t, store, "ipapiis_score")
	assertColumnExists(t, store, "ipapi_flags")
	// 旧行迁移后默认值。
	var score float64
	var flags string
	if err := store.db.QueryRow(`SELECT ipapiis_score, ipapi_flags FROM proxies WHERE address = '1.1.1.1:8080'`).Scan(&score, &flags); err != nil {
		t.Fatalf("read migrated risk columns: %v", err)
	}
	if score != -1 || flags != "" {
		t.Fatalf("migrated legacy risk = (%v,%q), want (-1,\"\")", score, flags)
	}
	store.Close()

	// 第二次打开同一库：addProxyColumnIfMissing 应识别列已存在，跳过 ALTER，不报错。
	store2, err := New(dbPath)
	if err != nil {
		t.Fatalf("New() second run error = %v", err)
	}
	defer store2.Close()
	assertColumnExists(t, store2, "ipapiis_score")
	assertColumnExists(t, store2, "ipapi_flags")

	// 确认两列各只有一列（未被重复添加）。
	for _, col := range []string{"ipapiis_score", "ipapi_flags"} {
		var count int
		if err := store2.db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('proxies') WHERE name = ?`, col).Scan(&count); err != nil {
			t.Fatalf("count %s column: %v", col, err)
		}
		if count != 1 {
			t.Fatalf("%s column count = %d, want 1", col, count)
		}
	}
}

// TestUpdateProxyExitInfoWritesRiskSignals 正向：成功路径写入两源风险信号，scanProxy 正确读回。
func TestUpdateProxyExitInfoWritesRiskSignals(t *testing.T) {
	store := newTestStorage(t)
	insertProxy(t, store, "risk-write:8080", "http", "us", SourceManual, 100, "active", 0)
	p, err := store.GetProxyByAddress("risk-write:8080")
	if err != nil {
		t.Fatalf("GetProxyByAddress() error = %v", err)
	}

	if err := store.UpdateProxyExitInfo(p.ID, "203.0.113.20", "US Ashburn", 42, 0.65, "proxy,hosting"); err != nil {
		t.Fatalf("UpdateProxyExitInfo() error = %v", err)
	}
	p, _ = store.GetProxyByAddress("risk-write:8080")
	if p.IPAPIIsScore != 0.65 {
		t.Fatalf("ipapiis_score after update = %v, want 0.65", p.IPAPIIsScore)
	}
	if p.IPAPIFlags != "proxy,hosting" {
		t.Fatalf("ipapi_flags after update = %q, want proxy,hosting", p.IPAPIFlags)
	}
	// 顺带确认出口信息与 fail_count 同路径写入正确（scanProxy 顺序对齐）。
	if p.ExitIP != "203.0.113.20" || p.Latency != 42 {
		t.Fatalf("exit info not written correctly: %#v", p)
	}
}

// TestUpdateProxyExitInfoNegativeScoreDoesNotOverwrite 负向：ipapiis_score 降级/未知(-1)
// 不得覆盖已有有效分；但 ipapi_flags 随每次成功探测覆盖（含清空，语义有效）。
func TestUpdateProxyExitInfoNegativeScoreDoesNotOverwrite(t *testing.T) {
	store := newTestStorage(t)
	insertProxy(t, store, "risk-keep:8080", "http", "us", SourceManual, 100, "active", 0)
	p, err := store.GetProxyByAddress("risk-keep:8080")
	if err != nil {
		t.Fatalf("GetProxyByAddress() error = %v", err)
	}

	// 先写入有效分 0.80 + flags。
	if err := store.UpdateProxyExitInfo(p.ID, "203.0.113.21", "US Ashburn", 42, 0.80, "proxy"); err != nil {
		t.Fatalf("UpdateProxyExitInfo(0.80) error = %v", err)
	}
	// 再以 -1（ipapi.is 探测降级/失败）更新：ipapiis_score 应保持 0.80 不被覆盖。
	if err := store.UpdateProxyExitInfo(p.ID, "203.0.113.21", "US Ashburn", 50, -1, ""); err != nil {
		t.Fatalf("UpdateProxyExitInfo(-1) error = %v", err)
	}
	p, _ = store.GetProxyByAddress("risk-keep:8080")
	if p.IPAPIIsScore != 0.80 {
		t.Fatalf("ipapiis_score after -1 update = %v, want 0.80 (must not overwrite)", p.IPAPIIsScore)
	}
	// 其它字段（如延迟）仍应正常更新，证明只有 ipapiis_score 被条件保护。
	if p.Latency != 50 {
		t.Fatalf("latency after -1 update = %d, want 50", p.Latency)
	}

	// 从 0 分再写 -1 也不得覆盖（0 是有效低分）。
	if err := store.UpdateProxyExitInfo(p.ID, "203.0.113.21", "US Ashburn", 50, 0, ""); err != nil {
		t.Fatalf("UpdateProxyExitInfo(0) error = %v", err)
	}
	if err := store.UpdateProxyExitInfo(p.ID, "203.0.113.21", "US Ashburn", 50, -1, ""); err != nil {
		t.Fatalf("UpdateProxyExitInfo(-1 after 0) error = %v", err)
	}
	p, _ = store.GetProxyByAddress("risk-keep:8080")
	if p.IPAPIIsScore != 0 {
		t.Fatalf("ipapiis_score after 0 then -1 = %v, want 0 (must not overwrite with -1)", p.IPAPIIsScore)
	}
}

// TestProxyColumnsMatchScanProxy 弥补现有缺口：proxyColumns 与 scanProxy 的列数/顺序
// 必须严格一致。用真实 SELECT proxyColumns 后交给 scanProxy 扫描——若列数或类型不匹配，
// rows.Scan 会报错，从而捕获两处失步。
func TestProxyColumnsMatchScanProxy(t *testing.T) {
	store := newTestStorage(t)
	insertProxy(t, store, "cols-sync:8080", "http", "us", SourceManual, 123, "active", 0)
	if err := store.UpdateProxyExitInfo(mustProxyID(t, store, "cols-sync:8080"), "203.0.113.30", "US Ashburn", 123, 0.55, "hosting"); err != nil {
		t.Fatalf("UpdateProxyExitInfo() error = %v", err)
	}

	rows, err := store.db.Query(`SELECT ` + proxyColumns + ` FROM proxies WHERE address = 'cols-sync:8080'`)
	if err != nil {
		t.Fatalf("query proxyColumns: %v", err)
	}
	defer rows.Close()

	// 列数与 scanProxy Scan 目标数须一致：rows.Columns 给出 SELECT 侧列数。
	cols, err := rows.Columns()
	if err != nil {
		t.Fatalf("rows.Columns(): %v", err)
	}
	// 当前模型：id..note(20) + ipapiis_score + ipapi_flags = 22 列。硬编码断言，防止任一侧漏改。
	if len(cols) != 22 {
		t.Fatalf("proxyColumns yields %d columns, want 22 (id..note + ipapiis_score + ipapi_flags)", len(cols))
	}
	// 最后两列必须是 ipapiis_score, ipapi_flags，与 scanProxy 末尾 &p.IPAPIIsScore,&p.IPAPIFlags 对齐。
	if cols[len(cols)-2] != "ipapiis_score" || cols[len(cols)-1] != "ipapi_flags" {
		t.Fatalf("last two SELECT columns = %q,%q, want ipapiis_score,ipapi_flags", cols[len(cols)-2], cols[len(cols)-1])
	}

	if !rows.Next() {
		t.Fatal("expected one row")
	}
	// 核心断言：scanProxy 能无错扫描全部列——列数/顺序不一致会导致 Scan 报错。
	p, err := scanProxy(rows)
	if err != nil {
		t.Fatalf("scanProxy failed (proxyColumns/scanProxy out of sync?): %v", err)
	}
	if p.Address != "cols-sync:8080" || p.IPAPIIsScore != 0.55 || p.IPAPIFlags != "hosting" || p.Latency != 123 {
		t.Fatalf("scanned proxy mismatch: %#v", p)
	}
	if rows.Err() != nil {
		t.Fatalf("rows error: %v", rows.Err())
	}
}

func mustProxyID(t *testing.T, store *Storage, address string) int64 {
	t.Helper()
	p, err := store.GetProxyByAddress(address)
	if err != nil {
		t.Fatalf("GetProxyByAddress(%s) error = %v", address, err)
	}
	return p.ID
}

// TestProxyColumnsIncludesRiskColumns 静态断言 proxyColumns 常量包含两个风险列，
// 便于在两处对齐测试之外提供一个直接的字符串防护。
func TestProxyColumnsIncludesRiskColumns(t *testing.T) {
	if !strings.Contains(proxyColumns, "ipapiis_score") {
		t.Fatal("proxyColumns constant missing ipapiis_score")
	}
	if !strings.Contains(proxyColumns, "ipapi_flags") {
		t.Fatal("proxyColumns constant missing ipapi_flags")
	}
}
