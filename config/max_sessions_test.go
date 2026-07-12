package config

import (
	"encoding/json"
	"os"
	"testing"
)

func TestDefaultMaxSessionsPerProxyIsUnlimited(t *testing.T) {
	clearConfigEnv(t)
	t.Setenv("DATA_DIR", t.TempDir())
	cfg := DefaultConfig()
	if cfg.MaxSessionsPerProxy != 0 {
		t.Fatalf("MaxSessionsPerProxy = %d, want 0 (unlimited)", cfg.MaxSessionsPerProxy)
	}
}

func TestMaxSessionsPerProxyEnvAndRoundTrip(t *testing.T) {
	clearConfigEnv(t)
	t.Setenv("DATA_DIR", t.TempDir())
	t.Setenv("MAX_SESSIONS_PER_PROXY", "3")

	cfg := DefaultConfig()
	if cfg.MaxSessionsPerProxy != 3 {
		t.Fatalf("from env MaxSessionsPerProxy = %d, want 3", cfg.MaxSessionsPerProxy)
	}

	cfg = Load()
	cfg.MaxSessionsPerProxy = 2
	if err := Save(cfg); err != nil {
		t.Fatalf("Save: %v", err)
	}
	// Clear env so reload uses file only for this field path check.
	t.Setenv("MAX_SESSIONS_PER_PROXY", "")
	reloaded := Load()
	if reloaded.MaxSessionsPerProxy != 2 {
		t.Fatalf("after reload MaxSessionsPerProxy = %d, want 2", reloaded.MaxSessionsPerProxy)
	}

	data, err := os.ReadFile(ConfigFile())
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("json: %v", err)
	}
	if raw["max_sessions_per_proxy"] != float64(2) {
		t.Fatalf("json max_sessions_per_proxy = %v, want 2", raw["max_sessions_per_proxy"])
	}
}

func TestMaxSessionsPerProxyZeroPersistsAsUnlimited(t *testing.T) {
	clearConfigEnv(t)
	t.Setenv("DATA_DIR", t.TempDir())
	cfg := Load()
	cfg.MaxSessionsPerProxy = 5
	if err := Save(cfg); err != nil {
		t.Fatalf("Save 5: %v", err)
	}
	cfg.MaxSessionsPerProxy = 0
	if err := Save(cfg); err != nil {
		t.Fatalf("Save 0: %v", err)
	}
	reloaded := Load()
	if reloaded.MaxSessionsPerProxy != 0 {
		t.Fatalf("MaxSessionsPerProxy = %d, want 0", reloaded.MaxSessionsPerProxy)
	}
}

func TestSaveRejectsNegativeMaxSessionsPerProxy(t *testing.T) {
	clearConfigEnv(t)
	t.Setenv("DATA_DIR", t.TempDir())
	cfg := Load()
	cfg.MaxSessionsPerProxy = -1
	if err := Save(cfg); err == nil {
		t.Fatal("Save(-1) error = nil, want rejection")
	}
}
