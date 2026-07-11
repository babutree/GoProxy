package custom

import "testing"

// TestAssembleConfigEmitsSingleMixedInboundPerNode 是端口合并(缺陷5,双端口2x→单端口)的
// 核心行为测试:每个 tunnel 节点应只生成一个 sing-box `mixed` 入站(单端口同时服务
// SOCKS5 与 HTTP),而非旧的 socks + http 两个入站两个端口。
//
// mixed 入站为 sing-box 官方类型,单端口同时接受 SOCKS4/4a/5 与 HTTP 连接(见官方 docs)。
// 本测试先驱动行为(单 mixed 入站/单端口),兼容当前 3 返回值签名(第三值先以 _ 承接);
// 后续 REFACTOR 步骤再把签名收敛为 (config, portMap) 两值并移除 httpPortMap。
func TestAssembleConfigEmitsSingleMixedInboundPerNode(t *testing.T) {
	s := NewSingBoxProcess("missing-sing-box", t.TempDir(), testSingBoxBasePort)
	nodeA := tunnelNode("a", "a.example.com", "pw-a")
	nodeB := tunnelNode("b", "b.example.com", "pw-b")

	config, portMap := s.assembleConfig([]ParsedNode{nodeA, nodeB})

	inbounds, ok := config["inbounds"].([]map[string]interface{})
	if !ok {
		t.Fatalf("config[inbounds] 类型错误: %T", config["inbounds"])
	}
	// 2 个节点 → 恰好 2 个入站(每节点 1 个),而非 4 个(每节点 socks+http)。
	if len(inbounds) != 2 {
		t.Fatalf("入站数=%d, 期望 2(每节点单 mixed 入站), 说明仍是双端口", len(inbounds))
	}
	for i, in := range inbounds {
		if in["type"] != "mixed" {
			t.Fatalf("入站 %d type=%v, 期望 mixed", i, in["type"])
		}
		if in["listen"] != "127.0.0.1" {
			t.Fatalf("入站 %d listen=%v, 期望 127.0.0.1", i, in["listen"])
		}
	}

	// 端口映射:每个节点恰好一个端口,两两不同,非零。
	if len(portMap) != 2 {
		t.Fatalf("portMap 大小=%d, 期望 2(每节点单端口)", len(portMap))
	}
	pa, pb := portMap[nodeA.NodeKey()], portMap[nodeB.NodeKey()]
	if pa == 0 || pb == 0 {
		t.Fatalf("节点未分配端口: pa=%d pb=%d", pa, pb)
	}
	if pa == pb {
		t.Fatalf("两个节点分配了相同端口 %d", pa)
	}
}
