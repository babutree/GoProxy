package auth

import (
	"reflect"
	"testing"
)

func TestParseUsernameUnlockFilters(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want ParsedUsername
	}{
		{
			name: "gpt only",
			raw:  "username-unlock-gpt",
			want: ParsedUsername{Base: "username", Unlock: []string{"openai"}},
		},
		{
			name: "openai alias",
			raw:  "username-unlock-openai",
			want: ParsedUsername{Base: "username", Unlock: []string{"openai"}},
		},
		{
			name: "claude",
			raw:  "username-unlock-claude",
			want: ParsedUsername{Base: "username", Unlock: []string{"claude"}},
		},
		{
			name: "cf only",
			raw:  "username-unlock-cf",
			want: ParsedUsername{Base: "username", Unlock: []string{"cf"}},
		},
		{
			name: "all expands to five requirements",
			raw:  "username-unlock-all",
			want: ParsedUsername{Base: "username", Unlock: []string{"openai", "claude", "grok", "gemini", "cf"}},
		},
		{
			name: "region then unlock then session",
			raw:  "username-region-us-unlock-gpt-session-abc",
			want: ParsedUsername{Base: "username", Region: "us", Session: "abc", Unlock: []string{"openai"}},
		},
		{
			name: "region unlock all",
			raw:  "username-region-jp-unlock-all",
			want: ParsedUsername{Base: "username", Region: "jp", Unlock: []string{"openai", "claude", "grok", "gemini", "cf"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseUsername(tt.raw)
			if err != nil {
				t.Fatalf("ParseUsername(%q) error = %v", tt.raw, err)
			}
			if got.Base != tt.want.Base || got.Region != tt.want.Region || got.Session != tt.want.Session {
				t.Fatalf("ParseUsername(%q) = %#v, want %#v", tt.raw, got, tt.want)
			}
			if !reflect.DeepEqual(got.Unlock, tt.want.Unlock) {
				t.Fatalf("Unlock = %#v, want %#v", got.Unlock, tt.want.Unlock)
			}
		})
	}
}

func TestParseUsernameRejectsInvalidUnlock(t *testing.T) {
	for _, raw := range []string{
		"username-unlock-",
		"username-unlock-netflix",
		"username-unlock-gpt-unlock-cf", // 单次仅允许一个 unlock 段
		"username-session-x-unlock-gpt", // 顺序：unlock 必须在 session 之前
	} {
		if _, err := ParseUsername(raw); err == nil {
			t.Fatalf("ParseUsername(%q) expected error", raw)
		}
	}
}
