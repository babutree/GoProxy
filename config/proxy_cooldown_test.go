package config

import (
	"encoding/json"
	"os"
	"testing"
)

func TestDefaultProxyCooldownMinutesIsZero(t *testing.T) {
	clearConfigEnv(t)
	t.Setenv("DATA_DIR", t.TempDir())
	cfg := DefaultConfig()
	if cfg.ProxyCooldownMinutes != 0 {
		t.Fatalf("ProxyCooldownMinutes = %d, want 0 (disabled)", cfg.ProxyCooldownMinutes)
	}
}

func TestProxyCooldownMinutesEnvAndRoundTrip(t *testing.T) {
	clearConfigEnv(t)
	t.Setenv("DATA_DIR", t.TempDir())
	t.Setenv("PROXY_COOLDOWN_MINUTES", "5")

	cfg := DefaultConfig()
	if cfg.ProxyCooldownMinutes != 5 {
		t.Fatalf("from env ProxyCooldownMinutes = %d, want 5", cfg.ProxyCooldownMinutes)
	}

	cfg = Load()
	cfg.ProxyCooldownMinutes = 3
	if err := Save(cfg); err != nil {
		t.Fatalf("Save: %v", err)
	}
	t.Setenv("PROXY_COOLDOWN_MINUTES", "")
	reloaded := Load()
	if reloaded.ProxyCooldownMinutes != 3 {
		t.Fatalf("after reload ProxyCooldownMinutes = %d, want 3", reloaded.ProxyCooldownMinutes)
	}

	data, err := os.ReadFile(ConfigFile())
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("json: %v", err)
	}
	if raw["proxy_cooldown_minutes"] != float64(3) {
		t.Fatalf("json proxy_cooldown_minutes = %v, want 3", raw["proxy_cooldown_minutes"])
	}
}

func TestProxyCooldownMinutesZeroPersists(t *testing.T) {
	clearConfigEnv(t)
	t.Setenv("DATA_DIR", t.TempDir())
	cfg := Load()
	cfg.ProxyCooldownMinutes = 5
	if err := Save(cfg); err != nil {
		t.Fatalf("Save 5: %v", err)
	}
	cfg.ProxyCooldownMinutes = 0
	if err := Save(cfg); err != nil {
		t.Fatalf("Save 0: %v", err)
	}
	reloaded := Load()
	if reloaded.ProxyCooldownMinutes != 0 {
		t.Fatalf("ProxyCooldownMinutes = %d, want 0", reloaded.ProxyCooldownMinutes)
	}
}

func TestSaveRejectsNegativeProxyCooldownMinutes(t *testing.T) {
	clearConfigEnv(t)
	t.Setenv("DATA_DIR", t.TempDir())
	cfg := Load()
	cfg.ProxyCooldownMinutes = -1
	if err := Save(cfg); err == nil {
		t.Fatal("Save(-1) error = nil, want rejection")
	}
}
