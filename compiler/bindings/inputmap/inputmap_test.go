package inputmap

import "testing"

func TestInputPressedLogic(t *testing.T) {
	st := &actionState{wasDown: false, isDown: true}
	// pressed = !wasDown && isDown
	pressed := !st.wasDown && st.isDown
	if !pressed {
		t.Fatal("expected pressed")
	}
}
