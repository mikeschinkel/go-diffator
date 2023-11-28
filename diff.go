package diffator

import (
	"reflect"
)

func Diff(v1, v2 any) string {
	d := New()
	return d.Diff(v1, v2)
}

func DiffWithFormat(v1, v2 reflect.Value, format string) string {
	d := New()
	return d.DiffWithFormat(v1, v2, format)
}

func ReflectValuesDiff(rv1, rv2 reflect.Value) string {
	d := New()
	return d.ReflectValuesDiff(rv1, rv2)
}

func ReflectValuesDiffWithFormat(rv1, rv2 reflect.Value, format string) string {
	d := New()
	return d.ReflectValuesDiffWithFormat(rv1, rv2, format)
}
