package diffator

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func panicf(msg string, args ...any) {
	panic(fmt.Sprintf(msg, args...))
}

func ContainsReflectValue(rvs []ReflectValuer, rv ReflectValuer) (contains bool) {
	for _, item := range rvs {
		if ReflectValuesEqual(item, rv) {
			contains = ReflectValuesEqual(item, rv)
			goto end
		}
	}
end:
	return contains
}

func ReflectValuesEqual(rv1, rv2 ReflectValuer) (found bool) {
	var s1, s2 string

	if rv1 == rv2 {
		found = true
		goto end
	}
	if reflect.DeepEqual(rv1, rv2) {
		found = true
		goto end
	}
	s1 = AsString(rv1)
	s2 = AsString(rv2)
	if s1 == s2 {
		found = true
		goto end
	}
end:
	return found
}

func SortReflectValues(rvs []ReflectValuer) []ReflectValuer {
	keys := make([]ReflectValuer, len(rvs))
	for i, k := range rvs {
		keys[i] = k
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i].String() < keys[j].String()
	})
	return keys
}

func SliceReduceFunc[S ~[]E, E any, R any](s S, f func(any, R) R) (r R) {
	for _, e := range s {
		r = f(e, r)
	}
	return r
}

func ReflectorsToNameString[S ~[]E, E any](in S) (names string) {
	names = SliceReduceFunc(in, func(rv any, r string) (name string) {
		switch t := rv.(type) {
		case ReflectTyper:
			name = t.Name()
		case ReflectValuer:
			name = t.ReflectType().Name()
		case reflect.Type:
			name = t.Name()
		case reflect.Value:
			name = t.Type().Name()
		default:
			panicf("ReflectorsToNameString() does not support a slice of type '%s'",
				reflect.TypeOf(t).String(),
			)
		}
		return fmt.Sprintf("%s,%s", r, name)
	})
	if names != "" {
		// Remove the leading comma
		names = names[1:]
	}
	return names
}

func AsString(a any) (s string) {
	var rv ReflectValuer
	switch t := a.(type) {
	case ReflectValuer:
		rv = t
	default:
		rv = NewDiffator().NewValue(t)
	}
	if !rv.IsValid() {
		s = "nil"
		goto end
	}
	switch rv.Kind() {
	case reflect.Interface:
		s = AsString(ChildOf(rv))
	case reflect.Pointer:
		//s = "*" + AsString(ChildOf(rv)) // TODO: Need to resolve infinite recursion
		s = fmt.Sprintf("%016x", rv.Pointer())
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
		sb.WriteString(TypenameOf(rv))
		keys := SortedMapKeys(rv)
		sb.WriteByte('{')
		raw := rv.ReflectValue()
		for _, key := range keys {
			sb.WriteString(AsString(key))
			sb.WriteByte(':')
			sb.WriteString(AsString(raw.MapIndex(key)))
			sb.WriteByte(',')
		}
		sb.WriteByte('}')
		s = sb.String()
	case reflect.Slice, reflect.Array:
		sb := strings.Builder{}
		sb.WriteString(TypenameOf(rv))
		sb.WriteByte('{')
		for i := 0; i < rv.Len(); i++ {
			sb.WriteString(AsString(rv.Index(i)))
			sb.WriteByte(',')
		}
		sb.WriteByte('}')
		s = sb.String()
	case reflect.Struct:
		sb := strings.Builder{}
		sb.WriteString(TypenameOf(rv))
		sb.WriteByte('{')
		rt := rv.Type()
		for i := 0; i < rv.NumField(); i++ {
			sb.WriteString(AsString(rt.Field(i).Name))
			sb.WriteByte(':')
			sb.WriteString(AsString(rv.Field(i)))
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

func TypenameOf(rv ReflectValuer) (n string) {
	var rt ReflectTyper

	if !rv.IsValid() {
		n = "nil"
		goto end
	}
	rt = rv.Type()
	switch rv.Kind() {
	case reflect.Interface:
		n = fmt.Sprintf("any(%s)", TypenameOf(ChildOf(rv)))
	case reflect.Pointer:
		n = fmt.Sprintf("*%s", TypenameOf(ChildOf(rv)))
	default:
		n = rt.String()
	}
end:
	return n
}

func ChildOf(rv ReflectValuer) (c ReflectValuer) {
	switch rv.Kind() {
	case reflect.Pointer, reflect.Interface:
		c = rv.Elem()
	}
	return c
}

func SortedMapKeys(a any) (keys []reflect.Value) {
	var rv reflect.Value
	switch t := a.(type) {
	case reflect.Value:
		rv = t
	case ReflectValuer:
		rv = t.ReflectValue()
	default:
		rv = reflect.ValueOf(a)
	}
	if rv.Kind() != reflect.Map {
		panicf("Value passed to SortedMapKeys() not a map: '%v'", a)
	}
	keyValues := rv.MapKeys()
	keys = make([]reflect.Value, len(keyValues))
	for i, k := range keyValues {
		keys[i] = k
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i].String() < keys[j].String()
	})
	return keys
}
