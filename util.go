package diffator

import (
	"fmt"
	"reflect"
	"sort"
)

func panicf(msg string, args ...any) {
	panic(fmt.Sprintf(msg, args...))
}

func ContainsReflectValue(rvs []reflect.Value, rv reflect.Value) (contains bool) {
	for _, item := range rvs {
		if ReflectValuesEqual(item, rv) {
			contains = true
			goto end
		}
	}
end:
	return contains
}

func ReflectValuesEqual(rv1, rv2 reflect.Value) (found bool) {
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
