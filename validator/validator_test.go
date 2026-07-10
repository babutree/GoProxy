package validator

import (
	"strings"
	"testing"

	"goproxy/config"
)

func TestDefaultValidationTargetsAvoidHttpbinSinglePoint(t *testing.T) {
	for _, target := range httpsTestTargets {
		if strings.Contains(target, "httpbin.org") {
			t.Fatalf("httpsTestTargets still contains httpbin single-point target: %q", target)
		}
	}
	if len(httpsTestTargets) < 3 {
		t.Fatalf("httpsTestTargets has %d targets, want multiple providers", len(httpsTestTargets))
	}

	cfg := config.DefaultConfig()
	if strings.Contains(cfg.ValidateURL, "httpbin.org") {
		t.Fatalf("default ValidateURL still contains httpbin: %q", cfg.ValidateURL)
	}
	if targets := parseValidateURLs(cfg.ValidateURL); len(targets) < 3 {
		t.Fatalf("default ValidateURL has %d targets, want multiple providers: %q", len(targets), cfg.ValidateURL)
	}
}

func TestParseValidateURLsTrimsEmptyParts(t *testing.T) {
	targets := parseValidateURLs(" http://a.test/ok, ,http://b.test/ok ")
	if len(targets) != 2 || targets[0] != "http://a.test/ok" || targets[1] != "http://b.test/ok" {
		t.Fatalf("parseValidateURLs() = %#v, want trimmed non-empty targets", targets)
	}
}

// TestGeoFilterReadDoesNotRaceWithConfigSave 复现 BUG-58：
// validator 缓存了 config.Get() 返回的 *Config 指针并在 passesGeoFilter 中无锁读取
// 国家名单 slice，同时 config.Save 并发更新全局配置。旧实现下 Save 原地改写
// *globalCfg（含 slice header）会与这里的读取产生 data race；修复后 Save 改为
// 替换 globalCfg 指针，validator 持有的旧快照不可变，-race 下必须干净通过。
func TestGeoFilterReadDoesNotRaceWithConfigSave(t *testing.T) {
	t.Setenv("DATA_DIR", t.TempDir())

	// 建立初始 globalCfg，使 config.Get() 返回一个真实指针。
	base := config.Load()
	base.BlockedCountries = []string{"CN"}
	base.AllowedCountries = nil
	if err := config.Save(base); err != nil {
		t.Fatalf("initial Save() error = %v", err)
	}

	// validator 在此刻捕获 config.Get() 的指针（与生产 New() 路径一致）。
	v := New(4, 1, "http://127.0.0.1/validate")

	const iterations = 2000
	done := make(chan struct{})

	// writer：反复 Save，交替改写国家名单（触发 globalCfg 指针替换）。
	go func() {
		defer close(done)
		for i := 0; i < iterations; i++ {
			cfg := *base
			if i%2 == 0 {
				cfg.BlockedCountries = []string{"CN", "RU", "IR"}
				cfg.AllowedCountries = nil
			} else {
				cfg.BlockedCountries = nil
				cfg.AllowedCountries = []string{"US", "JP", "SG"}
			}
			if err := config.Save(&cfg); err != nil {
				t.Errorf("Save() error = %v", err)
				return
			}
		}
	}()

	// reader：反复经 passesGeoFilter 读取 v.cfg 的国家名单 slice。
	for i := 0; i < iterations; i++ {
		_ = v.passesGeoFilter("US")
		_ = v.passesGeoFilter("CN")
	}

	<-done
}
