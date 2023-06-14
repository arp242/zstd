package zreflect

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"zgo.at/zstd/ztype"
)

func TestTag(t *testing.T) {
	tests := []struct {
		in      reflect.StructField
		tag     string
		wantTag string
		wantOpt []string
	}{
		{
			func() reflect.StructField {
				return reflect.TypeOf(struct {
					XXX string `json:"xxx" db:"yyy"`
				}{}).Field(0)
			}(),
			"json",
			"xxx", nil,
		},

		{
			func() reflect.StructField {
				return reflect.TypeOf(struct {
					XXX string `json:"xxx,opt1,opt2" db:"yyy"`
				}{}).Field(0)
			}(),
			"json",
			"xxx", []string{"opt1", "opt2"},
		},

		{
			func() reflect.StructField {
				return reflect.TypeOf(struct {
					XXX string `db:"yyy"`
				}{}).Field(0)
			}(),
			"json",
			"", nil,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			tag, opt := Tag(tt.in, tt.tag)

			if tag != tt.wantTag {
				t.Errorf("\nhave: %q\nwant: %q", tag, tt.wantTag)
			}
			if !reflect.DeepEqual(opt, tt.wantOpt) {
				t.Errorf("\nhave: %#v\nwant: %#v", opt, tt.wantOpt)
			}
		})
	}
}

type (
	EmbedMe struct {
		FE string `db:"fe"`
	}
	EmbedUnexport struct {
		FE    string `db:"fe"`
		unexp int
	}
	embedUnexport struct {
		FE    string `db:"fe"`
		unexp int
	}
)

func TestFields(t *testing.T) {
	tests := []struct {
		in         any
		wantNames  []string
		wantValues []any
		wantOpts   [][]string
	}{
		{struct{}{}, []string{}, []any{}, [][]string{}},

		{ztype.Ptr(struct {
			F1 int `db:"f1,opt"`
		}{42}),
			[]string{"f1"},
			[]any{42},
			[][]string{{"opt"}}},

		{struct {
			F1 int `db:"f1"`
			F2 int `db:"f2,skip"`
			F3 int
			F4 int `db:",skip"`
			F5 int `db:"-"`
			F6 int `db:"f6,opt"`
			F7 int `db:"f7,opt,other"`
		}{1, 2, 3, 4, 5, 6, 7},
			[]string{"f1", "F3", "f6", "f7"},
			[]any{1, 3, 6, 7},
			[][]string{nil, nil, []string{"opt"}, []string{"opt", "other"}}},

		{&struct {
			F1 int `db:"f1"`
			F2 int `db:"f2,skip"`
			F3 int
			F4 int `db:",skip"`
			F5 int `db:"-"`
			F6 *int
		}{1, 2, 3, 4, 5, ztype.Ptr(6)},
			[]string{"f1", "F3", "F6"},
			[]any{1, 3, ztype.Ptr(6)},
			[][]string{nil, nil, nil}},

		{struct {
			N struct{ I int }
		}{struct{ I int }{42}},
			[]string{"N"},
			[]any{struct{ I int }{42}},
			[][]string{nil}},

		{struct {
			EmbedMe
			F1 int `db:"f1"`
		}{EmbedMe{"XXX"}, 666},
			[]string{"fe", "f1"},
			[]any{"XXX", 666},
			[][]string{nil, nil}},

		{struct {
			EmbedMe EmbedMe
			F1      int `db:"f1"`
		}{EmbedMe{"XXX"}, 666},
			[]string{"EmbedMe", "f1"},
			[]any{EmbedMe{"XXX"}, 666},
			[][]string{nil, nil}},

		{struct {
			unexport1 int
			unexport2 int `db:"asd"`
		}{1, 2}, []string{}, []any{}, [][]string{}},

		{ztype.Ptr(struct {
			embedUnexport
			F1 int `db:"f1,opt"`
		}{embedUnexport{"X", 1}, 42}),
			[]string{"f1"},
			[]any{42},
			[][]string{{"opt"}}},

		{ztype.Ptr(struct {
			EmbedUnexport
			F1 int `db:"f1,opt"`
		}{EmbedUnexport{"X", 1}, 42}),
			[]string{"fe", "f1"},
			[]any{"X", 42},
			[][]string{nil, {"opt"}}},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			haveNames, haveValues, haveOpts := Fields(tt.in, "db", "skip")
			if !reflect.DeepEqual(haveNames, tt.wantNames) {
				t.Errorf("wrong names\nhave: %#v\nwant: %#v", haveNames, tt.wantNames)
			}
			if !reflect.DeepEqual(haveValues, tt.wantValues) {
				t.Errorf("wrong values\nhave: %#v\nwant: %#v", haveValues, tt.wantValues)
			}
			if !reflect.DeepEqual(haveOpts, tt.wantOpts) {
				t.Errorf("wrong options\nhave: %#v\nwant: %#v", haveOpts, tt.wantOpts)
			}

			if n := Names(tt.in, "db", "skip"); !reflect.DeepEqual(n, haveNames) {
				t.Errorf("Names() has different return\nhave:  %#v\nnames: %#v", haveNames, n)
			}
			if v := Values(tt.in, "db", "skip"); !reflect.DeepEqual(v, haveValues) {
				t.Errorf("Values() has different return\nhave:   %#v\nvalues: %#v", haveValues, v)
			}
		})
	}
}

