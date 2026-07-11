package storage

import (
	"strings"
	"testing"
)

// TestDualProtocolDefaultsFalse 新节点默认 dual_protocol=false（普通单协议节点）。
func TestDualProtocolDefaultsFalse(t *testing.T) {
	store := newTestStorage(t)
	insertProxy(t, store, "dual-default:8080", "socks5", "us", SourceManual, 100, "active", 0)
	p, err := store.GetProxyByAddress("dual-default:8080")
	if err != nil {
		t.Fatalf("GetProxyByAddress() error = %v", err)
	}
	if p.DualProtocol {
		t.Fatalf("default dual_protocol = true, want false")
	}
}

// TestSetProxyDualProtocolTogglesFlag 置位 dual_protocol 并被 scanProxy 读回。
// 这是 mixed 隧道节点（单端口同时服务 SOCKS5+HTTP）的显式能力标记，
// 供前端可靠区分双协议节点，而非靠地址长相猜测。
func TestSetProxyDualProtocolTogglesFlag(t *testing.T) {
	store := newTestStorage(t)
	insertProxy(t, store, "dual-toggle:8080", "socks5", "us", SourceManual, 100, "active", 0)
	p, err := store.GetProxyByAddress("dual-toggle:8080")
	if err != nil {
		t.Fatalf("GetProxyByAddress() error = %v", err)
	}

	if err := store.SetProxyDualProtocol(p.ID, true); err != nil {
		t.Fatalf("SetProxyDualProtocol(true) error = %v", err)
	}
	p, _ = store.GetProxyByAddress("dual-toggle:8080")
	if !p.DualProtocol {
		t.Fatalf("dual_protocol after set true = false, want true")
	}

	if err := store.SetProxyDualProtocol(p.ID, false); err != nil {
		t.Fatalf("SetProxyDualProtocol(false) error = %v", err)
	}
	p, _ = store.GetProxyByAddress("dual-toggle:8080")
	if p.DualProtocol {
		t.Fatalf("dual_protocol after set false = true, want false")
	}
}

// TestProxyColumnsIncludesDualProtocol 静态断言 proxyColumns 常量包含 dual_protocol 列。
func TestProxyColumnsIncludesDualProtocol(t *testing.T) {
	if !strings.Contains(proxyColumns, "dual_protocol") {
		t.Fatal("proxyColumns constant missing dual_protocol")
	}
}
