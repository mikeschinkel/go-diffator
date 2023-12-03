package diffator_test

import (
	"testing"

	"github.com/mikeschinkel/go-diffator"
	"github.com/stretchr/testify/assert"
)

type Recur1Struct struct {
	Name  string
	Child *Recur1Struct
}
type Recur2Struct struct {
	Name  string
	Slice []*Recur2Struct
}

func TestAsString(t *testing.T) {
	rs1 := Recur1Struct{
		Name: "Test",
	}
	rs1.Child = &rs1
	rs2 := Recur2Struct{
		Name:  "Test",
		Slice: make([]*Recur2Struct, 1),
	}
	rs2.Slice[0] = &rs2
	rs1.Child = &rs1
	tests := []struct {
		name  string
		value any
		wantS string
	}{
		{
			name:  `Direct recursion`,
			value: &rs1,
			wantS: `*diffator_test.Recur1Struct{"Name":"Test","Child":<recursion>,}`,
		},
		{
			name:  `Indirect recursion`,
			value: &rs2,
			wantS: `*diffator_test.Recur2Struct{"Name":"Test","Slice":[]*diffator_test.Recur2Struct{<recursion>,},}`,
		},
		{
			name:  "Int Seven (7)",
			value: 7,
			wantS: "7",
		},
		{
			name:  `String "Hello"`,
			value: "Hello",
			wantS: `"Hello"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := diffator.NewReflector().AsString(tt.value)
			assert.Equalf(t, tt.wantS, got, "AsString(%v)", tt.value)
		})
	}
}
