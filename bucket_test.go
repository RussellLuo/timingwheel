package timingwheel

import "testing"

func TestBucket_Flush(t *testing.T) {
	b := newBucket()

	b.Add(&Timer{})
	b.Add(&Timer{})
	l1 := b.timers.Len()
	if l1 != 2 {
		t.Fatalf("Got (%+v) != Want (%+v)", l1, 2)
	}

	b.Flush(func(*Timer) {})
	l2 := b.timers.Len()
	if l2 != 0 {
		t.Fatalf("Got (%+v) != Want (%+v)", l2, 0)
	}
}
