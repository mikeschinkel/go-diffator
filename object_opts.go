package diffator

import (
	"reflect"
)

type ObjectOpts struct {
	LevelIndent  *StringValue
	OutputFormat *StringValue
	PrettyPrint  *BoolValue
	CompareFuncs bool
	FormatFunc   func(reflect.Type, any) string
}

func (opts *ObjectOpts) SetDefaults() {
	if opts.LevelIndent == nil {
		opts.LevelIndent = String("  ")
	}
	if opts.OutputFormat == nil {
		opts.OutputFormat = String("%s")
	}
	if opts.PrettyPrint == nil {
		opts.PrettyPrint = Bool(false)
	}
}
