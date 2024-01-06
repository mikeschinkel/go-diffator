package diffator

var _ fixer = (*tree)(nil)

type tree struct {
	prefix fixer // Used for both or left & right
	infix  fixer // Used for both (matching string) or subtree
	suffix fixer // Used for both or left & right
	opts   *StringOpts
}

func (*tree) Fixer() {}

func (t *tree) String() string {
	return t.prefix.String() +
		t.infix.String() +
		t.suffix.String()
}

func newTree(opts *StringOpts) *tree {
	return &tree{
		prefix: newNode(opts),
		infix:  newNode(opts),
		suffix: newNode(opts),
		opts:   opts,
	}
}
