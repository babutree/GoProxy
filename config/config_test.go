package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultConfigUsesPRDPortsAndGatewaySettings(t *testing.T) {
	clearConfigEnv(t)
	t.Setenv("DATA_DIR", t.TempDir())

	cfg := DefaultConfig()

	if cfg.HTTPPort != ":7802" {
		t.Fatalf("HTTPPort = %q, want :7802", cfg.HTTPPort)
	}
	if cfg.SOCKS5Port != ":7801" {
		t.Fatalf("SOCKS5Port = %q, want :7801", cfg.SOCKS5Port)
	}
	if cfg.WebUIPort != ":7800" {
		t.Fatalf("WebUIPort = %q, want :7800", cfg.WebUIPort)
	}
	if cfg.SessionTTLMinutes != 10 {
		t.Fatalf("SessionTTLMinutes = %d, want 10", cfg.SessionTTLMinutes)
	}
	if cfg.DefaultRegion != "" {
		t.Fatalf("DefaultRegion = %q, want empty", cfg.DefaultRegion)
	}
	if cfg.ProxyAuthUsername != "acct" {
		t.Fatalf("ProxyAuthUsername = %q, want acct", cfg.ProxyAuthUsername)
	}
}

func TestSaveLoadRoundTripUsesActiveGatewayFieldsOnly(t *testing.T) {
	clearConfigEnv(t)
	t.Setenv("DATA_DIR", t.TempDir())

	cfg := Load()
	cfg.HTTPPort = ":9100"
	cfg.SOCKS5Port = ":9101"
	cfg.WebUIPort = ":9102"
	cfg.SessionTTLMinutes = 15
	cfg.DefaultRegion = " jp "
	cfg.HealthIntervalMinutes = 7
	cfg.MaxRetry = 4
	cfg.SingBoxPath = "D:/tools/sing-box.exe"
	cfg.AllowedCountries = []string{" jp ", "US", "us", "bad"}
	cfg.BlockedCountries = []string{" cn "}

	if err := Save(cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	reloaded := Load()
	if reloaded.HTTPPort != cfg.HTTPPort || reloaded.SOCKS5Port != cfg.SOCKS5Port || reloaded.WebUIPort != cfg.WebUIPort {
		t.Fatalf("ports after reload = %q/%q/%q", reloaded.HTTPPort, reloaded.SOCKS5Port, reloaded.WebUIPort)
	}
	if reloaded.SessionTTLMinutes != cfg.SessionTTLMinutes || reloaded.DefaultRegion != "JP" {
		t.Fatalf("session/default region after reload = %d/%q", reloaded.SessionTTLMinutes, reloaded.DefaultRegion)
	}
	if strings.Join(reloaded.AllowedCountries, ",") != "JP,US" || strings.Join(reloaded.BlockedCountries, ",") != "CN" {
		t.Fatalf("countries after reload = allowed:%#v blocked:%#v", reloaded.AllowedCountries, reloaded.BlockedCountries)
	}
	if reloaded.HealthIntervalMinutes != cfg.HealthIntervalMinutes || reloaded.MaxRetry != cfg.MaxRetry || reloaded.SingBoxPath != cfg.SingBoxPath {
		t.Fatalf("runtime settings after reload = %d/%d/%q", reloaded.HealthIntervalMinutes, reloaded.MaxRetry, reloaded.SingBoxPath)
	}

	assertConfigJSONOmitsLegacyFields(t, ConfigFile())
}

func TestSaveDoesNotPersistPlainProxyPassword(t *testing.T) {
	clearConfigEnv(t)
	t.Setenv("DATA_DIR", t.TempDir())
	cfg := Load()
	cfg.ProxyAuthPassword = "do-not-write"
	cfg.ProxyAuthPasswordHash = passwordHash("do-not-write")

	if err := Save(cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	data, err := os.ReadFile(ConfigFile())
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if strings.Contains(string(data), "do-not-write") || strings.Contains(string(data), "proxy_auth_password\"") {
		t.Fatalf("config persisted plain proxy password: %s", string(data))
	}
	if got := Get(); got.ProxyAuthPassword != "" || got.ProxyAuthPasswordHash == "" {
		t.Fatalf("global credential state = password:%q hash:%q", got.ProxyAuthPassword, got.ProxyAuthPasswordHash)
	}
}

func TestSaveFailureDoesNotPolluteGlobalConfig(t *testing.T) {
	clearConfigEnv(t)
	t.Setenv("DATA_DIR", t.TempDir())
	oldCfg := Load()
	oldCfg.ProxyAuthUsername = "old"
	if err := Save(oldCfg); err != nil {
		t.Fatalf("Save(old) error = %v", err)
	}

	badDataDir := filepath.Join(t.TempDir(), "not-a-dir")
	if err := os.WriteFile(badDataDir, []byte("file"), 0644); err != nil {
		t.Fatalf("create bad DATA_DIR file: %v", err)
	}
	t.Setenv("DATA_DIR", badDataDir)
	newCfg := *oldCfg
	newCfg.ProxyAuthUsername = "new"
	newCfg.SessionTTLMinutes = oldCfg.SessionTTLMinutes + 1

	if err := Save(&newCfg); err == nil {
		t.Fatal("Save() error = nil, want write failure")
	}
	if got := Get(); got.ProxyAuthUsername != "old" || got.SessionTTLMinutes != oldCfg.SessionTTLMinutes {
		t.Fatalf("global config after failed Save = user:%q ttl:%d, want old/%d", got.ProxyAuthUsername, got.SessionTTLMinutes, oldCfg.SessionTTLMinutes)
	}
}

func TestSaveLoadPersistsZeroMaxRetry(t *testing.T) {
	clearConfigEnv(t)
	t.Setenv("DATA_DIR", t.TempDir())
	cfg := Load()
	cfg.MaxRetry = 0

	if err := Save(cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	reloaded := Load()
	if reloaded.MaxRetry != 0 {
		t.Fatalf("MaxRetry after reload = %d, want 0", reloaded.MaxRetry)
	}
}

func clearConfigEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		"WEBUI_PASSWORD", "PROXY_AUTH_ENABLED", "PROXY_AUTH_USERNAME", "PROXY_AUTH_PASSWORD",
		"HTTP_PORT", "SOCKS5_PORT", "WEBUI_PORT", "SESSION_TTL_MINUTES", "DEFAULT_REGION",
		"ALLOWED_COUNTRIES", "BLOCKED_COUNTRIES", "HEALTH_CHECK_INTERVAL", "MAX_RETRY", "SINGBOX_PATH",
	} {
		t.Setenv(key, "")
	}
}

func assertConfigJSONOmitsLegacyFields(t *testing.T, path string) {
	t.Helper()
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read config file: %v", err)
	}
	for _, legacy := range []string{"pool_", "fetch_interval", "custom_proxy_mode"} {
		if strings.Contains(string(data), legacy) {
			t.Fatalf("saved config contains legacy field marker %q: %s", legacy, string(data))
		}
	}
	var parsed map[string]any
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("saved config is not valid json: %v", err)
	}
}
