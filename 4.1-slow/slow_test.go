package slow_test

import (
	"bytes"
	"io"
	"os"
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

func TestSlowPrintFileIn(t *testing.T) {
	t.Parallel()
	testDelay := 10 * time.Millisecond
	mockWriter := &bytes.Buffer{}
	testFilePath := "testdata/1_2_3.txt"
	fi, err := os.Stat(testFilePath)
	if err != nil {
		t.Fatal(err)
	}
	fileSize := fi.Size()
	args := []string{
		testFilePath,
	}
	printer := slow.NewPrinter(
		slow.WithArgs(args),
		slow.WithWriter(io.Writer(mockWriter)),
		slow.WithDelay(testDelay),
	)
	start := time.Now()
	err = printer.Print()
	if err != nil {
		t.Error(err)
	}
	wantElapsed := time.Duration(fileSize) * testDelay
	checkElapsed(t, start, wantElapsed, testDelay)

	bytes, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatal(err)
	}
	want := string(bytes)
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
