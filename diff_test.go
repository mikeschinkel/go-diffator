package diffator_test

import (
	"testing"

	"github.com/mikeschinkel/go-diffator"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Int    int
	String string
}

func TestDiff(t *testing.T) {
	tests := []struct {
		name       string
		v1         any
		v2         any
		wantDiff   string
		wantFailed bool
	}{
		{
			name:       "int100:matching",
			v1:         100,
			v2:         100,
			wantDiff:   "",
			wantFailed: false,
		},
		{
			name:       "int100vsint99:failing",
			v1:         100,
			v2:         99,
			wantDiff:   "(100!=99)",
			wantFailed: true,
		},
		{
			name:       "struct-vs-struct:matching",
			v1:         &TestStruct{},
			v2:         &TestStruct{},
			wantDiff:   "",
			wantFailed: false,
		},
		{
			name: "struct-vs-struct:failing",
			v1:   &TestStruct{},
			v2: &TestStruct{
				Int:    1,
				String: "hello",
			},
			wantDiff:   "*diffator_test.TestStruct{Int:(0!=1),String:(!=hello),}",
			wantFailed: true,
		},
		{
			name:       "map-vs-map:matching",
			v1:         map[string]int{"Foo": 1, "Bar": 2, "Baz": 3},
			v2:         map[string]int{"Foo": 1, "Bar": 2, "Baz": 3},
			wantDiff:   "",
			wantFailed: false,
		},
		{
			name:       "map-vs-map:failing",
			v1:         map[string]int{"Foo": 1, "Bar": 2, "Baz": 3, "Superman": 0},
			v2:         map[string]int{"Foo": 10, "Bar": 20, "Baz": 30, "Batman": 0},
			wantDiff:   "map[string]int{Bar:(2!=20),Baz:(3!=30),Foo:(1!=10),Superman:<missing:expected>,Batman:<missing:actual>,}",
			wantFailed: true,
		},
		{
			name:       "slice-vs-slice:matching",
			v1:         []int{1, 2, 3},
			v2:         []int{1, 2, 3},
			wantDiff:   "",
			wantFailed: false,
		},
		{
			name:       "slice-vs-slice:failing",
			v1:         []int{3, 4, 5},
			v2:         []int{3, 4, 6, 7},
			wantDiff:   "[]int{[2](5!=6),[3](<missing>!=7),}",
			wantFailed: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDiff := diffator.Diff(tt.v1, tt.v2)
			assert.Equalf(t, tt.wantDiff, gotDiff, "Diff(%v, %v)", tt.v1, tt.v2)
			assert.Equalf(t, tt.wantFailed, gotDiff != "", "Diff(%v, %v)", tt.v1, tt.v2)
		})
	}
}
