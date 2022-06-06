package zreflect

import (
	"reflect"
	"testing"
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
