package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	gitignore "github.com/sabhiram/go-gitignore"
	"github.com/spf13/cobra"
	"trigo/service"
)

// NewRootCommand creates a fresh root command instance.
func NewRootCommand() *cobra.Command {
	var (
		dir     string
		all     bool
		exclude []string
	)

	var excludePatterns []string

	cmd := &cobra.Command{
		Use:   "trigo",
		Short: "Prints out a tree structure of the current directory or given dir",
		Long:  `trigo: prints out the tree structure of the current directory or a given dir`,
		RunE: func(cmd *cobra.Command, args []string) error {
			absDir, err := filepath.Abs(dir)
			if err != nil {
				return err
			}
			compiled := make([]*regexp.Regexp, 0, len(excludePatterns))
			for _, p := range excludePatterns {
				re, err := regexp.Compile(p)
				if err != nil {
					return fmt.Errorf("invalid pattern %q: %w", p, err)
				}
				compiled = append(compiled, re)
			}

			info, err := os.Stat(absDir)
			if err != nil {
				return fmt.Errorf("invalid directory %q: %w", dir, err)
			}
			if !info.IsDir() {
				return fmt.Errorf("%q is not a directory", dir)
			}

			buf := bufio.NewWriter(cmd.OutOrStdout())
			cfg := &service.WalkerConfig{
				ShowHidden:     all,
				Root:           absDir,
				Exclude:        exclude,
				ExcludePattern: compiled,
				Out:            buf,
			}

			if _, err := os.Stat(filepath.Join(absDir, ".git")); err == nil {
				if ignore, err := gitignore.CompileIgnoreFile(filepath.Join(absDir, ".gitignore")); err == nil {
					cfg.Ignore = ignore
				}
			}

			root := service.NewNode(absDir, true)
			fmt.Fprintln(buf, absDir)
			service.WalkTree(root, absDir, "", false, cfg)
			err = buf.Flush()
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&dir, "dir", "d", ".", "Directory to print tree for")
	cmd.Flags().BoolVarP(&all, "all", "a", false, "Show hidden files and directories")
	cmd.Flags().StringArrayVarP(&exclude, "exclude", "e", []string{}, "Exclude files or directories by name")
	cmd.Flags().StringArrayVarP(&excludePatterns, "exclude-expression", "E", []string{}, "Exclude by regex pattern")

	return cmd
}

// Execute runs the root command.
func Execute() error {
	return NewRootCommand().Execute()
}
