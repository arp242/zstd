package zreflect

import (
	"fmt"
	"testing"
)

func TestSimplify(t *testing.T) {
	{
		v, k := Simplify(int8(123))
		fmt.Printf("%v %v, %v\n", v, v.Type(), k)
		fmt.Println(v.Float())

		v, k = Simplify(int16(123))
		fmt.Printf("%v %v, %v\n", v, v.Type(), k)
		fmt.Println(v.Float())
	}
	{
		v, k := Simplify([]byte("ASD"))
		fmt.Printf("%v %v, %v\n", v, v.Type(), k)
		fmt.Println(v.String())
	}
	{
		x := "asd"
		v, k := Simplify(&x)
		fmt.Printf("%v %v, %v\n", v, v.Type(), k)
		fmt.Println(v.String())
	}

	tests := []struct {
		in interface{}
	}{
		{""},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			_ = tt

			// if have != tt.want {
			// 	t.Errorf("\nhave: %q\nwant: %q", have, tt.want)
			// }
			// if !reflect.DeepEqual(have, tt.want) {
			// 	t.Errorf("\nhave: %#v\nwant: %#v", have, tt.want)
			// }
			// if d := ztest.Diff(have, tt.want); d != "" {
			// 	t.Errorf(d)
			// }
		})
	}
}
