package zgo

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/token"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	"zgo.at/zstd/ztest"
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

// This also tests ResolvePackage() and ResolveWildcard().
func TestExpand(t *testing.T) {
	t.Skip() // Broken with modules?

	cases := []struct {
		in      []string
		want    []string
		wantErr string
	}{
		{
			[]string{"fmt"},
			[]string{"fmt"},
			"",
		},
		{
			[]string{"fmt", "fmt"},
			[]string{"fmt"},
			"",
		},
		{
			[]string{"fmt", "net/http"},
			[]string{"fmt", "net/http"},
			"",
		},
		{
			[]string{"net/..."},
			[]string{"net", "net/http", "net/http/cgi", "net/http/cookiejar",
				"net/http/fcgi", "net/http/httptest", "net/http/httptrace",
				"net/http/httputil", "net/http/internal", "net/http/pprof",
				"net/internal/socktest", "net/mail", "net/rpc", "net/rpc/jsonrpc",
				"net/smtp", "net/textproto", "net/url",
			},
			"",
		},
		{
			[]string{"zgo.at/zstd"},
			[]string{"zgo.at/zstd"},
			"",
		},
		//{
		//	[]string{"."},
		//	[]string{"zgo.at/zstd/zgo"},
		//	"",
		//},
		//{
		//	[]string{".."},
		//	[]string{"zgo.at/zstd"},
		//	"",
		//},
		//{
		//	[]string{"../..."},
		//	[]string{
		//		"github.com/teamwork/utils",
		//		"github.com/teamwork/utils/aesutil",
		//		"github.com/teamwork/utils/dbg",
		//		"github.com/teamwork/utils/errorutil",
		//		"github.com/teamwork/utils/goutil",
		//		"github.com/teamwork/utils/httputilx",
		//		"github.com/teamwork/utils/httputilx/header",
		//		"github.com/teamwork/utils/imageutil",
		//		"github.com/teamwork/utils/ioutilx",
		//		"github.com/teamwork/utils/jsonutil",
		//		"github.com/teamwork/utils/maputil",
		//		"github.com/teamwork/utils/mathutil",
		//		"github.com/teamwork/utils/netutil",
		//		"github.com/teamwork/utils/raceutil",
		//		"github.com/teamwork/utils/sliceutil",
		//		"github.com/teamwork/utils/sqlutil",
		//		"github.com/teamwork/utils/stringutil",
		//		"github.com/teamwork/utils/syncutil",
		//		"github.com/teamwork/utils/timeutil",
		//	},
		//	"",
		//},

		// Errors
		{
			[]string{""},
			nil,
			"cannot resolve empty string",
		},
		{
			[]string{"thi.s/will/never/exist"},
			nil,
			`cannot find module providing package thi.s/will/never/exist`,
		},
		{
			[]string{"thi.s/will/never/exist/..."},
			nil,
			`cannot find module providing package thi.s/will/never/exist`,
		},
		{
			[]string{"./doesnt/exist"},
			nil,
			"cannot find package",
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			out, err := Expand(tc.in, build.FindOnly)
			if !ztest.ErrorContains(err, tc.wantErr) {
				t.Fatal(err)
			}

			sort.Strings(tc.want)
			var outPkgs []string
			for _, p := range out {
				outPkgs = append(outPkgs, p.ImportPath)
			}

			if !reflect.DeepEqual(tc.want, outPkgs) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", outPkgs, tc.want)
			}
		})
	}
}

func TestParseFiles(t *testing.T) {
	pkg, err := ResolvePackage("net/http", 0)
	if err != nil {
		t.Fatal(err)
	}

	fset := token.NewFileSet()
	out, err := ParseFiles(fset, pkg.Dir, pkg.GoFiles, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(out) != 1 {
		t.Fatalf("len(out) == %v", len(out))
	}

	for _, pkg := range out {
		if pkg.Name != "http" {
			t.Errorf("name == %v", pkg.Name)
		}

		if len(pkg.Files) < 10 {
			t.Errorf("len(pkg.Files) == %v", len(pkg.Files))
		}
	}
}

func TestResolveImport(t *testing.T) {
	cases := []struct {
		inFile, inPkg, want, wantErr string
	}{
		// Twice to test it works from cache
		{"package main\nimport \"net/http\"\n", "http", "net/http", ""},
		{"package main\nimport \"os\"\n", "os", "os", ""},
		{"package main\nimport xxx \"net/http\"\n", "xxx", "net/http", ""},
		{"package main\nimport \"net/http\"\n", "httpx", "", ""},

		// Make sure it works from vendor
		{"package main\n import \"github.com/teamwork/test\"\n", "test", "github.com/teamwork/test", ""},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			f := ztest.TempFile(t, tc.inFile)

			out, err := ResolveImport(f, tc.inPkg)
			if !ztest.ErrorContains(err, tc.wantErr) {
				t.Fatalf("wrong err: %v", err)
			}
			if out != tc.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}

	t.Run("cache", func(t *testing.T) {
		f := ztest.TempFile(t, "package main\nimport \"net/http\"\n")

		importsCache = make(map[string]map[string]string)
		out, err := ResolveImport(f, "http")
		if err != nil {
			t.Fatal(err)
		}
		if out != "net/http" {
			t.Fatalf("out wrong: %v", out)
		}

		// Second time
		out, err = ResolveImport(f, "http")
		if err != nil {
			t.Fatal(err)
		}
		if out != "net/http" {
			t.Fatalf("out wrong: %v", out)
		}

		if len(importsCache) != 1 {
			t.Error(importsCache)
		}
	})
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
