package diffator

import (
	"reflect"
	"unsafe"
)

type ValueIdTracker struct {
	seen ValueIdMap
}

func NewValueIdTracker() *ValueIdTracker {
	return NewValueIdWithSeen(make(ValueIdMap))
}

func NewValueIdWithSeen(seen ValueIdMap) *ValueIdTracker {
	return &ValueIdTracker{
		seen: seen,
	}
}

func (vId *ValueIdTracker) SetSeen(seen ValueIdMap) {
	vId.seen = seen
}

func (vId *ValueIdTracker) Seen() (seen ValueIdMap) {
	return vId.seen
}

// IdOf returns a comparable struct for any reflect type.
func (vId *ValueIdTracker) IdOf(rv reflect.Value) (id ValueId) {
	var ptr unsafe.Pointer
	var altId reflect.Value
	var ref bool

	switch rv.Kind() {
	case reflect.Slice, reflect.Map, reflect.Pointer:
		// Use unsafe.Pointer to support the GC
		ptr = unsafe.Pointer(rv.Pointer())
		ref = true
	default:
		altId = rv
	}
	var rt reflect.Type
	if rv.IsValid() {
		rt = rv.Type()
	}
	return ValueId{
		pointer:   ptr,
		Type:      rt,
		altId:     altId,
		reference: ref,
	}
}

func (vId *ValueIdTracker) HaveSeen(rv reflect.Value) (seen bool, id ValueId) {
	id = vId.IdOf(rv)
	return vId.HaveSeenId(id), id
}

func (vId *ValueIdTracker) HaveSeenId(id ValueId) (seen bool) {
	switch {
	case id.reference:
		_, seen = vId.seen[id]
	default:
		idCanInterface := id.altId.IsValid() && id.altId.CanInterface()
		for v := range vId.seen {
			if reflect.DeepEqual(id.altId, v.altId) {
				seen = true
				goto end
			}
			if !idCanInterface {
				continue
			}
			if !v.altId.IsValid() {
				continue
			}
			if !v.altId.CanInterface() {
				continue
			}
			if id.altId.Interface() != v.altId.Interface() {
				// TODO: Verify this handles every case except Pointer, Map, Slice
				continue
			}
			seen = true
			goto end
		}
	}
end:
	return seen
}

func (vId *ValueIdTracker) Push(rv reflect.Value) (seen bool, id ValueId) {
	seen, id = vId.HaveSeen(rv)
	if !seen {
		vId.seen[id] = struct{}{}
		goto end
	}
	seen = true
end:
	return seen, id
}

func (vId *ValueIdTracker) Pop(id ValueId) {
	vId.Delete(id)
}
func (vId *ValueIdTracker) Delete(id ValueId) {
	delete(vId.seen, id)
}
