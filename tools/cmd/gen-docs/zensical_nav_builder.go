//go:build tfcoach_tools

// WARNING!
// This file is vibe coded and has only a temporary purpose
// If zensical support awesome-pages, we will drop this file

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type pagesFile struct {
	Title string `yaml:"title"`
	Nav   []any  `yaml:"nav"`
}

// GenerateNav walks docs/pages, resolves .pages nav entries recursively,
// and patches the nav = [...] array inside an existing TOML file.
func GenerateNav(pagesRoot, tomlFile string) {
	nav := buildNav(pagesRoot, pagesRoot)

	var buf bytes.Buffer
	buf.WriteString("nav = ")
	writeTOMLNav(&buf, nav, 0)
	buf.WriteByte('\n')
	navBlock := buf.String()

	existing, err := os.ReadFile(tomlFile)
	if err != nil {
		log.Fatalf("failed to read %s: %v", tomlFile, err)
	}

	patched := replaceNavBlock(string(existing), navBlock)
	if err := os.WriteFile(tomlFile, []byte(patched), 0644); err != nil {
		log.Fatalf("failed to write %s: %v", tomlFile, err)
	}
}

// replaceNavBlock replaces the nav = [...] block in a TOML string.
func replaceNavBlock(toml, navBlock string) string {
	start := strings.Index(toml, "nav = [")
	if start == -1 {
		return toml + "\n" + navBlock
	}
	// find the matching closing bracket
	depth, end := 0, start
	for end < len(toml) {
		switch toml[end] {
		case '[':
			depth++
		case ']':
			depth--
			if depth == 0 {
				return toml[:start] + navBlock + strings.TrimLeft(toml[end+1:], "\n")
			}
		}
		end++
	}
	return toml
}

// writeTOMLNav serialises the nav tree into the TOML inline-table array format
// that zensical.toml uses: [{ "Label" = "path" }, { "Label" = [...] }]
func writeTOMLNav(b *bytes.Buffer, items []any, depth int) {
	indent := strings.Repeat("    ", depth)
	inner := strings.Repeat("    ", depth+1)

	flat := flattenNavItems(items)

	b.WriteString("[\n")
	for i, entry := range flat {
		for label, val := range entry {
			b.WriteString(inner + "{ ")
			fmt.Fprintf(b, "%q = ", label)
			switch child := val.(type) {
			case string:
				fmt.Fprintf(b, "%q", child)
			case []any:
				writeTOMLNav(b, child, depth+1)
			}
			b.WriteString(" }")
		}
		if i < len(flat)-1 {
			b.WriteByte(',')
		}
		b.WriteByte('\n')
	}
	b.WriteString(indent + "]")
}

// flattenNavItems normalises a []any nav slice into []map[string]any,
// expanding any nested []any produced by restEntries.
func flattenNavItems(items []any) []map[string]any {
	var flat []map[string]any
	for _, item := range items {
		switch v := item.(type) {
		case map[string]any:
			flat = append(flat, v)
		case []any:
			flat = append(flat, flattenNavItems(v)...)
		}
	}
	return flat
}

// buildNav resolves the nav for a single directory using its .pages file.
// Falls back to alphabetical order when no .pages file exists.
func buildNav(root, dir string) []any {
	pages := readPagesFile(dir)

	if len(pages.Nav) == 0 {
		return fallbackNav(root, dir)
	}

	var result []any
	for _, entry := range pages.Nav {
		resolved := resolveEntry(root, dir, entry)
		if resolved != nil {
			result = append(result, resolved)
		}
	}
	return result
}

