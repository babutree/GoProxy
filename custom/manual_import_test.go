package custom

import (
	"strings"
	"testing"

	"goproxy/storage"
)

func TestImportManualLinksAddsDirectAndSkipsDuplicates(t *testing.T) {
	store := newTestStorage(t)
	m := &Manager{storage: store}

	text := strings.Join([]string{
		"socks5://1.1.1.1:1080 韩国 [首尔]",
		"http://2.2.2.2:8080",
		"https://3.3.3.3:443 note",
		"socks5://1.1.1.1:1080", // duplicate in batch
		"",
		"# comment",
		"trojan://x@bad.example.com:443#no",
	}, "\n")

	r, err := m.ImportManualLinks(text, "us", "batch")
	if err != nil {
		t.Fatalf("ImportManualLinks: %v", err)
	}
	if r.Added != 3 {
		t.Fatalf("added=%d, want 3; result=%+v", r.Added, r)
	}
	if r.Skipped < 1 {
		t.Fatalf("skipped=%d, want >=1 for in-batch dup", r.Skipped)
	}
	if r.Failed < 1 {
		t.Fatalf("failed=%d, want >=1 for tunnel link", r.Failed)
	}

	// Second import of same set should skip all existing.
	r2, err := m.ImportManualLinks("socks5://1.1.1.1:1080\nhttp://2.2.2.2:8080\n", "", "")
	if err != nil {
		t.Fatalf("second import: %v", err)
	}
	if r2.Added != 0 || r2.Skipped < 2 {
		t.Fatalf("second import result=%+v, want added=0 skipped>=2", r2)
	}

	p, err := store.GetProxyByAddress("1.1.1.1:1080")
	if err != nil {
		t.Fatalf("GetProxyByAddress: %v", err)
	}
	if p.Source != storage.SourceManual || p.Protocol != "socks5" {
		t.Fatalf("proxy=%+v", p)
	}
	if p.Region != "us" || p.Note != "batch" {
		t.Fatalf("region/note=%q/%q", p.Region, p.Note)
	}
}

func TestDeleteManualNodesBatch(t *testing.T) {
	store := newTestStorage(t)
	m := &Manager{storage: store, singbox: newSpyShard()}
	if err := m.AddManualNode("http://9.9.9.1:8080", "", "a"); err != nil {
		t.Fatal(err)
	}
	if err := m.AddManualNode("http://9.9.9.2:8080", "", "b"); err != nil {
		t.Fatal(err)
	}
	a, _ := store.GetProxyByAddress("9.9.9.1:8080")
	b, _ := store.GetProxyByAddress("9.9.9.2:8080")
	deleted, errs := m.DeleteManualNodes([]int64{a.ID, b.ID, 999999})
	if deleted != 2 {
		t.Fatalf("deleted=%d, want 2; errs=%v", deleted, errs)
	}
	if len(errs) != 1 {
		t.Fatalf("errs=%v, want 1 failure for missing id", errs)
	}
}
