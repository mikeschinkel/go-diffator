package diffator

import (
	"reflect"
)

type Tracker struct {
	seen       ValueIdMap
	SortedKeys []reflect.Value
}

func NewTracker() *Tracker {
	return &Tracker{
		seen: make(ValueIdMap),
	}
}

func (vId *Tracker) SetSeen(seen ValueIdMap) {
	vId.seen = seen
}

func (vId *Tracker) Seen() (seen ValueIdMap) {
	return vId.seen
}

func (vId *Tracker) HaveSeen(rv *reflect.Value) (seen bool, id ValueId) {
	id = NewValueId(rv)
	return vId.HaveSeenId(id), id
}

func (vId *Tracker) HaveSeenId(id ValueId) (seen bool) {
	switch {
	case id.reference:
		_, seen = vId.seen[id]
	default:
		rv1 := id.altId
		rv1Comparable := rv1.IsValid() && rv1.Comparable()
		for v := range vId.seen {
			rv2 := v.altId
			if !rv2.IsValid() {
				continue
			}
			if rv1Comparable && rv2.Comparable() && rv1.Equal(rv2) {
				seen = true
				goto end
			}
			if reflect.DeepEqual(rv1, rv2) {
				seen = true
				goto end
			}
		}
	}
end:
	return seen
}

func (vId *Tracker) Push(rv *reflect.Value) (seen bool, id ValueId) {
	seen, id = vId.HaveSeen(rv)
	if !seen {
		vId.seen[id] = struct{}{}
		goto end
	}
	seen = true
end:
	return seen, id
}

func (vId *Tracker) Pop(id ValueId) {
	vId.Delete(id)
}
func (vId *Tracker) Delete(id ValueId) {
	delete(vId.seen, id)
}

// NewTrackerWithKeys returns sorted map keys as a slice, and a ValueIdTracker for the Value
func NewTrackerWithKeys(rv *reflect.Value) *Tracker {
	t := NewTracker()
	t.SortedKeys = SortReflectValues(rv.MapKeys())
	for _, key := range t.SortedKeys {
		rvId := NewValueId(&key)
		t.seen[rvId] = struct{}{}
	}
	return t
}
