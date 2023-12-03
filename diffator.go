package diffator

import (
	"fmt"
	"reflect"
	"strings"
)

type Diffator struct {
	values       [2]any
	seen         []reflect.Value
	CompareFuncs bool
	Pretty       bool
	Indent       string
	level        int
	FormatFunc   func(reflect.Type, any) string
	idTracker    *ValueIdTracker
}

func NewDiffator() *Diffator {
	return &Diffator{
		idTracker: NewValueIdTracker(),
		seen:      make([]reflect.Value, 0),
		Indent:    "  ",
	}
}

func (d *Diffator) Diff(v1, v2 any) string {
	return d.DiffWithFormat(v1, v2, "%s")
}

func (d *Diffator) DiffWithFormat(v1, v2 any, format string) string {
	switch v1.(type) {
	case reflect.Value:
		// We got what we need, do nothing
	default:
		v1 = reflect.ValueOf(v1)
	}
	switch v2.(type) {
	case reflect.Value:
		// We got what we need, do nothing
	default:
		v2 = reflect.ValueOf(v2)
	}
	// Copy original values so they are available during debugging, or if needed
	// later for other things.
	d.values[0] = v1
	d.values[1] = v2
	return d.ReflectValuesDiffWithFormat(
		v1.(reflect.Value),
		v2.(reflect.Value),
		format,
	)
}

func (d *Diffator) ReflectValuesDiff(rv1, rv2 reflect.Value) string {
	return ReflectValuesDiffWithFormat(rv1, rv2, "%s")
}

