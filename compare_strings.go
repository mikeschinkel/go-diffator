package diffator

func CompareStrings(s1, s2 string, opts *StringOpts) (s string) {
	c := NewStringComparator(s1, s2, opts)
	return c.Compare()
}
