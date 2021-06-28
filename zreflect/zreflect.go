package zreflect

import "reflect"

// Simplify a type.
//
// This reduces the Go types a limited set of values, making type switches a bit
// easier.
//
//    everything else       No conversation, return as-is.
//    bool                  reflect.Bool
//    complex*              reflect.Complex128
//    int*, uint*, float*   reflect.Float64; as floats are reliable for natural
//                          numbers up to ~9e15 (9,007,199 billion) just
//                          converting it to a float will be fine in most use
//                          cases.
//    string, []byte        reflect.String; note this also matches []uint8 and []uint32,
//            []rune        as byte and rune are just aliases for that with no way to
//                          distinguish between the two.
//
// The practical value of this is that it makes it a lot easier to deal with
// different types:
//
//     switch vv, t := Simplify(v); t {
//     case reflect.String:
//         ..
//     case reflect.Float64:
//         ..
//     case reflect.Bool:
//         ..
//     case reflect.Complex128:
//         ..
//
//     default:
//         if vv.Type() != reflect.TypeOf(time.Time{}) {
//             return vv.String()
//         }
//     }
func Simplify(value interface{}) (reflect.Value, reflect.Kind) {
	v := reflect.ValueOf(value)
	// Actually, this cases problems with some types, like time.Time
	//
	// We should probably make this easier/simpler anyway:
	//
	// 1. I want to be able to do "case time.Time{}".
	//
	// 2. "reflect.String" should catch all strings, including "type x string",
	//    but I also want to be able to "case x" specifically by putting it
	//    earlier in the switch
	//
	// 3. Also want to make it easy to switch on "does it implement this
	//    interface?" "case fmt.Stringer"
	//
	// if v.Type().Implements(reflect.TypeOf((*fmt.Stringer)(nil)).Elem()) {
	// 	v = v.MethodByName("String").Call(nil)[0]
	// 	return v, reflect.String
	// }

top:
	switch v.Kind() {
	default:
		return v, v.Kind()
	case reflect.Ptr:
		v = v.Elem()
		goto top

	case reflect.Bool:
		return v, reflect.Bool
	case reflect.Complex64, reflect.Complex128:
		return v.Convert(reflect.TypeOf(complex128(0))), reflect.Complex128
	case reflect.String:
		return v, reflect.String

	case reflect.SliceOf(reflect.TypeOf([]byte{})).Kind(): // []uint8 matches this as well.
		return v.Convert(reflect.TypeOf("")), reflect.String
	case reflect.SliceOf(reflect.TypeOf([]rune{})).Kind(): // []uint32 matches this as well.
		return v.Convert(reflect.TypeOf("")), reflect.String

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return v.Convert(reflect.TypeOf(float64(0))), reflect.Float64
	}
}
