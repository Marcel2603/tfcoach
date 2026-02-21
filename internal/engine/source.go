package engine

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/codeglyph/go-dotignore/v2"
)

type Source interface {
	List(root string) ([]string, error)
	ReadFile(path string) ([]byte, error)
}

type FileSystem struct {
	SkipDirs []string
}

func (f FileSystem) List(root string) ([]string, error) {
	ignorer, err := dotignore.NewRepositoryMatcherWithConfig(root, &dotignore.RepositoryConfig{IgnoreFileName: ".tfcoachignore"})
	if err != nil {
		return nil, err
	}

	// TODO idea: maybe allow 2 ignore lists: completely skip (which for now is statically defined in config.go
	//   + optionally TG cache) and ignore in issue reporting only

	// TODO later: switch .terragrunt-cache from "completely skipped" to "ignored in issue reporting"?

	skip := map[string]struct{}{}
	for _, d := range f.SkipDirs {
		skip[d] = struct{}{}
	}
	var out []string
	err = filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		shouldIgnore, ignoreErr := ignorer.Matches(p)
		if ignoreErr != nil {
			return ignoreErr
		}
		if shouldIgnore {
			return nil
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
