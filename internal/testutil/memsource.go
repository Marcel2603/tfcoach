//go:build test

package testutil

import "github.com/Marcel2603/tfcoach/internal/engine"

type MemSource struct {
	Files map[string]string
}

func (m MemSource) List(_ string) (*engine.FileList, error) {
	var paths []string
	for p := range m.Files {
		paths = append(paths, p)
	}
	return &engine.FileList{TerraformFiles: paths}, nil
}

func (m MemSource) ReadFile(path string) ([]byte, error) {
	return []byte(m.Files[path]), nil
}
