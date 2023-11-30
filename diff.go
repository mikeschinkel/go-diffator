package diffator

func Diff(v1, v2 any) string {
	d := New()
	return d.Diff(v1, v2)
}

func DiffWithFormat(v1, v2 ReflectValuer, format string) string {
	d := New()
	return d.DiffWithFormat(v1, v2, format)
}

func ReflectValuesDiff(rv1, rv2 ReflectValuer) string {
	d := New()
	return d.ReflectValuesDiff(rv1, rv2)
}

func ReflectValuesDiffWithFormat(rv1, rv2 ReflectValuer, format string) string {
	d := New()
	return d.ReflectValuesDiffWithFormat(rv1, rv2, format)
}
