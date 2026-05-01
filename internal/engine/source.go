package engine

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Source interface {
	List(root string) (*FileList, error)
	ReadFile(path string) ([]byte, error)
}

type FileList struct {
	TerraformFiles     []string
	TFCoachIgnoreFiles []string
}

type FileSystem struct {
	SkipDirs []string
}

func (f FileSystem) List(root string) (*FileList, error) {
	// TODO later: switch .terragrunt-cache from "completely skipped" to "ignored in issue reporting"?

	skip := map[string]struct{}{}
	for _, d := range f.SkipDirs {
		skip[d] = struct{}{}
	}
	var foundTerraformFiles []string
	var foundIgnoreFiles []string
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if _, ok := skip[d.Name()]; ok {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(p, ".tf") {
			foundTerraformFiles = append(foundTerraformFiles, p)
		}
		if d.Name() == ".tfcoachignore" {
			foundIgnoreFiles = append(foundIgnoreFiles, p)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Also search parent directories of root for .tfcoachignore files, so that
	// running `tfcoach lint subdir` still respects ignore files at the repo root.
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve root path: %w", err)
	}
	for d := filepath.Dir(absRoot); d != filepath.Dir(d); d = filepath.Dir(d) {
		p := filepath.Join(d, ".tfcoachignore")
		if _, err := os.Stat(p); err == nil {
			foundIgnoreFiles = append(foundIgnoreFiles, p)
		}
	}

	sort.Strings(foundTerraformFiles) // deterministic order
	sort.Strings(foundIgnoreFiles)
	return &FileList{
		TerraformFiles:     foundTerraformFiles,
		TFCoachIgnoreFiles: foundIgnoreFiles,
	}, nil
}

func (FileSystem) ReadFile(path string) ([]byte, error) { return os.ReadFile(path) }
