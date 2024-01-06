package diffator

// IntValue enables a nil length, so we can set defaults on nil
type IntValue struct {
	Value int
}

func Int(n int) *IntValue {
	return &IntValue{Value: n}
}

// BoolValue enables a nil length, so we can set defaults on nil
type BoolValue struct {
	Value bool
}

func Bool(b bool) *BoolValue {
	return &BoolValue{Value: b}
}

// StringValue enables a nil length so we can set defaults on nil
type StringValue struct {
	Value string
}

func String(s string) *StringValue {
	return &StringValue{Value: s}
}
