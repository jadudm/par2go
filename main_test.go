package par2go

import "testing"

func TestHello(t *testing.T) {
	want := "Hi!"
	if got := Hello(); got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}
