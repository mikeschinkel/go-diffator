package diffator

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
)

func panicf(msg string, args ...any) {
	panic(fmt.Sprintf(msg, args...))
}

func ContainsReflectValue(rvs []ReflectValuer, rv ReflectValuer) (contains bool) {
	for _, item := range rvs {
		if ReflectValuesEqual(item, rv) {
			contains = true
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
	s1 = fmt.Sprintf("%v", rv1)
	s2 = fmt.Sprintf("%v", rv2)
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
		// TODO: What about named interfaces?
		s = fmt.Sprintf("any(%s)", AsString(ChildOf(rv)))
	case reflect.Pointer:
		// TODO: This is probably wrong
		s = fmt.Sprintf("*%s", AsString(ChildOf(rv)))
	case reflect.String:
		s = strconv.Quote(rv.String())
	case reflect.Int, reflect.Int8, reflect.Int16:
		s = strconv.Itoa(int(rv.Int()))
	case reflect.Int32, reflect.Int64:
		s = strconv.FormatInt(rv.Int(), 10)
	case reflect.Float32:
		s = strconv.FormatFloat(rv.Float(), 'g', 10, 32)
	case reflect.Float64:
		s = strconv.FormatFloat(rv.Float(), 'g', 10, 64)
	case reflect.Map, reflect.Slice, reflect.Struct:
		s = fmt.Sprintf("%s{...}", TypenameOf(rv))
	case reflect.Bool:
		if rv.Bool() {
			s = "true"
		} else {
			s = "false"
		}
	default:
		panicf("Unhandled (s of yet) reflect value kind: %s", rv.Kind())
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
