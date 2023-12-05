package diffator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Reflector struct {
	*reflect.Value
	original any
	tracker  *Tracker
}

func NewReflector(value any) *Reflector {
	var rv *reflect.Value
	switch t := value.(type) {
	case reflect.Value:
		rv = &t
	case *reflect.Value:
		rv = t
	default:
		tmp := reflect.ValueOf(value)
		rv = &tmp
	}
	return &Reflector{
		Value:    rv,
		original: value,
		tracker:  NewTracker(),
	}
}

func NewReflectorFromValue(rv *reflect.Value) *Reflector {
	return &Reflector{
		Value:   rv,
		tracker: NewTracker(),
	}
}

func (r *Reflector) String() (s string) {
	return r.AsString(r.Value)
}
func (r *Reflector) Child() *reflect.Value {
	return r.ChildOf(r.Value)
}
func (r *Reflector) Typename() string {
	return r.TypenameOf(r.Value)
}
func (r *Reflector) Any() any {
	return r.AsAny(r.Value)
}

func (r *Reflector) AsAny(rv *reflect.Value) (a any) {
	if !rv.IsValid() {
		a = nil
		goto end
	}
	if rv.CanInterface() {
		a = rv.Interface()
		goto end
	}
	switch rv.Kind() {
	case reflect.Bool:
		a = rv.Bool()
	case reflect.String:
		a = rv.String()
	case reflect.Int:
		a = int(rv.Int())
	case reflect.Int8:
		a = int8(rv.Int())
	case reflect.Int16:
		a = int16(rv.Int())
	case reflect.Int32:
		a = int32(rv.Int())
	case reflect.Int64:
		a = rv.Int()
	case reflect.Float32:
		a = float32(rv.Float())
	case reflect.Float64:
		a = rv.Float()
	case reflect.Map:
		a = map[any]any{"<example>": "<example>"}
	case reflect.Slice:
		a = []any{"<example>"}
	case reflect.UnsafePointer:
		a = "*<example>"
	case reflect.Interface, reflect.Pointer:
		a = r.AsAny(r.ChildOf(rv))
	default:
		panicf("Unhandled reflect value kind (as of yet): %s", rv.Kind())
	}
end:
	return a
}

func (r *Reflector) TypenameOf(rv *reflect.Value) (n string) {
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
	case reflect.UnsafePointer:
		n = "<unsafe-pointer>"
	default:
		n = rt.String()
	}
end:
	return n
}

func (r *Reflector) ChildOf(rv *reflect.Value) (c *reflect.Value) {
	switch rv.Kind() {
	case reflect.Pointer, reflect.Interface:
		tmp := rv.Elem()
		c = &tmp
	}
	return c
}

func (r *Reflector) AsString(rv *reflect.Value) (s string) {
	seen, id := r.tracker.Push(rv)
	if seen && isReference(rv.Kind()) {
		s = "<recursion>"
		goto end
	}
	defer r.tracker.Pop(id)
	if !rv.IsValid() {
		s = "nil"
		goto end
	}
	switch rv.Kind() {
	case reflect.Func:
		s = "func()error" // TODO: Flesh this out
	case reflect.Interface:
		s = r.AsString(r.ChildOf(rv))
	case reflect.Pointer:
		s = "*" + r.AsString(r.ChildOf(rv))
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
			sb.WriteString(r.AsString(&key))
			sb.WriteByte(':')
			idx := rv.MapIndex(key)
			sb.WriteString(r.AsString(&idx))
			sb.WriteByte(',')
		}
		sb.WriteByte('}')
		s = sb.String()
	case reflect.Slice, reflect.Array:
		sb := strings.Builder{}
		sb.WriteString(r.TypenameOf(rv))
		sb.WriteByte('{')
		for i := 0; i < rv.Len(); i++ {
			idx := rv.Index(i)
			sb.WriteString(r.AsString(&idx))
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
			sb.WriteString(rt.Field(i).Name)
			sb.WriteByte(':')
			fld := rv.Field(i)
			sb.WriteString(r.AsString(&fld))
			sb.WriteByte(',')
		}
		sb.WriteByte('}')
		s = sb.String()
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
