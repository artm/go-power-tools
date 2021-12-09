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
	printer, err := slow.NewPrinter(
		slow.WithReader(io.Reader(mockReader)),
		slow.WithArgs([]string{}),
		slow.WithWriter(io.Writer(mockWriter)),
		slow.WithDelay(testDelay),
	)
	if err != nil {
		t.Error(err)
	}
	start := time.Now()
	err = printer.Print()
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

func TestSlowPrintFile(t *testing.T) {
	t.Parallel()
	ignoredText := "foo"
	ignoredReader := strings.NewReader(ignoredText)
	testDelay := 10 * time.Millisecond
	mockWriter := &bytes.Buffer{}
	testFilePath1 := "testdata/1_2_3.txt"
	fi, err := os.Stat(testFilePath1)
	if err != nil {
		t.Fatal(err)
	}
	totalSize := fi.Size()
	args := []string{
		testFilePath1,
	}
	printer, err := slow.NewPrinter(
		slow.WithReader(ignoredReader),
		slow.WithArgs(args),
		slow.WithWriter(io.Writer(mockWriter)),
		slow.WithDelay(testDelay),
	)
	if err != nil {
		t.Error(err)
	}
	start := time.Now()
	err = printer.Print()
	if err != nil {
		t.Error(err)
	}
	wantElapsed := time.Duration(totalSize) * testDelay
	checkElapsed(t, start, wantElapsed, testDelay)

	bytes, err := os.ReadFile(testFilePath1)
	if err != nil {
		t.Fatal(err)
	}
	want := string(bytes)
	got := mockWriter.String()
	if got != want {
		t.Errorf("wanted: %#v but got: %#v", want, got)
	}
}

func TestSlowPrintTwoFiles(t *testing.T) {
	t.Parallel()
	ignoredText := "foo"
	ignoredReader := strings.NewReader(ignoredText)
	testDelay := 10 * time.Millisecond
	mockWriter := &bytes.Buffer{}

	testFilePath1 := "testdata/1_2_3.txt"
	fi, err := os.Stat(testFilePath1)
	if err != nil {
		t.Fatal(err)
	}
	totalSize := fi.Size()

	testFilePath2 := "testdata/456.txt"
	fi, err = os.Stat(testFilePath2)
	if err != nil {
		t.Fatal(err)
	}
	totalSize += fi.Size()

	args := []string{
		testFilePath1,
		testFilePath2,
	}
	printer, err := slow.NewPrinter(
		slow.WithReader(ignoredReader),
		slow.WithArgs(args),
		slow.WithWriter(io.Writer(mockWriter)),
		slow.WithDelay(testDelay),
	)
	if err != nil {
		t.Error(err)
	}
	start := time.Now()
	err = printer.Print()
	if err != nil {
		t.Error(err)
	}
	wantElapsed := time.Duration(totalSize) * testDelay
	checkElapsed(t, start, wantElapsed, testDelay)

	bytes, err := os.ReadFile(testFilePath1)
	if err != nil {
		t.Fatal(err)
	}
	want := string(bytes)
	bytes, err = os.ReadFile(testFilePath2)
	if err != nil {
		t.Fatal(err)
	}
	want += string(bytes)

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
