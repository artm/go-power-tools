package slow_test

import (
	"bytes"
	"io"
	"slow"
	"strings"
	"testing"
	"time"
)

func TestSlowPrintReaderIn(t *testing.T) {
	t.Parallel()
	testText := "1\n2\n3"
	testDelay := 10 * time.Millisecond
	mockReader := strings.NewReader(testText)
	mockWriter := &bytes.Buffer{}
	printer := slow.NewPrinter(
		slow.WithReader(io.Reader(mockReader)),
		slow.WithWriter(io.Writer(mockWriter)),
		slow.WithDelay(testDelay),
	)
	start := time.Now()
	err := printer.Print()
	if err != nil {
		t.Error(err)
	}
	wantElapsed := time.Duration(len(testText)) * testDelay
	checkElapsed(t, start, wantElapsed, testDelay)

	want := testText
	got := mockWriter.String()
	if got != want {
		t.Errorf("wanted: %#v but got: %#v", want, got)
	}
}

func checkElapsed(
	t *testing.T,
	start time.Time,
	want time.Duration,
	accuracy time.Duration,
) {
	got := time.Since(start).Round(accuracy)
	want = want.Round(accuracy)
	if got != want {
		t.Errorf("wanted: %v but got: %v", want, got)
	}
}
