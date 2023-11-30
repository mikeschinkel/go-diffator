package diffator

import (
	"reflect"
)

var _ ReflectValuer = (*Value)(nil)

type Value struct {
	diffator *Diffator
	reflect.Value
	id int
}

func (v Value) ValueString() (s string) {
	return AsString(v)
}

func (v Value) ReflectTyper() ReflectTyper {
	return v.diffator.NewType(v.Value.Type())
}

func (v Value) ReflectValue() reflect.Value {
	return v.Value
}

func (v Value) ReflectType() reflect.Type {
	return v.Value.Type()
}

func (v Value) Type() ReflectTyper {
	return v.diffator.NewType(v.Value.Type())
}

func (v Value) Index(i int) ReflectValuer {
	return v.diffator.NewValue(v.Value.Index(i))
}

func (v Value) Field(i int) ReflectValuer {
	return v.diffator.NewValue(v.Value.Field(i))
}

func (v Value) Elem() ReflectValuer {
	return v.diffator.NewValue(v.Value.Elem())
}

func (v Value) MapIndex(key ReflectValuer) ReflectValuer {
	return v.diffator.NewValue(v.Value.MapIndex(key.ReflectValue()))
}

func (v Value) MapKeys() (rvs []ReflectValuer) {
	mk := v.Value.MapKeys()
	rvs = make([]ReflectValuer, len(mk))
	for i, k := range mk {
		rvs[i] = v.diffator.NewValue(k)
	}
	return rvs
}
