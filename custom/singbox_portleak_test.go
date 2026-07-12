package custom

import (
	"strings"
	"testing"
)

// TestAssembleConfigReusesFreedLowPort 验证缺陷4（端口泄漏）：
// 当某节点被移除后，其占用的低位端口应被释放并可被新节点复用，
// 而不是让端口分配单调向上爬升、永不回收，最终耗尽端口段。
//
// 场景：第一次仅加载 A（拿到低位端口 pA）；随后 A 下线、B 上线。
// B 作为新节点应复用被释放的 pA，而不是拿到 pA+1。
func TestAssembleConfigReusesFreedLowPort(t *testing.T) {
	s := NewSingBoxProcess("missing-sing-box", t.TempDir(), testSingBoxBasePort)
	nodeA := tunnelNode("a", "a.example.com", "pw-a")

	_, m1 := s.assembleConfig([]ParsedNode{nodeA})
	pA := m1[nodeA.NodeKey()]
	if pA == 0 {
		t.Fatalf("节点 A 未分配端口: m1=%v", m1)
	}

	// 预置上一轮运行态：仅 A 占用端口。
	s.portMap = m1

	// A 下线、B 上线。B 是新节点，应复用被释放的低位端口 pA。
	nodeB := tunnelNode("b", "b.example.com", "pw-b")
	_, m2 := s.assembleConfig([]ParsedNode{nodeB})
	pB := m2[nodeB.NodeKey()]
	if pB == 0 {
		t.Fatalf("节点 B 未分配端口: m2=%v", m2)
	}
	if pB != pA {
		t.Fatalf("端口泄漏：B 应复用被释放的低位端口 %d，实际拿到 %d（端口只增不减）", pA, pB)
	}
}

// TestAssembleConfigStaysWithinSegmentAndSkipsOnOverflow 验证分片端口段溢出保护：
// 每个分片实例独占端口段 [basePort, basePort+portRangeSpan)。
// 当端口段被占满时，新节点必须被跳过（同 buildOutbound 失败的处理方式：
// 不进入 portMap、不生成 inbound/outbound/rule），而不能越界爬进下一个分片的段。
func TestAssembleConfigStaysWithinSegmentAndSkipsOnOverflow(t *testing.T) {
	s := NewSingBoxProcess("missing-sing-box", t.TempDir(), testSingBoxBasePort)

	// 预置一个已有节点 E 占用端口段最顶端（段内最后一个合法端口）。
	nodeE := tunnelNode("e", "e.example.com", "pw-e")
	topPort := testSingBoxBasePort + portRangeSpan - 1
	s.portMap = map[string]int{nodeE.NodeKey(): topPort}

	// E 保留 + 一个不同的新节点。段已满，新节点应被跳过。
	newNode := tunnelNode("new", "new.example.com", "pw-new")
	_, m := s.assembleConfig([]ParsedNode{nodeE, newNode})

	// 段不变式：所有端口必须落在本分片段内，绝不越界到下一段。
	lo := testSingBoxBasePort
	hi := testSingBoxBasePort + portRangeSpan
	for key, p := range m {
		if p < lo || p >= hi {
			t.Fatalf("端口 %d (节点 %s) 越出分片段 [%d, %d)", p, key, lo, hi)
		}
	}

	// E 应保留其原端口。
	if got := m[nodeE.NodeKey()]; got != topPort {
		t.Fatalf("已有节点 E 端口=%d，期望保持 %d", got, topPort)
	}

	// 段已满，新节点必须被跳过（不出现在 portMap 中）。
	if p, ok := m[newNode.NodeKey()]; ok {
		t.Fatalf("段已满时新节点未被跳过，仍被分配端口 %d", p)
	}
}

// TestReloadSegmentFullReturnsIncompleteError 覆盖段满跳过假成功：
// 端口段已满导致新节点不进 portMap 时，Reload 必须返回 error（在启动进程前即可判定），
// 不得 return nil 让上层 RefreshSubscription 删旧代理。
func TestReloadSegmentFullReturnsIncompleteError(t *testing.T) {
	s := NewSingBoxProcess("missing-sing-box", t.TempDir(), testSingBoxBasePort)

	nodeE := tunnelNode("e", "e.example.com", "pw-e")
	topPort := testSingBoxBasePort + portRangeSpan - 1
	s.portMap = map[string]int{nodeE.NodeKey(): topPort}
	s.nodes = []ParsedNode{nodeE}

	newNode := tunnelNode("new", "new.example.com", "pw-new")
	err := s.Reload([]ParsedNode{nodeE, newNode})
	if err == nil {
		t.Fatal("段满跳过节点时 Reload 应返回 incomplete error，得到 nil")
	}
	// 必须是端口不完整类错误，而非仅 binary_not_found（否则段满在缺二进制时被掩盖，
	// 有二进制时仍会假成功）。
	msg := err.Error()
	if !strings.Contains(msg, "不完整") && !strings.Contains(msg, "incomplete") &&
		!strings.Contains(msg, "未分配") && !strings.Contains(msg, "段") {
		t.Fatalf("段满应在启动前以 incomplete/未分配失败，实际: %v", err)
	}
	// 失败后不得把未分配节点伪装进 portMap（restore 后可能仍是旧 E 映射）。
	if _, ok := s.GetPortMap()[newNode.NodeKey()]; ok {
		t.Fatal("失败路径不应保留未分配新节点的端口映射")
	}
}

// TestAssembleConfigExistingNodesKeepStablePorts 回归保护：
// 跨两次重载，已有节点端口必须保持稳定（新增节点不得改动既有节点端口），
// 否则会破坏已入库的代理地址。此用例在修复前后都应为 GREEN（纯回归护栏）。
func TestAssembleConfigExistingNodesKeepStablePorts(t *testing.T) {
	s := NewSingBoxProcess("missing-sing-box", t.TempDir(), testSingBoxBasePort)
	nodeA := tunnelNode("a", "a.example.com", "pw-a")

	_, m1 := s.assembleConfig([]ParsedNode{nodeA})
	pA := m1[nodeA.NodeKey()]
	if pA == 0 {
		t.Fatalf("节点 A 未分配端口: m1=%v", m1)
	}

	s.portMap = m1
	nodeB := tunnelNode("b", "b.example.com", "pw-b")
	_, m2 := s.assembleConfig([]ParsedNode{nodeA, nodeB})

	if got := m2[nodeA.NodeKey()]; got != pA {
		t.Fatalf("已有节点 A 端口从 %d 变为 %d，破坏稳定性", pA, got)
	}
	pB := m2[nodeB.NodeKey()]
	if pB == 0 {
		t.Fatal("新节点 B 未分配端口")
	}
	if pB == pA {
		t.Fatalf("A 与 B 分配了相同端口 %d", pA)
	}
}
