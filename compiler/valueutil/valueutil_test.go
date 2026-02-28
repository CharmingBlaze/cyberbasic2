package valueutil

import "testing"

func TestIsTruthy(t *testing.T) {
	tests := []struct {
		name string
		v    interface{}
		want bool
	}{
		{"nil", nil, false},
		{"false", false, false},
		{"true", true, true},
		{"zero int", 0, false},
		{"non-zero int", 1, true},
		{"zero float", 0.0, false},
		{"non-zero float", 1.5, true},
		{"empty string", "", false},
		{"non-empty string", "x", true},
		{"negative int", -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsTruthy(tt.v); got != tt.want {
				t.Errorf("IsTruthy(%v) = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}
