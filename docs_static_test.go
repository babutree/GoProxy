package main

import (
	"os"
	"strings"
	"testing"
)

func TestReadmeProxyExamplesAvoidHttpbinSinglePoint(t *testing.T) {
	data, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatalf("read README.md: %v", err)
	}
	readme := string(data)
	if strings.Contains(readme, "httpbin.org") {
		t.Fatal("README.md still uses httpbin.org in proxy examples")
	}
	if !strings.Contains(readme, "https://www.gstatic.com/generate_204") {
		t.Fatal("README.md missing stable HTTPS proxy example target")
	}
}
