package proxy

import (
	"net"
	"strings"
)

// isBypassTarget 判定目标是否为内网/本地地址，这类目标应直连而非经上游节点转发。
// 目标可能是 "host" 或 "host:port"，host 可能是域名或 IPv4/IPv6 字面量。
//
// 命中规则：
//   - 主机名 localhost（或其子域）、.local 后缀（mDNS）；
//   - IPv4/IPv6 回环（127.0.0.0/8、::1）；
//   - RFC1918 私网（10/8、172.16/12、192.168/16）与 IPv6 ULA（fc00::/7）；
//   - link-local（169.254.0.0/16、fe80::/10）。
//
// 公网域名与公网 IP 一律返回 false，必须经上游节点。
func isBypassTarget(target string) bool {
	host := hostOnly(target)
	if host == "" {
		return false
	}

	lower := strings.ToLower(strings.TrimSuffix(host, "."))
	if lower == "localhost" || strings.HasSuffix(lower, ".localhost") {
		return true
	}
	if strings.HasSuffix(lower, ".local") {
		return true
	}

	ip := net.ParseIP(host)
	if ip == nil {
		// 非 IP 字面量的域名：除上面的本地名外，一律经上游。
		return false
	}
	return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast()
}

// hostOnly 从 "host" 或 "host:port" 中取出 host（去掉 IPv6 方括号）。
// 裸 IPv6 字面量（如 "fc00::1"，含多个冒号且无端口）不会被 SplitHostPort 接受，
// 此时回退为原串去方括号处理。
func hostOnly(target string) string {
	target = strings.TrimSpace(target)
	if target == "" {
		return ""
	}
	if host, _, err := net.SplitHostPort(target); err == nil {
		return host
	}
	// SplitHostPort 失败：可能是裸 host（无端口）或裸 IPv6 字面量。
	return strings.Trim(target, "[]")
}
