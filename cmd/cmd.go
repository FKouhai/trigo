package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	gitignore "github.com/sabhiram/go-gitignore"
	"github.com/spf13/cobra"
	"trigo/service"
)

// initialize rootCmd

var (
	rootCmd = &cobra.Command{
		Use:   "trigo",
		Short: "Prints out a tree structure of the current directory or given dir",
		Long:  `trigo: prints out the tree structure of the current directory or a given dir`,
		Run: func(cmd *cobra.Command, args []string) {
			absDir, err := filepath.Abs(dir)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			cfg := &service.WalkerConfig{
				ShowHidden: all,
				Root:       absDir,
				Exclude:    exclude,
			}

			if _, err := os.Stat(filepath.Join(absDir, ".git")); err == nil {
				if ignore, err := gitignore.CompileIgnoreFile(filepath.Join(absDir, ".gitignore")); err == nil {
					cfg.Ignore = ignore
				}
			}

			root := service.NewNode(absDir, true)
			fmt.Println(absDir)
			service.WalkTree(root, absDir, "", false, cfg)
		},
	}
	dir     string
	all     bool
	exclude []string
)

func init() {
	rootCmd.Flags().StringVarP(&dir, "dir", "d", ".", "Directory to print tree for")
	rootCmd.Flags().BoolVarP(&all, "all", "a", false, "Show hidden files and directories")
	rootCmd.Flags().StringArrayVarP(&exclude, "exclude", "e", []string{}, "Exclude files or directories by name")
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
