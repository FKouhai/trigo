package service_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
	"trigo/service"
)

func defaultCfg(root string) *service.WalkerConfig {
	return &service.WalkerConfig{
		ShowHidden: false,
		Root:       root,
		Ignore:     nil,
	}
}

func TestNewNode(t *testing.T) {
	tests := []struct {
		name  string
		isDir bool
	}{
		{"file.txt", false},
		{"mydir", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.NewNode(tt.name, tt.isDir)
			if got.Name != tt.name {
				t.Errorf("NewNode().Name = %q, want %q", got.Name, tt.name)
			}
			if got.IsDir != tt.isDir {
				t.Errorf("NewNode().IsDir = %v, want %v", got.IsDir, tt.isDir)
			}
			if len(got.Children) != 0 {
				t.Errorf("NewNode().Children = %v, want empty", got.Children)
			}
		})
	}
}

func TestFSNode_DirContents(t *testing.T) {
	tmpDir := t.TempDir()
	os.MkdirAll(filepath.Join(tmpDir, "subdir"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("hello"), 0644)
	os.WriteFile(filepath.Join(tmpDir, ".hidden"), []byte("hidden"), 0644)

	t.Run("hidden files excluded by default", func(t *testing.T) {
		f := service.NewNode(tmpDir, true)
		f.DirContents(tmpDir, defaultCfg(tmpDir))
		if len(f.Children) != 2 {
			t.Errorf("expected 2 children (hidden excluded), got %d", len(f.Children))
		}
	})

	t.Run("hidden files included with ShowHidden", func(t *testing.T) {
		f := service.NewNode(tmpDir, true)
		f.DirContents(tmpDir, &service.WalkerConfig{ShowHidden: true, Root: tmpDir})
		if len(f.Children) != 3 {
			t.Errorf("expected 3 children (hidden included), got %d", len(f.Children))
		}
	})

	t.Run("excluded entries are skipped", func(t *testing.T) {
		f := service.NewNode(tmpDir, true)
		f.DirContents(tmpDir, &service.WalkerConfig{ShowHidden: true, Root: tmpDir, Exclude: []string{"subdir"}})
		for _, child := range f.Children {
			if child.Name == "subdir" {
				t.Error("expected subdir to be excluded")
			}
		}
	})

	t.Run("excluded files are also skipped", func(t *testing.T) {
		f := service.NewNode(tmpDir, true)
		f.DirContents(tmpDir, &service.WalkerConfig{ShowHidden: true, Root: tmpDir, Exclude: []string{"file.txt"}})
		for _, child := range f.Children {
			if child.Name == "file.txt" {
				t.Error("expected file.txt to be excluded")
			}
		}
	})

	t.Run("children not recursively populated", func(t *testing.T) {
		f := service.NewNode(tmpDir, true)
		f.DirContents(tmpDir, defaultCfg(tmpDir))
		for _, child := range f.Children {
			if len(child.Children) != 0 {
				t.Errorf("child %q should have no pre-populated children, got %d", child.Name, len(child.Children))
			}
		}
	})
}

func TestWalkTree(t *testing.T) {
	tmpDir := t.TempDir()
	os.MkdirAll(filepath.Join(tmpDir, "subdir"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("hello"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "subdir", "nested.txt"), []byte("world"), 0644)
	os.WriteFile(filepath.Join(tmpDir, ".hidden"), []byte("hidden"), 0644)

	captureOutput := func(fn func()) string {
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w
		fn()
		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		os.Stdout = old
		return buf.String()
	}

	t.Run("produces output", func(t *testing.T) {
		root := service.NewNode(tmpDir, true)
		out := captureOutput(func() {
			service.WalkTree(root, tmpDir, "", false, defaultCfg(tmpDir))
		})
		if out == "" {
			t.Error("WalkTree() produced no output")
		}
	})

	t.Run("hidden files not in output by default", func(t *testing.T) {
		root := service.NewNode(tmpDir, true)
		out := captureOutput(func() {
			service.WalkTree(root, tmpDir, "", false, defaultCfg(tmpDir))
		})
		if contains(out, ".hidden") {
			t.Error("expected .hidden to be excluded from output")
		}
	})

	t.Run("hidden files in output with ShowHidden", func(t *testing.T) {
		root := service.NewNode(tmpDir, true)
		out := captureOutput(func() {
			service.WalkTree(root, tmpDir, "", false, &service.WalkerConfig{ShowHidden: true, Root: tmpDir})
		})
		if !contains(out, ".hidden") {
			t.Error("expected .hidden to appear in output with ShowHidden=true")
		}
	})

	t.Run("excluded entries not in output", func(t *testing.T) {
		root := service.NewNode(tmpDir, true)
		out := captureOutput(func() {
			service.WalkTree(root, tmpDir, "", false, &service.WalkerConfig{ShowHidden: true, Root: tmpDir, Exclude: []string{"subdir"}})
		})
		if contains(out, "subdir") {
			t.Error("expected subdir to be excluded from output")
		}
	})

	t.Run("nil node does not panic", func(t *testing.T) {
		service.WalkTree(nil, "", "", false, defaultCfg(""))
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
