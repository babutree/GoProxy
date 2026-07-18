package auth

import (
	"encoding/base64"
	"fmt"
	"strings"
	"unicode"
)

const maxSessionLength = 64

// ParsedUsername is the routing DSL carried in proxy auth username.
// Syntax: <base>[-region-<cc>][-unlock-<token>][-session-<id>][-node-<addr|key-<nodeKey>>]
// unlock token: gpt|openai|chatgpt|claude|gemini|grok|cf|all
// all => openai+claude+grok+gemini+cf (AND).
type ParsedUsername struct {
	Base    string
	Region  string
	Session string
	// Unlock is the normalized list of required unlock signals (openai/claude/grok/gemini/cf).
	Unlock []string
	// Node pins routing. Forms:
	//   - host:port          当前拨号入口（兼容旧复制；隧道本地端口可能被重分配）
	//   - key-<nodeKey>      稳定配置身份（推荐；刷新后端口变仍命中同一上游配置）
	// Not an observed exit IP. Chained/realm upstreams may still change public exit.
	Node string
}

// maxNodeLength caps the pinned node token length (host:port or key-<nodeKey>).
const maxNodeLength = 255
const nodeKeyPrefix = "key-"

func ParseUsername(raw string) (ParsedUsername, error) {
	if raw == "" {
		return ParsedUsername{}, fmt.Errorf("username is empty")
	}

	regionStart := strings.Index(raw, "-region-")
	unlockStart := strings.Index(raw, "-unlock-")
	sessionStart := strings.Index(raw, "-session-")
	nodeStart := strings.Index(raw, "-node-")
	baseEnd := firstMarker(regionStart, unlockStart, sessionStart, nodeStart, len(raw))
	if baseEnd == 0 {
		return ParsedUsername{}, fmt.Errorf("base username is empty")
	}

	parsed := ParsedUsername{Base: raw[:baseEnd]}
	remainder := raw[baseEnd:]
	if remainder == "" {
		return parsed, nil
	}

	if strings.HasPrefix(remainder, "-region-") {
		region, rest, err := parseRegion(remainder[len("-region-"):])
		if err != nil {
			return ParsedUsername{}, err
		}
		parsed.Region = region
		remainder = rest
	}

	if strings.HasPrefix(remainder, "-unlock-") {
		unlock, rest, err := parseUnlock(remainder[len("-unlock-"):])
		if err != nil {
			return ParsedUsername{}, err
		}
		parsed.Unlock = unlock
		remainder = rest
	}

	if strings.HasPrefix(remainder, "-node-") {
		node, rest, err := parseNode(remainder[len("-node-"):])
		if err != nil {
			return ParsedUsername{}, err
		}
		parsed.Node = node
		remainder = rest
	}

	if strings.HasPrefix(remainder, "-session-") {
		session, rest, err := parseSession(remainder[len("-session-"):])
		if err != nil {
			return ParsedUsername{}, err
		}
		parsed.Session = session
		remainder = rest
	}

	if remainder != "" {
		return ParsedUsername{}, fmt.Errorf("invalid username DSL suffix: %s", remainder)
	}

	return parsed, nil
}

func firstMarker(positions ...int) int {
	end := -1
	for _, p := range positions {
		if p < 0 {
			continue
		}
		if end < 0 || p < end {
			end = p
		}
	}
	if end < 0 {
		// last arg is default length when no markers
		return positions[len(positions)-1]
	}
	return end
}

func parseRegion(raw string) (string, string, error) {
	value, rest := splitValue(raw)
	if len(value) != 2 || !isAlpha(value) {
		return "", "", fmt.Errorf("region must be 2 alpha characters")
	}
	return strings.ToLower(value), rest, nil
}

func parseSession(raw string) (string, string, error) {
	value, rest := splitValue(raw)
	if value == "" {
		return "", "", fmt.Errorf("session is empty")
	}
	if len(value) > maxSessionLength {
		return "", "", fmt.Errorf("session exceeds %d characters", maxSessionLength)
	}
	if !isSessionID(value) {
		return "", "", fmt.Errorf("session contains invalid characters")
	}
	return value, rest, nil
}

func parseUnlock(raw string) ([]string, string, error) {
	value, rest := splitValue(raw)
	if value == "" {
		return nil, "", fmt.Errorf("unlock filter is empty")
	}
	// 单段 unlock，token 内用 + 或 , 组合：gpt+cf / gpt,cf
	// 也允许纯 all。
	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == '+' || r == ','
	})
	if len(parts) == 0 {
		return nil, "", fmt.Errorf("unlock filter is empty")
	}
	out := make([]string, 0, 5)
	seen := map[string]bool{}
	add := func(token string) error {
		norm, err := normalizeUnlockToken(token)
		if err != nil {
			return err
		}
		for _, n := range norm {
			if seen[n] {
				continue
			}
			seen[n] = true
			out = append(out, n)
		}
		return nil
	}
	for _, p := range parts {
		if err := add(p); err != nil {
			return nil, "", err
		}
	}
	if len(out) == 0 {
		return nil, "", fmt.Errorf("unlock filter is empty")
	}
	return out, rest, nil
}

