package inputmap

import "testing"

func TestInputPressedLogic(t *testing.T) {
	st := &actionState{wasDown: false, isDown: true}
	if st.wasDown || !st.isDown {
		// pressed = !was && is
	}
	pressed := !st.wasDown && st.isDown
	if !pressed {
		t.Fatal("expected pressed")
	}
}