func (d *Diffator) ReflectValuesDiffWithFormat(rv1, rv2 reflect.Value, format string) (diff string) {
	var sb strings.Builder
	var alreadySeen bool
	var id ValueId

	if !d.checkValid(rv1, rv2, sb) {
		goto end
	}

	if !d.checkKind(rv1, rv2, sb) {
		goto end
	}

	alreadySeen, id = d.idTracker.Push(rv1)
	if alreadySeen {
		goto end
	}
	defer d.idTracker.Pop(id)

	switch rv1.Kind() {
	case reflect.Pointer:
		rv1 = rv1.Elem()
		rv2 = rv2.Elem()
		switch {
		case !rv1.IsValid() && !rv2.IsValid():
			// Do nothing
		case !rv1.IsValid():
			r := NewReflector()
			diff = d.notEqualDiff(rv2.Type(), "nil", r.AsString(rv2))
		case !rv2.IsValid():
			r := NewReflector()
			diff = d.notEqualDiff(rv1.Type(), r.AsString(rv1), "nil")
		default:
			diff = d.ReflectValuesDiffWithFormat(rv1, rv2, "*%s")
		}
		if diff != "" {
			sb.WriteString(fmt.Sprintf(format, diff))
		}

	case reflect.Interface:
		diff := d.ReflectValuesDiffWithFormat(rv1.Elem(), rv2.Elem(), "any(%s)")
		if diff != "" {
			sb.WriteString(fmt.Sprintf(format, diff))
		}

	case reflect.Struct:
		diff := d.diffStruct(rv1, rv2)
		if len(diff) > 0 {
			f := "%s{%s}"
			if d.Pretty {
				r := strings.Repeat(d.Indent, d.level)
				f = "%s{\n%s" + r + "}"
			}
			diff = fmt.Sprintf(f, rv1.Type().String(), diff)
			sb.WriteString(fmt.Sprintf(format, diff))
		}

	case reflect.Slice, reflect.Array:
		diff := d.diffElements(rv1, rv2)
		if len(diff) > 0 {
			switch rv1.Kind() {
			case reflect.Slice:
				diff = fmt.Sprintf("%s{%s}", rv1.Type().String(), diff)
			case reflect.Array:
				// TODO: Handle mismatches lengths
				diff = fmt.Sprintf("%s%s", rv1.Type().String(), diff)
			}
			sb.WriteString(fmt.Sprintf(format, diff))
		}

	case reflect.Map:
		diff := d.diffMaps(rv1, rv2)
		if len(diff) > 0 {
			diff = fmt.Sprintf("map[%s]%s{%s}",
				rv1.Type().Key(),
				rv1.Type().Elem(),
				diff,
			)
			sb.WriteString(fmt.Sprintf(format, diff))
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if rv1.Int() != rv2.Int() {
			diff := fmt.Sprintf(format, d.notEqualDiff(
				rv1.Type(),
				rv1.Int(),
				rv2.Int(),
			))
			sb.WriteString(diff)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if rv1.Uint() != rv2.Uint() {
			diff := fmt.Sprintf(format, d.notEqualDiff(
				rv1.Type(),
				rv1.Uint(),
				rv2.Uint(),
			))
			sb.WriteString(diff)
		}

	case reflect.Func:
		diff := d.diffFuncs(rv1, rv2)
		if len(diff) > 0 {
			diff = fmt.Sprintf("func(%s)%s{%s}",
				d.funcParams(rv1),
				d.funcReturns(rv1),
				diff,
			)
			sb.WriteString(fmt.Sprintf(format, diff))
		}

	case reflect.String:
		if rv1.String() != rv2.String() {
			diff := fmt.Sprintf(format, d.notEqualDiff(
				rv1.Type(),
				rv1.String(),
				rv2.String(),
			))
			sb.WriteString(diff)
		}

	case reflect.Bool:
		if rv1.Bool() != rv2.Bool() {
			diff := fmt.Sprintf(format, d.notEqualDiff(
				rv1.Type(),
				rv1.String(),
				rv2.String(),
			))
			sb.WriteString(diff)
		}

	case reflect.Float32, reflect.Float64:
		if rv1.Float() != rv2.Float() {
			diff := fmt.Sprintf(format, d.notEqualDiff(
				rv1.Type(),
				rv1.String(),
				rv2.String(),
			))
			sb.WriteString(diff)
		}

	case reflect.Invalid, reflect.UnsafePointer:
		if !reflect.DeepEqual(rv1, rv2) {
			diff := fmt.Sprintf(format, d.notEqualDiff(
				reflect.TypeOf(nil),
				rv1,
				rv2,
			))
			sb.WriteString(diff)
		}

	default:
		panicf("Unhandled kind '%s'", rv1.Kind())
	}

	if sb.Len() > 0 {
		diff = sb.String()
	}

end:
	if d.Pretty && d.level == 0 && diff != "" {
		diff = "\n" + diff
	}
	return diff
}

func (d *Diffator) diffStruct(rv1 reflect.Value, rv2 reflect.Value) string {
	d.level++
	diff := ""
	sb := strings.Builder{}
	for i := 0; i < rv1.NumField(); i++ {
		diff = d.ReflectValuesDiffWithFormat(
			rv1.Field(i),
			rv2.Field(i),
			fmt.Sprintf("%v:%s,", rv1.Type().Field(i).Name, "%v"),
		)
		if diff == "" {
			continue
		}
		if !d.Pretty {
			sb.WriteString(diff)
			continue
		}
		sb.WriteString(strings.Repeat(d.Indent, d.level))
		sb.WriteString(diff)
		sb.WriteByte('\n')
	}
	diff = sb.String()
	d.level--
	return diff
}

func (d *Diffator) indent() string {
	return strings.Repeat(d.Indent, d.level)
}

func (d *Diffator) diffElements(rv1, rv2 reflect.Value) (diff string) {
	d.level++
	sb := strings.Builder{}
	cnt := max(rv1.Len(), rv2.Len())
	for i := 0; i < cnt; i++ {
		switch {
		case i >= rv1.Len():
			diff = d.DiffWithFormat(
				"<missing>",
				fmt.Sprintf("%v", NewReflector().AsString(rv2.Index(i))),
				"%s,",
			)
		case i >= rv2.Len():
			diff = d.DiffWithFormat(
				fmt.Sprintf("%v", NewReflector().AsString(rv1.Index(i))),
				"<missing>",
				"%s,",
			)
		default:
			diff = d.ReflectValuesDiffWithFormat(rv1.Index(i), rv2.Index(i), "%s,")
		}
		if diff != "" {
			f := "[%d]%s"
			if d.Pretty {
				f = "\n" + d.indent() + f
			}
			sb.WriteString(fmt.Sprintf(f, i, diff))
		}
	}
	d.level--
	if diff != "" && d.Pretty {
		sb.WriteByte('\n')
		sb.WriteString(d.indent())
	}
	diff = sb.String()
	return diff
}

func (d *Diffator) diffMaps(rv1, rv2 reflect.Value) (diff string) {
	sb := strings.Builder{}
	keys1, nId1 := d.prepMapKeysForValue(rv1)
	keys2, nId2 := d.prepMapKeysForValue(rv2)

	for _, key := range keys1 {
		seen, id := nId2.HaveSeen(key)
		if !seen {
			sb.WriteString(fmt.Sprintf("%v:<missing:expected>,", key))
			continue
		}
		nId2.Delete(id)
		diff = d.ReflectValuesDiffWithFormat(
			rv1.MapIndex(key),
			rv2.MapIndex(key),
			fmt.Sprintf("%v:%s,", key, "%v"),
		)
		if diff != "" {
			sb.WriteString(diff)
		}
	}
	for _, key := range keys2 {
		seen, _ := nId1.HaveSeen(key)
		if !seen {
			sb.WriteString(fmt.Sprintf("%v:<missing:actual>,", key))
		}
	}
	diff = sb.String()
	return diff
}

// prepMapKeysForValue returns sorted map keys as a slice, and a ValueIdTracker for the Value
func (d *Diffator) prepMapKeysForValue(rv reflect.Value) (keys []reflect.Value, vId ValueIdTracker) {
	m := make(ValueIdMap)
	keys = rv.MapKeys()
	nId := NewValueIdTracker()
	for _, k := range keys {
		rvId := nId.IdOf(k)
		m[rvId] = struct{}{}
	}
	vId.SetSeen(m)
	keys = SortReflectValues(keys)
	return keys, vId
}

func (d *Diffator) checkValid(rv1, rv2 reflect.Value, sb strings.Builder) bool {
	if rv1.IsValid() != rv2.IsValid() {
		sb.WriteString(d.notEqualDiff(reflect.TypeOf(nil),
			fmt.Sprintf("Valid:%t", rv1.IsValid()),
			fmt.Sprintf("Valid:%t", rv2.IsValid()),
		))
		return false
	}
	return true
}

func (d *Diffator) checkKind(rv1, rv2 reflect.Value, sb strings.Builder) bool {
	if rv1.Kind() != rv2.Kind() {
		sb.WriteString(d.notEqualDiff(rv1.Type(),
			fmt.Sprintf("Kind:%s", rv1.Kind().String()),
			fmt.Sprintf("Kind:%s", rv2.Kind().String()),
		))
		return false
	}
	return true
}

func (d *Diffator) notEqualDiff(rt reflect.Type, v1, v2 any) (diff string) {
	if d.FormatFunc == nil {
		diff = fmt.Sprintf("(%v!=%v)", v1, v2)
		goto end
	}
	diff = fmt.Sprintf("(%s!=%s)",
		d.FormatFunc(rt, v1),
		d.FormatFunc(rt, v2),
	)
end:
	return diff
}

func (d *Diffator) diffFuncs(rv1 reflect.Value, rv2 reflect.Value) (diff string) {
	if rv1.IsNil() && rv2.IsNil() {
		goto end
	}
	if rv1.IsNil() {
		diff = fmt.Sprintf("(nil!=func(%s)%s)", d.funcParams(rv2), d.funcReturns(rv2))
		goto end
	}
	if rv2.IsNil() {
		diff = fmt.Sprintf("(func(%s)%s)!=nil)", d.funcParams(rv1), d.funcReturns(rv1))
		goto end
	}
	if !d.CompareFuncs {
		goto end
	}
end:
	return diff
}

func (d *Diffator) funcParams(rv reflect.Value) string {
	rt := rv.Type()
	cnt := rt.NumIn()
	last := cnt - 1
	in := make([]reflect.Type, cnt)
	for i := 0; i < cnt; i++ {
		in[i] = reflect.TypeOf(rt.In(i))
		if !rt.IsVariadic() {
			continue
		}
		if i < last {
			continue
		}
		// A variadic
		in[i] = in[i].Elem()
	}
	return ReflectorsToNameString(in)
}

func (d *Diffator) funcReturns(rv reflect.Value) (returns string) {
	rt := rv.Type()
	cnt := rt.NumOut()
	out := make([]reflect.Value, cnt)
	for i := 0; i < cnt; i++ {
		out[i] = reflect.ValueOf(rt.Out(i))
	}
	return ReflectorsToNameString(out)
}
