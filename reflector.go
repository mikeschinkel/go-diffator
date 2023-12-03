package diffator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Reflector struct {
	ValueIdTracker
}

func NewReflector() *Reflector {
	return &Reflector{
		ValueIdTracker: *NewValueIdTracker(),
	}
}

func (r *Reflector) AsString(a any) (s string) {
	var rv reflect.Value
	switch t := a.(type) {
	case reflect.Value:
		rv = t
	default:
		rv = reflect.ValueOf(t)
	}
	seen, id := r.Push(rv)
	if seen {
		s = "<recursion>"
		goto end
	}
	defer r.Pop(id)
	if !rv.IsValid() {
		s = "nil"
		goto end
	}
	switch rv.Kind() {
	case reflect.Interface:
		s = r.AsString(r.ChildOf(rv))
	case reflect.Pointer:
		s = "*" + r.AsString(r.ChildOf(rv)) // TODO: Need to resolve infinite recursion
	case reflect.String:
		s = strconv.Quote(rv.String())
	case reflect.Int, reflect.Int8, reflect.Int16:
		s = strconv.Itoa(int(rv.Int()))
	case reflect.Int32, reflect.Int64:
		s = strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16:
		s = strconv.Itoa(int(rv.Uint()))
	case reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		s = strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float32:
		s = strconv.FormatFloat(rv.Float(), 'g', 10, 32)
	case reflect.Float64:
		s = strconv.FormatFloat(rv.Float(), 'g', 10, 64)
	case reflect.Map:
		sb := strings.Builder{}
		sb.WriteString(r.TypenameOf(rv))
		keys := SortedMapKeys(rv)
		sb.WriteByte('{')
		for _, key := range keys {
			sb.WriteString(r.AsString(key))
			sb.WriteByte(':')
			sb.WriteString(r.AsString(rv.MapIndex(key)))
			sb.WriteByte(',')
		}
		sb.WriteByte('}')
		s = sb.String()
	case reflect.Slice, reflect.Array:
		sb := strings.Builder{}
		sb.WriteString(r.TypenameOf(rv))
		sb.WriteByte('{')
		for i := 0; i < rv.Len(); i++ {
			sb.WriteString(r.AsString(rv.Index(i)))
			sb.WriteByte(',')
		}
		sb.WriteByte('}')
		s = sb.String()
	case reflect.Struct:
		sb := strings.Builder{}
		sb.WriteString(r.TypenameOf(rv))
		sb.WriteByte('{')
		rt := rv.Type()
		for i := 0; i < rv.NumField(); i++ {
			sb.WriteString(r.AsString(rt.Field(i).Name))
			sb.WriteByte(':')
			sb.WriteString(r.AsString(rv.Field(i)))
			sb.WriteByte(',')
		}
		sb.WriteByte('}')
		s = sb.String()
	case reflect.Func:
		s = "func()error" // TODO: Flesh this out
	case reflect.UnsafePointer:
		s = fmt.Sprintf("%p", rv.UnsafePointer())
	case reflect.Bool:
		if rv.Bool() {
			s = "true"
		} else {
			s = "false"
		}
	default:
		panicf("Unhandled (as of yet) reflect value kind: %s", rv.Kind())
	}
end:
	return s
}

func (r *Reflector) TypenameOf(rv reflect.Value) (n string) {
	var rt reflect.Type

	if !rv.IsValid() {
		n = "nil"
		goto end
	}
	rt = rv.Type()
	switch rv.Kind() {
	case reflect.Interface:
		n = fmt.Sprintf("any(%s)", r.TypenameOf(r.ChildOf(rv)))
	case reflect.Pointer:
		n = fmt.Sprintf("*%s", r.TypenameOf(r.ChildOf(rv)))
	default:
		n = rt.String()
	}
end:
	return n
}

func (r *Reflector) ChildOf(rv reflect.Value) (c reflect.Value) {
	switch rv.Kind() {
	case reflect.Pointer, reflect.Interface:
		c = rv.Elem()
	}
	return c
}
