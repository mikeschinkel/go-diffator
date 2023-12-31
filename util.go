package diffator

import (
	"fmt"
	"reflect"
	"sort"
)

func panicf(msg string, args ...any) {
	panic(fmt.Sprintf(msg, args...))
}

func SortReflectValues(rvs []reflect.Value) []reflect.Value {
	keys := make([]reflect.Value, len(rvs))
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
		case reflect.Type:
			name = t.Name()
		case *reflect.Value:
			name = t.Type().Name()
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

func SortedMapKeys(a any) (keys []reflect.Value) {
	var rv reflect.Value
	switch t := a.(type) {
	case *reflect.Value:
		rv = *t
	case reflect.Value:
		rv = t
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

func isReference(rk reflect.Kind) bool {
	switch rk {
	case reflect.Pointer, reflect.Map, reflect.Slice:
		return true
	}
	return false
}

func Equivalent(v1, v2 any) (same bool) {
	rv1 := reflect.ValueOf(v1)
	rv2 := reflect.ValueOf(v2)
	switch {
	case !rv1.IsValid() && !rv2.IsValid():
		same = true
		goto end
	case !rv1.IsValid() || !rv2.IsValid():
		goto end
	}
	if rv1.Kind() != rv2.Kind() {
		goto end
	}
	if rv1.Kind() == reflect.Pointer {
		rv1 = rv1.Elem()
	}
	if rv2.Kind() == reflect.Pointer {
		rv2 = rv2.Elem()
	}
	if rv1.Comparable() && rv2.Comparable() && !rv1.Equal(rv2) {
		goto end
	}
	if !reflect.DeepEqual(rv1, rv2) {
		goto end
	}
	same = true
end:
	return same
}
