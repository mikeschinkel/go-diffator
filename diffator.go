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
	tracker      *Tracker
}

func NewDiffator() *Diffator {
	return &Diffator{
		tracker: NewTracker(),
		Indent:  "  ",
	}
}

func (d *Diffator) Diff(v1, v2 any) string {
	return d.DiffWithFormat(v1, v2, "%s")
}

func (d *Diffator) DiffWithFormat(v1, v2 any, format string) string {
	switch v1.(type) {
	case reflect.Value, *reflect.Value:
		// We got what we need, do nothing
	default:
		v1 = reflect.ValueOf(v1)
	}
	switch v2.(type) {
	case reflect.Value, *reflect.Value:
		// We got what we need, do nothing
	default:
		v2 = reflect.ValueOf(v2)
	}
	// Copy original values to ensure they are available during debugging, or if
	// needed later for other things.
	d.values[0] = v1
	d.values[1] = v2
	rv1 := v1.(reflect.Value)
	rv2 := v2.(reflect.Value)
	return d.ReflectValuesDiffWithFormat(&rv1, &rv2, format)
}

func (d *Diffator) ReflectValuesDiff(rv1, rv2 *reflect.Value) string {
	return ReflectValuesDiffWithFormat(rv1, rv2, "%s")
}

func (d *Diffator) ReflectValuesDiffWithFormat(rv1, rv2 *reflect.Value, format string) (diff string) {
	var sb strings.Builder
	var alreadySeen bool
	var id ValueId

	if !d.checkValid(rv1, rv2, sb) {
		diff = "<invalid>"
		goto end
	}

	if !d.checkKind(rv1, rv2, sb) {
		diff = fmt.Sprintf("<type-mismatch>:%s", d.notEqualDiff(
			rv1.Type(),
			NewReflector(rv1).String(),
			NewReflector(rv2).String(),
		))
		goto end
	}

	alreadySeen, id = d.tracker.Push(rv1)
	if alreadySeen && isReference(rv1.Kind()) {
		goto end
	}
	defer d.tracker.Pop(id)

	switch rv1.Kind() {
	case reflect.Pointer, reflect.Interface:
		elem1 := rv1.Elem()
		elem2 := rv2.Elem()
		switch {
		case !elem1.IsValid() && !elem2.IsValid():
			// Do nothing
		case !elem1.IsValid():
			r := NewReflector(elem2)
			diff = d.notEqualDiff(elem2.Type(), "nil", r.String())
		case !elem2.IsValid():
			r := NewReflector(elem1)
			diff = d.notEqualDiff(elem1.Type(), r.String(), "nil")
		default:
			switch rv1.Kind() {
			case reflect.Pointer:
				diff = d.ReflectValuesDiffWithFormat(&elem1, &elem2, "*%s")
			case reflect.Interface:
				diff = d.ReflectValuesDiffWithFormat(&elem1, &elem2, "any(%s)")
			}
		}
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

	case reflect.UnsafePointer:
		// Do nothing, we cannot compare anyway, nor should we

	case reflect.Invalid:
		if !reflect.DeepEqual(rv1, rv2) {
			diff := fmt.Sprintf(format, d.notEqualDiff(nil, rv1, rv2))
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

func (d *Diffator) diffStruct(rv1, rv2 *reflect.Value) string {
	d.level++
	diff := ""
	sb := strings.Builder{}
	for i := 0; i < rv1.NumField(); i++ {
		fld1 := rv1.Field(i)
		fld2 := rv2.Field(i)
		name := rv1.Type().Field(i).Name
		diff = d.ReflectValuesDiffWithFormat(&fld1, &fld2, fmt.Sprintf("%v:%s,", name, "%v"))
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

func (d *Diffator) diffElements(rv1, rv2 *reflect.Value) (diff string) {
	d.level++
	sb := strings.Builder{}
	cnt := max(rv1.Len(), rv2.Len())
	for i := 0; i < cnt; i++ {
		switch {
		case i >= rv1.Len():
			idx := rv2.Index(i)
			diff = d.DiffWithFormat(
				"<missing>",
				fmt.Sprintf("%v", NewReflector(&idx).String()),
				"%s,",
			)
		case i >= rv2.Len():
			idx := rv1.Index(i)
			diff = d.DiffWithFormat(
				fmt.Sprintf("%v", NewReflector(&idx).String()),
				"<missing>",
				"%s,",
			)
		default:
			idx1 := rv1.Index(i)
			idx2 := rv2.Index(i)
			diff = d.ReflectValuesDiffWithFormat(&idx1, &idx2, "%s,")
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

func (d *Diffator) diffMaps(rv1, rv2 *reflect.Value) (diff string) {
	sb := strings.Builder{}
	tkr1 := NewTrackerWithKeys(rv1)
	tkr2 := NewTrackerWithKeys(rv2)

	for _, key := range tkr1.SortedKeys {
		seen, id := tkr2.HaveSeen(&key)
		if !seen {
			sb.WriteString(fmt.Sprintf("%v:<missing:expected>,", key))
			continue
		}
		tkr2.Delete(id)
		key1 := rv1.MapIndex(key)
		key2 := rv2.MapIndex(key)
		diff = d.ReflectValuesDiffWithFormat(&key1, &key2, fmt.Sprintf("%v:%s,", key, "%v"))
		if diff != "" {
			sb.WriteString(diff)
		}
	}
	for _, key := range tkr2.SortedKeys {
		seen, _ := tkr1.HaveSeen(&key)
		if !seen {
			sb.WriteString(fmt.Sprintf("%v:<missing:actual>,", key))
		}
	}
	diff = sb.String()
	return diff
}

func (d *Diffator) checkValid(rv1, rv2 *reflect.Value, sb strings.Builder) bool {
	if rv1.IsValid() != rv2.IsValid() {
		sb.WriteString(d.notEqualDiff(reflect.TypeOf(nil),
			fmt.Sprintf("Valid:%t", rv1.IsValid()),
			fmt.Sprintf("Valid:%t", rv2.IsValid()),
		))
		return false
	}
	return true
}

func (d *Diffator) checkKind(rv1, rv2 *reflect.Value, sb strings.Builder) bool {
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

func (d *Diffator) diffFuncs(rv1, rv2 *reflect.Value) (diff string) {
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

func (d *Diffator) funcParams(rv *reflect.Value) string {
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

func (d *Diffator) funcReturns(rv *reflect.Value) (returns string) {
	rt := rv.Type()
	cnt := rt.NumOut()
	out := make([]reflect.Value, cnt)
	for i := 0; i < cnt; i++ {
		out[i] = reflect.ValueOf(rt.Out(i))
	}
	return ReflectorsToNameString(out)
}
