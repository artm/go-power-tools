package text_clock_test

import (
	"bytes"
	"testing"
	text_clock "text-clock"
	"time"
)

func TestPrintsTime(t *testing.T) {
	fakeOut := &bytes.Buffer{}
	when, _ := time.Parse("15:04", "23:08")
	text_clock.PrintTime(when, fakeOut)
	got := fakeOut.String()
	want := "It's 8 minutes past 23\n"
	if got != want {
		t.Errorf(`want: "%s", got: "%s"`, want, got)
	}
}
