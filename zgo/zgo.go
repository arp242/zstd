// Package zgo provides functions to work with Go source files.
package zgo

import (
	"fmt"
	"go/ast"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// ModuleRoot gets the full path to the module root directory.
//
// Returns empty string if it can't find a module.
func ModuleRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	dir, err = filepath.Abs(dir)
	if err != nil {
		return ""
	}

	pdir := dir
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)

		/// Parent directory is identical: we reached the top of the filesystem
		/// hierarchy and didn't find anything.
		if dir == pdir {
			return ""
		}
		pdir = dir
	}
}

// Tag gets the tag name for a struct field and all attributes.
//
// It will return the struct field name if there is no tag. This function does
// not do any validation on the tag format. Use go vet!
func Tag(f *ast.Field, n string) (string, []string) {
	// For e.g.:
	//  A, B string `json:"x"`
	//
	// Most (all?) marshallers and such will simply skip this anyway as
	// duplicate keys usually doesn't make too much sense.
	if len(f.Names) > 1 {
		panic(fmt.Sprintf("cannot use TagName on struct with more than one name: %v",
			f.Names))
	}

	if f.Tag == nil {
		if len(f.Names) == 0 {
			return getEmbedName(f.Type), nil
		}
		return f.Names[0].Name, nil
	}

	tag := reflect.StructTag(strings.Trim(f.Tag.Value, "`")).Get(n)
	if tag == "" {
		if len(f.Names) == 0 {
			return getEmbedName(f.Type), nil
		}
		return f.Names[0].Name, nil
	}

	if p := strings.Index(tag, ","); p != -1 {
		return tag[:p], strings.Split(tag[p+1:], ",")
	}
	return tag, nil
}

// TagName gets the tag name for a struct field without any attributes.
//
// It will return the struct field name if there is no tag. This function does
// not do any validation on the tag format. Use go vet!
func TagName(f *ast.Field, n string) string {
	name, _ := Tag(f, n)
	return name
}

// StarExpr resolves *ast.StarExpr
func StarExpr(e ast.Expr) (ast.Expr, bool) {
	s, ok := e.(*ast.StarExpr)
	if ok {
		return s.X, true
	}
	return e, false
}

// Embedded struct:
//
//	Foo `json:"foo"`
func getEmbedName(f ast.Expr) string {
start:
	switch t := f.(type) {
	case *ast.StarExpr:
		f = t.X
		goto start
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return t.Sel.Name
	default:
		panic(fmt.Sprintf("can't get name for %#v", f))
	}
}

// PredeclaredType reports if a type is a predeclared built-in type.
//
// Note that this excludes composite types, such as maps, slices, channels, etc.
//
// https://golang.org/ref/spec#Predeclared_identifiers
func PredeclaredType(n string) bool {
	switch n {
	case "bool", "byte", "complex64", "complex128", "error", "float32",
		"float64", "int", "int8", "int16", "int32", "int64", "rune", "string",
		"uint", "uint8", "uint16", "uint32", "uint64", "uintptr",
		"any", "comparable":
		return true
	default:
		return false
	}
}
