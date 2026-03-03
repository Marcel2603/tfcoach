package engine_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Marcel2603/tfcoach/internal/engine"
)

func createFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir -p %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

// TestTree:
// tempDir/
//   a.tf
//   a.txt
//   modules/
//     m1.tf
//   vendor/
//     v1.tf           <-- should be skipped
//   nested/
//     vendor/
//       v2.tf         <-- should be skipped (matches by name at any depth)
//     deeper.tf
//     a.tf

func TestFileSystem_List_BasicAndSkipDirs(t *testing.T) {
	root := t.TempDir()

	createFile(t, filepath.Join(root, "a.tf"), "a")
	createFile(t, filepath.Join(root, "a.txt"), "not tf")
	createFile(t, filepath.Join(root, "modules", "m1.tf"), "m1")
	createFile(t, filepath.Join(root, "vendor", "v1.tf"), "v1")
	createFile(t, filepath.Join(root, "nested", "vendor", "v2.tf"), "v2")
	createFile(t, filepath.Join(root, "nested", "deeper.tf"), "deep")
	createFile(t, filepath.Join(root, "nested", "a.tf"), "a")

	fs := engine.FileSystem{SkipDirs: []string{"vendor"}}

	got, err := fs.List(root)
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}

	// Expect only .tf files outside any "vendor" directory; order is sorted.
	want := []string{
		filepath.Join(root, "a.tf"),
		filepath.Join(root, "modules", "m1.tf"),
		filepath.Join(root, "nested", "a.tf"),
		filepath.Join(root, "nested", "deeper.tf"),
	}
	if len(got) != len(want) {
		t.Fatalf("List() length = %d, want %d; got=%v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("List()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestFileSystem_ReadFile(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "main.tf")
	content := "terraform {}"
	createFile(t, path, content)

	fs := engine.FileSystem{}
	got, err := fs.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	if string(got) != content {
		t.Errorf("ReadFile() = %q, want %q", got, content)
	}
}
