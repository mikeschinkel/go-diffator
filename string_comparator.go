package diffator

import (
	"unicode/utf8"

	"github.com/mikeschinkel/go-lib"
)

var _ Comparator = (*StringComparator)(nil)

type StringComparator struct {
	*tree
	s1  string
	s2  string
	og1 string
	og2 string
}

func NewStringComparator(s1, s2 string, opts *StringOpts) *StringComparator {
	if opts == nil {
		opts = &StringOpts{}
	}
	opts.SetDefaults()
	return &StringComparator{
		s1:   s1,
		s2:   s2,
		tree: newTree(opts),
	}
}

func (c *StringComparator) Opts() *StringOpts {
	return c.opts
}

func (c *StringComparator) Compare() (s string) {
	var ok bool

	if s, ok = c.handleEmptyString(); !ok {
		goto end
	}
	c = c.findPrefixes()
	c = c.findSuffixes()
	c = c.findInfixes()
end:
	return c.String()
}

// findPrefixes finds the initial suffixes. This could be handled by logic in
// findInfixes, but then the logic for trimming the prefixes to pad length
// becomes much more complicated
func (c *StringComparator) findPrefixes() *StringComparator {
	n1, n2 := 0, 0
	s1 := c.s1
	s2 := c.s2
	prefix := newNode(c.opts)
	for {
		if len(s1) == 0 {
			goto end
		}
		if len(s2) == 0 {
			goto end
		}
		r1, sz1 := utf8.DecodeRuneInString(s1)
		if r1 == utf8.RuneError {
			lib.Panicf("ERROR: Attempting to retrieve last rune in '%s'", s1)
		}
		r2, sz2 := utf8.DecodeRuneInString(s2)
		if r2 == utf8.RuneError {
			lib.Panicf("ERROR: Attempting to retrieve last rune in '%s'", s2)
		}
		if r1 != r2 {
			goto end
		}
		prefix.AddBoth(string(r1))
		s1 = s1[sz1:]
		n1 += sz1
		s2 = s2[sz2:]
		n2 += sz2
	}
end:
	c.s1 = s1
	c.s2 = s2

	pad := c.opts.MatchingPadLen.Value
	// Trim the prefix if longer than the pad amount.
	if pad > 0 && len(prefix.both) > pad {
		prefix.both = prefix.both[len(prefix.both)-pad:]
	}

	c.prefix = prefix
	return c
}

// findSuffixes finds the initial suffixes. This could be handled by logic in
// findInfixes, but then the logic for trimming the suffixes to pad length
// becomes much more complicated
func (c *StringComparator) findSuffixes() *StringComparator {
	n1, n2 := 0, 0
	s1 := c.s1
	s2 := c.s2
	suffix := newNode(c.opts)
	for {
		if len(s1) == 0 {
			goto end
		}
		if len(s2) == 0 {
			goto end
		}
		r1, sz1 := utf8.DecodeLastRuneInString(s1)
		if r1 == utf8.RuneError {
			panicf("ERROR: Attempting to retrieve last rune in '%s'", s1)
		}
		r2, sz2 := utf8.DecodeLastRuneInString(s2)
		if r2 == utf8.RuneError {
			panicf("ERROR: Attempting to retrieve last rune in '%s'", s2)
		}
		if r1 != r2 {
			goto end
		}
		suffix.InsertBoth(string(r1))
		s1 = s1[:len(s1)-sz1]
		n1 += sz1
		s2 = s2[:len(s2)-sz2]
		n2 += sz2
	}
end:
	c.s1 = c.s1[:len(c.s1)-n1]
	c.s2 = c.s2[:len(c.s2)-n2]

	pad := c.opts.MatchingPadLen.Value
	// Trim the prefix if longer than the pad amount.
	if pad > 0 && len(suffix.both) > pad {
		suffix.both = suffix.both[:pad]
	}

	c.suffix = suffix
	return c
}

func (c *StringComparator) findInfixes() *StringComparator {
	ft := c.opts.findInfixes(c.s1, c.s2)
	c.infix = ft
	return c
}

// handleEmptyString upfront handles empty strings on left, right or both. This
// could be handled by logic in findInfixes, but then the logic for trimming the
// suffixes to pad length becomes much more complicated.
func (c *StringComparator) handleEmptyString() (s string, ok bool) {
	s1 := c.s1
	s2 := c.s2
	switch {
	case len(s1) == 0 && len(s2) == 0:
		goto end
	case len(s1) == 0:
		c.prefix.(*node).AddRight(s2)
		goto end
	case len(s2) == 0:
		c.prefix.(*node).AddLeft(s1)
		goto end
	default:
		ok = true
	}
end:
	return s, ok
}