func normalizeUnlockToken(token string) ([]string, error) {
	t := strings.ToLower(strings.TrimSpace(token))
	switch t {
	case "gpt", "openai", "chatgpt":
		return []string{"openai"}, nil
	case "claude":
		return []string{"claude"}, nil
	case "gemini":
		return []string{"gemini"}, nil
	case "grok":
		return []string{"grok"}, nil
	case "cf", "cloudflare":
		return []string{"cf"}, nil
	case "all":
		return []string{"openai", "claude", "grok", "gemini", "cf"}, nil
	default:
		return nil, fmt.Errorf("unknown unlock filter: %s", token)
	}
}

func splitValue(raw string) (string, string) {
	regionStart := strings.Index(raw, "-region-")
	unlockStart := strings.Index(raw, "-unlock-")
	sessionStart := strings.Index(raw, "-session-")
	nodeStart := strings.Index(raw, "-node-")
	valueEnd := firstMarker(regionStart, unlockStart, sessionStart, nodeStart, len(raw))
	return raw[:valueEnd], raw[valueEnd:]
}

// parseNode 解析 -node- 后的锁定令牌：
//   - host:port              兼容旧复制（本地端口可能被重分配）
//   - key-<base64url(nodeKey)>  稳定身份（wire 形态不含 :，避免主机名里的 -session-/-region- 误切 DSL）
//
// 不校验/不保证最终出口 IP。
func parseNode(raw string) (string, string, error) {
	// key- 令牌：仅 base64url 字符，可用 splitValue 安全切开后续 -session- 等 marker。
	if strings.HasPrefix(raw, nodeKeyPrefix) {
		value, rest := splitValue(raw)
		enc := value[len(nodeKeyPrefix):]
		if enc == "" || !isNodeKeyWireToken(enc) {
			return "", "", fmt.Errorf("node key must be non-empty base64url")
		}
		decoded, err := base64.RawURLEncoding.DecodeString(enc)
		if err != nil || len(decoded) == 0 {
			return "", "", fmt.Errorf("node key is not valid base64url")
		}
		// 规范化为 key-<decoded 原文> 供选路层统一 TrimPrefix("key-") 后查库。
		return nodeKeyPrefix + string(decoded), rest, nil
	}
	value, rest := splitValue(raw)
	if value == "" {
		return "", "", fmt.Errorf("node pin is empty")
	}
	if len(value) > maxNodeLength {
		return "", "", fmt.Errorf("node pin exceeds %d characters", maxNodeLength)
	}
	if !isNodeAddress(value) {
		return "", "", fmt.Errorf("node pin must be host:port or key-<base64url>")
	}
	return value, rest, nil
}

// EncodeNodeKeyPin 将存储态 NodeKey 编码为 DSL 安全令牌（无 key- 前缀由调用方拼接）。
func EncodeNodeKeyPin(nodeKey string) string {
	nodeKey = strings.TrimSpace(nodeKey)
	if nodeKey == "" {
		return ""
	}
	return base64.RawURLEncoding.EncodeToString([]byte(nodeKey))
}

// isNodeKeyWireToken 校验 DSL 线上 key 令牌：仅 base64url 字符，避免与 -session- 等 marker 碰撞。
func isNodeKeyWireToken(value string) bool {
	if value == "" {
		return false
	}
	for _, r := range value {
		if unicode.IsDigit(r) || isASCIIAlpha(r) || r == '-' || r == '_' {
			continue
		}
		return false
	}
	return true
}

// isNodeAddress 校验 host:port 形态：必须恰有一个冒号，端口为数字，host 非空。
// host 允许字母数字、点、连字符（域名/IPv4）；不支持 IPv6 字面量（含多个冒号），
// 与现有节点 Address 存储口径（host:port 文本键）一致。
func isNodeAddress(value string) bool {
	colon := strings.LastIndex(value, ":")
	if colon <= 0 || colon == len(value)-1 {
		return false
	}
	host, port := value[:colon], value[colon+1:]
	if strings.Contains(host, ":") {
		return false
	}
	for _, r := range host {
		if r == '.' || r == '-' || unicode.IsDigit(r) || isASCIIAlpha(r) {
			continue
		}
		return false
	}
	for _, r := range port {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func isAlpha(value string) bool {
	for _, r := range value {
		if !unicode.IsLetter(r) || r > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func isSessionID(value string) bool {
	for _, r := range value {
		if r == '-' || r == '_' || unicode.IsDigit(r) || isASCIIAlpha(r) {
			continue
		}
		return false
	}
	return true
}

func isASCIIAlpha(r rune) bool {
	return r <= unicode.MaxASCII && unicode.IsLetter(r)
}
