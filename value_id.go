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
