package diffator

func CompareObjects(v1, v2 any, opts *ObjectOpts) string {
	c := NewObjectComparator(v1, v2, opts)
	return c.Compare()
}