// resolveEntry handles a single nav entry which can be:
//   - "..." (rest placeholder — include unlisted files/dirs alphabetically)
//   - "filename.md" (bare file, use its H1 title)
//   - "Label: filename.md" (explicit label)
//   - "subdirname" (recurse into subdirectory)
//   - map with a single key (Label: path)
func resolveEntry(root, dir string, entry any) any {
	switch v := entry.(type) {
	case string:
		if v == "..." {
			return restEntries(root, dir)
		}
		if strings.HasSuffix(v, ".md") {
			return map[string]any{titleFromFile(filepath.Join(dir, v)): relPath(root, filepath.Join(dir, v))}
		}
		// subdirectory
		subDir := filepath.Join(dir, v)
		if info, err := os.Stat(subDir); err == nil && info.IsDir() {
			label := dirLabel(subDir)
			return map[string]any{label: buildNav(root, subDir)}
		}
	case map[string]any:
		for label, target := range v {
			t, ok := target.(string)
			if !ok {
				continue
			}
			if strings.HasSuffix(t, ".md") {
				return map[string]any{label: relPath(root, filepath.Join(dir, t))}
			}
			subDir := filepath.Join(dir, t)
			if info, err := os.Stat(subDir); err == nil && info.IsDir() {
				return map[string]any{label: buildNav(root, subDir)}
			}
		}
	}
	return nil
}

// restEntries returns nav entries for all files/dirs in dir not already listed in .pages nav.
func restEntries(root, dir string) []any {
	pages := readPagesFile(dir)
	listed := listedEntries(pages.Nav)

	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("failed to read dir %s: %v", dir, err)
	}

	var result []any
	for _, e := range entries {
		name := e.Name()
		if name == ".pages" || listed[name] {
			continue
		}
		if e.IsDir() {
			subDir := filepath.Join(dir, name)
			if !hasMDContent(subDir) {
				continue
			}
			label := dirLabel(subDir)
			result = append(result, map[string]any{label: buildNav(root, subDir)})
		} else if strings.HasSuffix(name, ".md") {
			full := filepath.Join(dir, name)
			result = append(result, map[string]any{titleFromFile(full): relPath(root, full)})
		}
	}
	return result
}

// fallbackNav returns alphabetically sorted nav entries for a directory without a .pages file.
func fallbackNav(root, dir string) []any {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("failed to read dir %s: %v", dir, err)
	}

	var result []any
	for _, e := range entries {
		name := e.Name()
		if name == ".pages" {
			continue
		}
		if e.IsDir() {
			subDir := filepath.Join(dir, name)
			if !hasMDContent(subDir) {
				continue
			}
			label := dirLabel(subDir)
			result = append(result, map[string]any{label: buildNav(root, subDir)})
		} else if strings.HasSuffix(name, ".md") {
			full := filepath.Join(dir, name)
			result = append(result, map[string]any{titleFromFile(full): relPath(root, full)})
		}
	}
	return result
}

// hasMDContent reports whether a directory (recursively) contains any .md files.
func hasMDContent(dir string) bool {
	has := false
	_ = filepath.WalkDir(dir, func(_ string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".md") {
			has = true
			return filepath.SkipAll
		}
		return nil
	})
	return has
}

func readPagesFile(dir string) pagesFile {
	data, err := os.ReadFile(filepath.Join(dir, ".pages"))
	if err != nil {
		return pagesFile{}
	}
	var p pagesFile
	if err := yaml.Unmarshal(data, &p); err != nil {
		log.Fatalf("failed to parse .pages in %s: %v", dir, err)
	}
	return p
}

// listedEntries returns a set of filenames/dirnames explicitly listed in a nav.
func listedEntries(nav []any) map[string]bool {
	listed := make(map[string]bool)
	for _, entry := range nav {
		switch v := entry.(type) {
		case string:
			if v != "..." {
				listed[v] = true
			}
		case map[string]any:
			for _, target := range v {
				if t, ok := target.(string); ok {
					listed[t] = true
				}
			}
		}
	}
	return listed
}

// titleFromFile reads the first H1 heading from a markdown file, falling back to the filename stem.
func titleFromFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return stemName(path)
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}
	return stemName(path)
}

// dirLabel returns the title from a directory's .pages file, falling back to the dir name.
func dirLabel(dir string) string {
	p := readPagesFile(dir)
	if p.Title != "" {
		return p.Title
	}
	return filepath.Base(dir)
}

func stemName(path string) string {
	base := filepath.Base(path)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

func relPath(root, full string) string {
	rel, err := filepath.Rel(root, full)
	if err != nil {
		log.Fatalf("failed to compute relative path: %v", err)
	}
	return rel
}
