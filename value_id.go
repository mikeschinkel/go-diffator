package diffator

import (
	"reflect"
	"unsafe"
)

type ValueIdMap map[ValueId]struct{}

// ValueId is used to compare the source value for original pointers
type ValueId struct {
	pointer unsafe.Pointer
	reflect.Type
	altId reflect.Value
	// reference is true for pointer, map, slice
	reference bool
}

// NewValueId returns a comparable struct for any reflect type.
func NewValueId(rv *reflect.Value) (id ValueId) {
	var ptr unsafe.Pointer
	var altId reflect.Value
	var ref bool

	switch rv.Kind() {
	case reflect.Slice, reflect.Map, reflect.Pointer:
		// Use unsafe.Pointer to support the GC
		ptr = unsafe.Pointer((*rv).Pointer())
		ref = true
	default:
		altId = *rv
	}
	var rt reflect.Type
	if rv.IsValid() {
		rt = (*rv).Type()
	}
	return ValueId{
		pointer:   ptr,
		Type:      rt,
		altId:     altId,
		reference: ref,
	}
}
