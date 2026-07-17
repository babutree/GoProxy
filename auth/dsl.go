package auth

import (
	"fmt"
	"strings"
	"unicode"
)

const maxSessionLength = 64

// ParsedUsername is the routing DSL carried in proxy auth username.
// Syntax: <base>[-region-<cc>][-unlock-<token>][-session-<id>][-node-<addr>]
// unlock token: gpt|openai|chatgpt|claude|gemini|grok|cf|all
// all => openai+claude+grok+gemini+cf (AND).
type ParsedUsername struct {
	Base    string
	Region  string
	Session string
	// Unlock is the normalized list of required unlock signals (openai/claude/grok/gemini/cf).
	Unlock []string
	// Node pins routing to a single node by its dial address (host:port), i.e. the
	// gateway ENTRANCE address (the node's own identity), not the observed exit IP.
	// When set, selection bypasses session affinity and picks exactly this node,
	// provided it still passes availability/region/unlock/subscription checks.
	// Chained/realm upstreams mean the final exit IP may still differ from this
	// address and may drift; the gateway can only guarantee the entrance node.
	Node string
}

// maxNodeLength caps the pinned node address token length (host:port).
const maxNodeLength = 255

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

// parseNode 解析 -node- 后的入口地址（host:port）。
// 仅锁定网关拨号的入口节点身份；不校验/不保证最终出口 IP（链式/realm 上游可能不同或漂移）。
// host:port 中含冒号，不能复用 splitValue（其按 DSL marker 切分即可，冒号不是 marker）。
func parseNode(raw string) (string, string, error) {
	value, rest := splitValue(raw)
	if value == "" {
		return "", "", fmt.Errorf("node address is empty")
	}
	if len(value) > maxNodeLength {
		return "", "", fmt.Errorf("node address exceeds %d characters", maxNodeLength)
	}
	if !isNodeAddress(value) {
		return "", "", fmt.Errorf("node address must be host:port")
	}
	return value, rest, nil
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
