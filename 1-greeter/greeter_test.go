package greeter_test

import (
	"bytes"
	"greeter/greeter"
	"testing"
)

func TestGreeterAsksNameAndGreetsBack(t *testing.T) {
	fakeOut := &bytes.Buffer{}
	fakeIn := bytes.NewBufferString("artm\n")
	greeter.Greet(fakeOut, fakeIn)
	got := fakeOut.String()
	want := "What's your name? Hello, artm!\n"
	if got != want {
		t.Errorf("want: '%s', got: '%s'", want, got)
	}
}