func TestFieldsPanic(t *testing.T) {
	t.Run("not a struct", func(t *testing.T) {
		var rec any
		func() {
			defer func() { rec = recover() }()
			Fields(123, "", "")
		}()
		if rec == nil {
			t.Fatal("no panic")
		}
		if !strings.Contains(rec.(string), "not a struct") {
			t.Fatalf("wrong panic: %s", rec)
		}
	})

	t.Run("skip without tag", func(t *testing.T) {
		var rec any
		func() {
			defer func() { rec = recover() }()
			Fields(struct{}{}, "", "skip")
		}()
		if rec == nil {
			t.Fatal("no panic")
		}
		if !strings.Contains(rec.(string), "skip without tagname") {
			t.Fatalf("wrong panic: %s", rec)
		}
	})
}

var g1, g2, g3 any

type Strukt struct {
	ID                   string    `db:"log_postgres_id,noinsert"`
	LogTime              time.Time `db:"log_time"`
	UserName             *string   `db:"user_name"`
	DatabaseName         *string   `db:"database_name"`
	ProcessID            int       `db:"process_id"`
	ConnectionFrom       *string   `db:"connection_from"`
	SessionID            string    `db:"session_id"`
	SessionLineNum       int64     `db:"session_line_num"`
	CommandTag           *string   `db:"command_tag"`
	SessionStartTime     time.Time `db:"session_start_time"`
	VirtualTransactionID *string   `db:"virtual_transaction_id"`
	TransactionID        int64     `db:"transaction_id"`
	ErrorSeverity        string    `db:"error_severity"`
	SQLStateCode         string    `db:"sql_state_code"`
	Message              string    `db:"message"`
	Detail               *string   `db:"detail"`
	Hint                 *string   `db:"hint"`
	InternalQuery        *string   `db:"internal_query"`
	InternalQueryPos     *int      `db:"internal_query_pos"`
	Context              *string   `db:"context"`
	Query                *string   `db:"query"`
	QueryPos             *int      `db:"query_pos"`
	Location             *string   `db:"location"`
	ApplicationName      string    `db:"application_name"`
	BackendType          string    `db:"backend_type"`
	LeaderPID            *int      `db:"leader_pid"`
	QueryID              int64     `db:"query_id"`
}

func BenchmarkTag(b *testing.B) {
	f := reflect.TypeOf(Strukt{}).Field(0)
	var v1, v2 any
	for n := 0; n < b.N; n++ {
		v1, v2 = Tag(f, "db")
	}
	g1, g2 = v1, v2
}

func BenchmarkFields(b *testing.B) {
	var s Strukt
	var v1, v2, v3 any
	for n := 0; n < b.N; n++ {
		v1, v2, v3 = Fields(s, "db", "noinsert")
	}
	g1, g2, g3 = v1, v2, v3
}

func BenchmarkNames(b *testing.B) {
	var s Strukt
	var v1 any
	for n := 0; n < b.N; n++ {
		v1 = Names(s, "db", "noinsert")
	}
	g1 = v1
}

func BenchmarkValues(b *testing.B) {
	var s Strukt
	var v1 any
	for n := 0; n < b.N; n++ {
		v1 = Values(s, "db", "noinsert")
	}
	g1 = v1
}
