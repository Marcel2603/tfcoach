package testutil

type MemSource struct {
	Files map[string]string
}

func (m MemSource) List(_ string) ([]string, error) {
	var paths []string
	for p := range m.Files {
		paths = append(paths, p)
	}
	return paths, nil
}

func (m MemSource) ReadFile(path string) ([]byte, error) {
	return []byte(m.Files[path]), nil
}
