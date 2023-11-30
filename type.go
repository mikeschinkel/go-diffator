package diffator

import (
	"reflect"
)

var _ ReflectTyper = (*Type)(nil)

type Type struct {
	diffator *Diffator
	reflect.Type
	id int
}

func (t Type) ReflectType() reflect.Type {
	return t.Type
}

func (t Type) In(i int) ReflectTyper {
	return t.diffator.NewType(t.Type.In(i))
}

func (t Type) Out(i int) ReflectTyper {
	return t.diffator.NewType(t.Type.Out(i))
}

func (t Type) Key() ReflectTyper {
	return t.diffator.NewType(t.Type.Key())
}

func (t Type) Elem() ReflectTyper {
	return t.diffator.NewType(t.Type.Elem())
}
