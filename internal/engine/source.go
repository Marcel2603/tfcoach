package engine

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Source interface {
	List(root string) ([]string, error)
	ReadFile(path string) ([]byte, error)
}

type FileSystem struct {
	SkipDirs []string
}

func (f FileSystem) List(root string) ([]string, error) {
	skip := map[string]struct{}{}
	for _, d := range f.SkipDirs {
		skip[d] = struct{}{}
	}
	var out []string
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
			out = append(out, p)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(out) // deterministic order
	return out, nil
}

func (FileSystem) ReadFile(path string) ([]byte, error) { return os.ReadFile(path) }
