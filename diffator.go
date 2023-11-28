package diffator

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

type Diffator struct {
	seen []reflect.Value
}

func New() *Diffator {
	return &Diffator{
		seen: make([]reflect.Value, 0),
	}
}

func (d *Diffator) Diff(v1, v2 any) string {
	return d.ReflectValuesDiffWithFormat(
		reflect.ValueOf(v1),
		reflect.ValueOf(v2),
		"%s",
	)
}

func (d *Diffator) DiffWithFormat(v1, v2 reflect.Value, format string) string {
	return d.ReflectValuesDiffWithFormat(
		reflect.ValueOf(v1),
		reflect.ValueOf(v2),
		format,
	)
}

func (d *Diffator) ReflectValuesDiff(rv1, rv2 reflect.Value) string {
	return ReflectValuesDiffWithFormat(rv1, rv2, "%s")
}

func (d *Diffator) ReflectValuesDiffWithFormat(rv1, rv2 reflect.Value, format string) (diff string) {
	var sb strings.Builder
	var alreadySeen bool

	if !d.checkValid(rv1, rv2, sb) {
		goto end
	}

	if !d.checkKind(rv1, rv2, sb) {
		goto end
	}

	if d.alreadySeen(rv1) {
		alreadySeen = true
		goto end
	}

	d.push(rv1)

	switch rv1.Kind() {
	case reflect.Pointer:
		diff := d.ReflectValuesDiffWithFormat(rv1.Elem(), rv2.Elem(), "*%s")
		if diff != "" {
			sb.WriteString(fmt.Sprintf(format, diff))
		}

	case reflect.Interface:
		diff := d.ReflectValuesDiffWithFormat(rv1.Elem(), rv2.Elem(), "any(%s)")
		if diff != "" {
			sb.WriteString(fmt.Sprintf(format, diff))
		}

	case reflect.Struct:
		diff := d.diffStruct(rv1, rv2, format)
		if len(diff) > 0 {
			diff = fmt.Sprintf("%s{%s}", rv1.Type().String(), diff)
			sb.WriteString(fmt.Sprintf(format, diff))
		}

	case reflect.Slice, reflect.Array:
		diff := d.diffElements(rv1, rv2, format)
		if len(diff) > 0 {
			switch rv1.Kind() {
			case reflect.Slice:
				diff = fmt.Sprintf("[]%s", diff)
			case reflect.Array:
				diff = fmt.Sprintf("[%d]%s", rv1.Len(), diff)
			}
			sb.WriteString(fmt.Sprintf(format, diff))
		}

	case reflect.Map:
		diff := d.diffMaps(rv1, rv2, format)
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
			diff := fmt.Sprintf(format, d.notEqualDiff(rv1.Int(), rv2.Int()))
			sb.WriteString(diff)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if rv1.Uint() != rv2.Uint() {
			diff := fmt.Sprintf(format, d.notEqualDiff(rv1.Uint(), rv2.Uint()))
			sb.WriteString(diff)
		}

	case reflect.String:
		if rv1.String() != rv2.String() {
			diff := fmt.Sprintf(format, d.notEqualDiff(rv1.String(), rv2.String()))
			sb.WriteString(diff)
		}

	case reflect.Invalid:
		if reflect.DeepEqual(rv1, rv2) {
			diff := fmt.Sprintf(format, d.notEqualDiff(rv1, rv2))
			sb.WriteString(diff)
		}

	default:
		panicf("Unhandled kind '%s'", rv1.Kind())
	}

	if sb.Len() > 0 {
		diff = sb.String()
	}

end:
	if !alreadySeen {
		d.pop()
	}
	return diff
}

func (d *Diffator) diffStruct(rv1 reflect.Value, rv2 reflect.Value, format string) string {
	diff := ""
	tmpSB := strings.Builder{}
	for i := 0; i < rv1.NumField(); i++ {
		diff = d.ReflectValuesDiffWithFormat(
			rv1.Field(i),
			rv2.Field(i),
			fmt.Sprintf("%v:%s,", rv1.Field(i).Type().Name(), "%v"),
		)
		if diff != "" {
			tmpSB.WriteString(diff)
		}
	}
	diff = tmpSB.String()
	return diff
}

func (d *Diffator) diffElements(rv1 reflect.Value, rv2 reflect.Value, format string) (diff string) {
	sb := strings.Builder{}
	for i := 0; i < rv1.Len(); i++ {
		diff = d.ReflectValuesDiffWithFormat(rv1.Index(i), rv2.Index(i), "%s")
		if diff != "" {
			sb.WriteString(diff)
		}
	}
	diff = sb.String()
	return diff
}

func (d *Diffator) diffMaps(rv1 reflect.Value, rv2 reflect.Value, format string) (diff string) {
	sb := strings.Builder{}
	keys1 := SortReflectValues(rv1.MapKeys())
	keys2 := SortReflectValues(rv2.MapKeys())
	for i, k := range keys1 {
		if !ContainsReflectValue(keys2, k) {
			sb.WriteString(fmt.Sprintf("%v:<missing:expected>,", k))
			continue
		}
		slices.DeleteFunc(keys2, func(value reflect.Value) bool {
			return ReflectValuesEqual(k, keys2[i])
		})
		diff = d.ReflectValuesDiffWithFormat(
			rv1.MapIndex(k),
			rv2.MapIndex(k),
			fmt.Sprintf("%v:%s,", k, "%v"),
		)
		if diff != "" {
			sb.WriteString(diff)
		}
	}
	for _, k := range keys2 {
		if !ContainsReflectValue(keys1, k) {
			sb.WriteString(fmt.Sprintf("%v:<missing:actual>,", k))
		}
	}
	diff = sb.String()
	return diff
}

func (d *Diffator) checkValid(rv1, rv2 reflect.Value, sb strings.Builder) bool {
	if rv1.IsValid() != rv2.IsValid() {
		sb.WriteString(d.notEqualDiff(
			fmt.Sprintf("Valid:%t", rv1.IsValid()),
			fmt.Sprintf("Valid:%t", rv2.IsValid()),
		))
		return false
	}
	return true
}

func (d *Diffator) checkKind(rv1, rv2 reflect.Value, sb strings.Builder) bool {
	if rv1.Kind() != rv2.Kind() {
		sb.WriteString(d.notEqualDiff(
			fmt.Sprintf("Kind:%s", rv1.Kind().String()),
			fmt.Sprintf("Kind:%s", rv2.Kind().String()),
		))
		return false
	}
	return true
}

func (d *Diffator) notEqualDiff(v1, v2 any) (diff string) {
	return fmt.Sprintf("(%v!=%v)", v1, v2)
}

func (d *Diffator) push(rv reflect.Value) {
	d.seen = append(d.seen, rv)
}

func (d *Diffator) pop() {
	d.seen = d.seen[:len(d.seen)-1]
}
func (d *Diffator) alreadySeen(rv reflect.Value) bool {
	return ContainsReflectValue(d.seen, rv)
}
