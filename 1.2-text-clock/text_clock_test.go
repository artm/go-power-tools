package text_clock_test

import (
	"bytes"
	"io"
	"testing"
	text_clock "text-clock"
	"time"
)

func TestPrinter(t *testing.T) {
	mockWriter := &bytes.Buffer{}
	when, _ := time.Parse("15:04", "23:08")

	printer := text_clock.Printer{
		Writer: io.Writer(mockWriter),
	}

	printer.Print(when)
	got := mockWriter.String()
	want := "It's 8 minutes past 23\n"
	if got != want {
		t.Errorf(`want: "%s", got: "%s"`, want, got)
	}
}
