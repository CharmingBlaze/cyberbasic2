package runtime

import "testing"

func TestDetectWindowMode(t *testing.T) {
	tests := []struct {
		src  string
		want WindowMode
	}{
		{`InitWindow(800, 600, "x")`, ModeExplicit},
		{`initwindow(1,2,"a")`, ModeExplicit},
		{`ON UPDATE
END ON`, ModeImplicit},
		{`SUB OnUpdate(dt)
ENDSUB`, ModeImplicit},
		{`ON DRAW
END ON`, ModeImplicit},
		{`PRINT 1`, ModeConsole},
		{`InitWindow(1,1,"a")
ON UPDATE
END ON`, ModeExplicit},
	}
	for _, tt := range tests {
		if got := DetectWindowMode(tt.src); got != tt.want {
			t.Errorf("DetectWindowMode(%q) = %v, want %v", tt.src, got, tt.want)
		}
	}
}
