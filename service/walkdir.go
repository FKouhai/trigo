// Package service contains the glue/logic for trigo
package service

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	gitignore "github.com/sabhiram/go-gitignore"
)

// WalkerConfig type to add extra options to our tree walker
type WalkerConfig struct {
	ShowHidden bool
	Root       string
	Ignore     *gitignore.GitIgnore
	Exclude    []string
}

// FSNode struct represents the filesystem into a tree holding the required information
// to to add new fs nodes
type FSNode struct {
	IsDir    bool
	Name     string
	Children []*FSNode
}

// NewNode creates a new node inside of our FS tree
func NewNode(name string, isDir bool) *FSNode {
	return &FSNode{
		Name:     name,
		IsDir:    isDir,
		Children: []*FSNode{},
	}
}

// DirContents lazily populates the tree with child nodes
func (f *FSNode) DirContents(dirName string, cfg *WalkerConfig) {
	dir, err := os.Open(dirName)
	if err != nil {
		return
	}
	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return
	}

	for _, fileInfo := range fileInfos {
		name := fileInfo.Name()

		if !cfg.ShowHidden && strings.HasPrefix(name, ".") {
			continue
		}

		// Check if the --exclude flag has been passed if so skip those directories
		// Even when --all is given
		if slices.Contains(cfg.Exclude, name) {
			continue
		}

		if cfg.Ignore != nil {
			relPath, _ := filepath.Rel(cfg.Root, filepath.Join(dirName, name))
			// Excludes files found in the relative path based on the .gitignore
			if cfg.Ignore.MatchesPath(relPath) {
				continue
			}
		}

		child := NewNode(fileInfo.Name(), fileInfo.IsDir())
		f.Children = append(f.Children, child)
	}
}

// WalkTree walks through the dynamically created tree by NewNode and DirContents
func WalkTree(f *FSNode, dirPath string, prefix string, isLast bool, cfg *WalkerConfig) {
	if f == nil {
		return
	}
	connector := "|--"
	childPrefix := prefix + "|   "

	if isLast {
		connector = "L__"
		childPrefix = prefix + "    "
	}

	fmt.Println(prefix + connector + f.Name)

	// in case a directory is found lazily create new nodes in the tree
	if f.IsDir {
		f.DirContents(dirPath, cfg)
	}

	for i, child := range f.Children {
		WalkTree(child, filepath.Join(dirPath, child.Name), childPrefix, i == len(f.Children)-1, cfg)
	}
}
