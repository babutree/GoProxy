package proxy

import "testing"

// TestIsBypassTarget 验证内网/本地目标判定：这些地址应直连，不经上游节点转发。
// 覆盖 localhost、.local、IPv4 私网段(RFC1918)、回环、link-local，以及 IPv6 回环/ULA/link-local。
// 目标可能带端口(host:port)或裸 host，可能是域名或 IP。
func TestIsBypassTarget(t *testing.T) {
	bypass := []string{
		// 本地主机名
		"localhost",
		"localhost:8080",
		"myprinter.local",
		"myprinter.local:631",
		// IPv4 回环 127.0.0.0/8
		"127.0.0.1",
		"127.0.0.1:5432",
		"127.1.2.3",
		// IPv4 私网 10.0.0.0/8
		"10.0.0.1",
		"10.255.255.255:22",
		// IPv4 私网 172.16.0.0/12
		"172.16.0.1",
		"172.31.255.255:3306",
		// IPv4 私网 192.168.0.0/16
		"192.168.1.1",
		"192.168.0.100:80",
		// IPv4 link-local 169.254.0.0/16
		"169.254.1.1",
		// IPv6 回环 / ULA / link-local
		"::1",
		"[::1]:8080",
		"fc00::1",
		"[fd12:3456::1]:443",
		"fe80::1",
	}
	for _, h := range bypass {
		if !isBypassTarget(h) {
			t.Errorf("isBypassTarget(%q) = false, want true (内网/本地应直连)", h)
		}
	}

	direct := []string{
		// 公网域名与 IP 必须经上游节点，不 bypass
		"www.google.com",
		"www.google.com:443",
		"1.1.1.1",
		"8.8.8.8:53",
		// 172.32.x 不在私网 172.16/12 段内
		"172.32.0.1",
		// 11.x 不是私网
		"11.0.0.1:80",
		// 公网 IPv6
		"[2001:4860:4860::8888]:443",
	}
	for _, h := range direct {
		if isBypassTarget(h) {
			t.Errorf("isBypassTarget(%q) = true, want false (公网目标必须经上游)", h)
		}
	}
}
