package webui

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
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
