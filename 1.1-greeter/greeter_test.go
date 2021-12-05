package greeter_test

import (
	"bytes"
	"greeter"
	"io"
	"testing"
)

func TestGreeterAsksNameAndGreetsBack(t *testing.T) {
	mockReader := bytes.NewBufferString("artm\n")
	mockWriter := &bytes.Buffer{}
	mockGreeter := greeter.Greeter{
		In:  io.Reader(mockReader),
		Out: io.Writer(mockWriter),
	}
	mockGreeter.Greet()
	got := mockWriter.String()
	want := "What's your name? Hello, artm!\n"
	if got != want {
		t.Errorf("want: '%s', got: '%s'", want, got)
	}
}
