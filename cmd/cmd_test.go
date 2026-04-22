package cmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"trigo/cmd"
)

func execute(c *cobra.Command, args ...string) (stdout, stderr string, err error) {
	var outBuf, errBuf bytes.Buffer
	c.SetOut(&outBuf)
	c.SetErr(&errBuf)
	c.SetArgs(args)
	err = c.Execute()
	return outBuf.String(), errBuf.String(), err
}

func TestRootCommand(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		setup      func(t *testing.T, dir string)
		wantOut    []string
		wantNotOut []string
		wantErr    bool
	}{
		{
			name: "default dir prints tree",
			args: []string{},
			setup: func(t *testing.T, dir string) {
				os.WriteFile(filepath.Join(dir, "file.txt"), []byte("hello"), 0o644)
			},
			wantOut: []string{"file.txt"},
		},
		{
			name: "--all includes hidden files",
			args: []string{"--all"},
			setup: func(t *testing.T, dir string) {
				os.WriteFile(filepath.Join(dir, ".hidden"), []byte("secret"), 0o644)
			},
			wantOut: []string{".hidden"},
		},
		{
			name: "--exclude skips entries",
			args: []string{"--exclude", "skipme"},
			setup: func(t *testing.T, dir string) {
				os.MkdirAll(filepath.Join(dir, "skipme"), 0o755)
				os.WriteFile(filepath.Join(dir, "keep.txt"), []byte("keep"), 0o644)
			},
			wantOut:    []string{"keep.txt"},
			wantNotOut: []string{"skipme"},
		},
		{
			name:    "invalid directory returns error",
			args:    []string{"--dir", "/nonexistent/path/12345"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			if len(tt.args) == 0 || (len(tt.args) > 0 && tt.args[0] != "--dir") {
				origWd, _ := os.Getwd()
				os.Chdir(tmpDir)
				defer os.Chdir(origWd)
			}

			if tt.setup != nil {
				tt.setup(t, tmpDir)
			}

			c := cmd.NewRootCommand()
			stdout, _, err := execute(c, tt.args...)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			for _, s := range tt.wantOut {
				if !strings.Contains(stdout, s) {
					t.Errorf("expected stdout to contain %q, got:\n%s", s, stdout)
				}
			}
			for _, s := range tt.wantNotOut {
				if strings.Contains(stdout, s) {
					t.Errorf("expected stdout NOT to contain %q, got:\n%s", s, stdout)
				}
			}
		})
	}
}
