package webui

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"goproxy/config"
	"goproxy/custom"
	"goproxy/storage"
	"goproxy/validator"
)

type failingSubscriptionFile struct {
	path        string
	writeErr    error
	chmodMode   os.FileMode
	closeCalled bool
}

func (f *failingSubscriptionFile) Name() string { return f.path }

func (f *failingSubscriptionFile) Chmod(mode os.FileMode) error {
	f.chmodMode = mode
	return nil
}

func (f *failingSubscriptionFile) WriteString(string) (int, error) {
	return 0, f.writeErr
}

func (f *failingSubscriptionFile) Close() error {
	f.closeCalled = true
	return nil
}

func TestWriteSubscriptionFilePropagatesWriteErrorAndCleansUp(t *testing.T) {
	path := filepath.Join(t.TempDir(), "subscription.yaml")
	if err := os.WriteFile(path, nil, 0600); err != nil {
		t.Fatalf("create test file: %v", err)
	}
	wantErr := errors.New("forced write failure")
	file := &failingSubscriptionFile{path: path, writeErr: wantErr}

	err := writeSubscriptionFile(file, "subscription content")

	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
	if file.chmodMode != 0644 {
		t.Fatalf("chmod mode = %04o, want 0644", file.chmodMode)
	}
	if !file.closeCalled {
		t.Fatal("file was not closed after write failure")
	}
	if _, statErr := os.Stat(path); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("temporary file still exists after write failure: %v", statErr)
	}
}

// TestAPISubscriptionDeleteDoesNotPreDeleteProxiesWhenSubscriptionDeleteFails
// 回归：handler 不得在 DeleteSubscription 事务外先 DeleteBySubscriptionID。
// 否则 DeleteSubscription 失败后会留下“空订阅 + 代理已丢”的半删状态。
func TestAPISubscriptionDeleteDoesNotPreDeleteProxiesWhenSubscriptionDeleteFails(t *testing.T) {
	server := newTestServer(t)
	// customMgr 非 nil 才会走旧双删路径中的前置清理分支
	server.customMgr = custom.NewManager(server.storage, validator.New(1, 1, "http://127.0.0.1"), &config.Config{})

	subID, err := server.storage.AddSubscription("sub", "https://example.test/webui-delete.yaml", "", "auto", 60, "")
	if err != nil {
		t.Fatalf("AddSubscription() error = %v", err)
	}
	const proxyAddr = "198.51.100.50:8080"
	if err := server.storage.AddProxyWithSource(proxyAddr, "http", storage.SourceSubscription, subID); err != nil {
		t.Fatalf("AddProxyWithSource() error = %v", err)
	}

	// 强制订阅行删除失败，模拟 DeleteSubscription 事务中途失败
	if _, err := server.storage.GetDB().Exec(`
		CREATE TRIGGER fail_webui_subscription_delete
		BEFORE DELETE ON subscriptions
		WHEN OLD.id = ` + fmt.Sprint(subID) + `
		BEGIN
			SELECT RAISE(ABORT, 'forced subscription delete failure');
		END`); err != nil {
		t.Fatalf("create trigger: %v", err)
	}

	req := authenticatedJSONRequest(http.MethodPost, "/api/subscription/delete", fmt.Sprintf(`{"id":%d}`, subID))
	rec := httptest.NewRecorder()
	server.routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusInternalServerError, rec.Body.String())
	}
	if _, err := server.storage.GetSubscription(subID); err != nil {
		t.Fatalf("subscription should remain after failed delete: %v", err)
	}
	if _, err := server.storage.GetProxyByIdentity(proxyAddr, storage.SourceSubscription, subID); err != nil {
		t.Fatalf("proxies must not be pre-deleted outside DeleteSubscription transaction: %v", err)
	}
}
