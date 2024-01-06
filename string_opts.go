package diffator

import (
	"strings"
)

type StringOpts struct {
	MatchingPadLen  *IntValue
	MinSubstrLen    *IntValue
	LeftRightFormat *StringValue
}

// findInfixes finds the strings and substrings after prefixes and suffixes are
// found. It creates a down-growth tree structure where differing prefixes and
// suffixes are found and common values stored in infix property of the `node`
// struct.
func (opts *StringOpts) findInfixes(s1, s2 string) (ifx fixer) {
	var t *tree
	var pos2 int

	ss, pos1 := longestCommonSubstr(s1, s2)
	pos2 = strings.Index(s2, ss)
	switch opts.hasCommonSubstr(ss) {
	case true:
		//goland:noinspection GoAssignmentToReceiver
		t = newTree(opts)
		t.prefix = opts.findInfixes(s1[:pos1], s2[:pos2])
		t.infix.(*node).AddBoth(ss)
		t.suffix = opts.findInfixes(s1[len(ss)+pos1:], s2[len(ss)+pos2:])
		ifx = t
	case false:
		n := newNode(opts)
		n.AddLeft(s1)
		n.AddRight(s2)
		ifx = n
		goto end
	}
end:
	return ifx
}

func (opts *StringOpts) SetDefaults() {
	if opts.MinSubstrLen == nil {
		opts.MinSubstrLen = Int(MinSubstrLen)
	}
	if opts.LeftRightFormat == nil {
		opts.LeftRightFormat = String(LeftRightFormat)
	}
	if opts.MatchingPadLen == nil {
		opts.MatchingPadLen = Int(0)
	}
}

// hasCommonSubstr returns true is a "common substring" — see `const
// MinSubstrLen` for definition of common substring. Used to decide is we capture
// two strings are left vs. right or attempt to subdivide them again.
func (opts *StringOpts) hasCommonSubstr(ss string) (has bool) {
	if len(ss) <= opts.MinSubstrLen.Value {
		goto end
	}
	if ss == "" {
		goto end
	}
	if ss == " " {
		goto end
	}
	has = true
end:
	return has
}
