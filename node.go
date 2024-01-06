package diffator

import (
	"fmt"
)

var _ fixer = (*node)(nil)

type node struct {
	left  []rune
	both  []rune
	right []rune
	opts  *StringOpts
}

func (n *node) InsertBoth(s string) {
	n.both = append([]rune(s), n.both...)
}

func newNode(opts *StringOpts) *node {
	return &node{
		left:  make([]rune, 0),
		both:  make([]rune, 0),
		right: make([]rune, 0),
		opts:  opts,
	}
}

func (*node) Fixer() {}

func (n *node) AddLeft(s string) {
	n.left = append(n.left, []rune(s)...)
}

func (n *node) AddBoth(s string) {
	n.both = append(n.both, []rune(s)...)
}

func (n *node) AddRight(s string) {
	n.right = append(n.right, []rune(s)...)
}

func (n *node) bitMap() (bits int8) {
	if len(n.left) > 0 {
		bits |= 4
	}
	if len(n.both) > 0 {
		bits |= 2
	}
	if len(n.right) > 0 {
		bits |= 1
	}
	return bits
}

func (n *node) String() (s string) {
	format := n.opts.LeftRightFormat
	switch n.bitMap() {
	case 0b000:
		s = ""
	case 0b010:
		s = string(n.both)
	case 0b101, 0b100, 0b001:
		s = fmt.Sprintf(format.Value, string(n.left), string(n.right))
	case 0b111, 0b110, 0b011:
		s = fmt.Sprintf(format.Value+"%s"+format.Value,
			string(n.left), "", string(n.both), "", string(n.right),
		)
	}
	return s
}
