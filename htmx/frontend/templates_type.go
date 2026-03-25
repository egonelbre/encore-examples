package frontend

import (
	"fmt"
	"io/fs"
	"reflect"
	"regexp"
)

var reGoType = regexp.MustCompile(`/\*gotype:\s*(\S+)\s*\*/`)

// goTypeName returns the fully qualified type name for T, dereferencing pointers.
func goTypeName[T any]() string {
	typ := reflect.TypeFor[T]()
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ.PkgPath() + "." + typ.Name()
}

// checkFileGoType returns a check that verifies a page template file
// contains a matching /*gotype: ...*/ comment.
func checkFileGoType[T any](path string) func(fsys fs.FS) error {
	expected := goTypeName[T]()
	return func(fsys fs.FS) error {
		content, err := fs.ReadFile(fsys, path)
		if err != nil {
			return fmt.Errorf("gotype check %s: %w", path, err)
		}
		m := reGoType.FindSubmatch(content)
		if m == nil {
			return fmt.Errorf("template %s: missing /*gotype: ...*/ declaration", path)
		}
		declared := string(m[1])
		if declared != expected {
			return fmt.Errorf("template %s: declared gotype %q, expected %q", path, declared, expected)
		}
		return nil
	}
}

// checkBlockGoType returns a check that verifies a {{define "name"}} block
// contains a matching /*gotype: ...*/ comment.
func checkBlockGoType[T any](name string) func(fsys fs.FS) error {
	expected := goTypeName[T]()
	re := regexp.MustCompile(
		`\{\{-?\s*(?:define|block)\s+"` + regexp.QuoteMeta(name) + `"[^}]*\}\}\s*\{\{-?\s*/\*gotype:\s*(\S+)\s*\*/\s*-?\}\}`,
	)
	return func(fsys fs.FS) error {
		patterns := []string{"templates/*.html", "templates/partials/*.html"}
		for _, pattern := range patterns {
			files, _ := fs.Glob(fsys, pattern)
			for _, f := range files {
				content, err := fs.ReadFile(fsys, f)
				if err != nil {
					continue
				}
				m := re.FindSubmatch(content)
				if m != nil {
					declared := string(m[1])
					if declared != expected {
						return fmt.Errorf("fragment %q in %s: declared gotype %q, expected %q", name, f, declared, expected)
					}
					return nil
				}
			}
		}
		return fmt.Errorf("fragment %q: no /*gotype: ...*/ declaration found", name)
	}
}
