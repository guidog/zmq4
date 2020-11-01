package zmq4

import (
	"testing"
)

func TestMsgFrames(t *testing.T) {
	rawA := []byte("Frame 0")
	rawB := []byte("Frame A")
	msg := NewMsgFrom(rawA, rawB)
	if len(msg.Frames) != 2 {
		t.Errorf("wrong number of frames, expected 2, got %d\n", len(msg.Frames))
	}
	msg.Reset()
	if len(msg.Frames) != 0 {
		t.Errorf("wrong number of frames, expected 0, got %d\n", len(msg.Frames))
	}
	if cap(msg.Frames) != 2 {
		t.Errorf("wrong cap of frames, expected 2, got %d\n", cap(msg.Frames))
	}
	// meddle w/ Frames so the unerlying []bytes can be inspected
	msg.Frames = msg.Frames[:2]
	// check len and cap of frames storage
	for i := range msg.Frames {
		fr := msg.Frames[i]
		if len(fr) != 0 {
			t.Errorf("wrong len of frame %d, expected 0, got %d\n", i, cap(fr))
		}
		if cap(fr) != 0 {
			t.Errorf("wrong cap of frame %d, expected 0, got %d\n", i, cap(fr))
		}
	}
}
