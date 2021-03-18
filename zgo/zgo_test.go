package zgo

import (
	"fmt"
	"go/ast"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestModuleRoot(t *testing.T) {
	if r := ModuleRoot(); r != "/home/martin/code/zstd" {
		t.Error(r)
	}

	wd, _ := os.Getwd()
	defer os.Chdir(wd)
	os.Chdir("/etc")
	if r := ModuleRoot(); r != "" {
		t.Error(r)
	}
}

func TestTag(t *testing.T) {
	cases := []struct {
		in, inName, want string
		wantAttr         []string
	}{
		{`json:"w00t"`, "json",
			"w00t", nil},
		{`yaml:"w00t"`, "json",
			"Original", nil},
		{`json:"w00t" yaml:"xxx""`, "yaml",
			"xxx", nil},
		{`JSON:"w00t"`, "json",
			"Original", nil},
		{`JSON: "w00t"`, "json",
			"Original", nil},
		{`json:"w00t,omitempty"`, "json",
			"w00t", []string{"omitempty"}},
		{`json:"w00t,omitempty,readonly"`, "json",
			"w00t", []string{"omitempty", "readonly"}},
		{`json:"w00t,"`, "json",
			"w00t", []string{""}},
		{`json:"-"`, "json",
			"-", nil},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			f := &ast.Field{
				Names: []*ast.Ident{&ast.Ident{Name: "Original"}},
				Tag:   &ast.BasicLit{Value: fmt.Sprintf("`%v`", tc.in)}}

			out, attr := Tag(f, tc.inName)
			if out != tc.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}

			if !reflect.DeepEqual(attr, tc.wantAttr) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", attr, tc.wantAttr)
			}
		})
	}

	t.Run("nil", func(t *testing.T) {
		f := &ast.Field{
			Names: []*ast.Ident{&ast.Ident{Name: "Original"}},
		}

		out := TagName(f, "json")
		if out != "Original" {
			t.Errorf("\nout:  %#v\nwant: %#v\n", out, "Original")
		}
	})

	t.Run("nil", func(t *testing.T) {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("didn't panic")
			}
			if !strings.HasPrefix(r.(string), "cannot use TagName on struct with more than one name: ") {
				t.Errorf("wrong message: %#v", r)
			}
		}()

		f := &ast.Field{
			Names: []*ast.Ident{&ast.Ident{Name: "Original"},
				&ast.Ident{Name: "Second"}}}
		_ = TagName(f, "json")
	})

	t.Run("embed", func(t *testing.T) {
		cases := []struct {
			name string
			in   *ast.Field
			want string
		}{
			{
				"notag",
				&ast.Field{
					Tag:  &ast.BasicLit{Value: "`b:\"Bar\"`"},
					Type: &ast.Ident{Name: "Foo"},
				},
				"Foo",
			},
			{
				"ident",
				&ast.Field{Type: &ast.Ident{Name: "Foo"}},
				"Foo",
			},
			{
				"pointer",
				&ast.Field{Type: &ast.StarExpr{X: &ast.Ident{Name: "Foo"}}},
				"Foo",
			},
			{
				"pkg",
				&ast.Field{Type: &ast.SelectorExpr{Sel: &ast.Ident{Name: "Foo"}}},
				"Foo",
			},
			{
				"pkg-pointer",
				&ast.Field{
					Type: &ast.StarExpr{
						X: &ast.SelectorExpr{Sel: &ast.Ident{Name: "Foo"}},
					},
				},
				"Foo",
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				out := TagName(tc.in, "a")
				if out != tc.want {
					t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
				}
			})
		}

	})
}
