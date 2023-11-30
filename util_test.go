package diffator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsString(t *testing.T) {
	tests := []struct {
		name  string
		value any
		wantS string
	}{
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
			got := AsString(tt.value)
			assert.Equalf(t, tt.wantS, got, "AsString(%v)", tt.value)
		})
	}
}
