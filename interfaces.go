package diffator

import (
	"fmt"
	"reflect"
	"unsafe"
)

type ReflectValuer interface {
	fmt.Stringer
	Type() ReflectTyper
	NumField() int
	MapKeys() []ReflectValuer
	MapIndex(ReflectValuer) ReflectValuer
	Field(int) ReflectValuer
	Kind() reflect.Kind
	Elem() ReflectValuer
	Int() int64
	Float() float64
	Bool() bool
	Pointer() uintptr
	UnsafePointer() unsafe.Pointer
	Len() int
	Index(int) ReflectValuer
	Uint() uint64
	IsValid() bool
	IsNil() bool
	ReflectValue() reflect.Value
	ReflectType() reflect.Type
	ReflectTyper() ReflectTyper
	ValueString() string
}

type ReflectTyper interface {
	fmt.Stringer
	Name() string
	NumIn() int
	Field(int) reflect.StructField
	In(int) ReflectTyper
	NumOut() int
	Out(int) ReflectTyper
	Key() ReflectTyper
	Elem() ReflectTyper
	IsVariadic() bool
	ReflectType() reflect.Type
}
