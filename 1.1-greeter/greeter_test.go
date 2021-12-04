package greeter_test

import (
	"bytes"
	"greeter"
	"testing"
)

func TestGreeterAsksNameAndGreetsBack(t *testing.T) {
	fakeIn := bytes.NewBufferString("artm\n")
	fakeOut := &bytes.Buffer{}
	greeter.Greet(fakeIn, fakeOut)
	got := fakeOut.String()
	want := "What's your name? Hello, artm!\n"
	if got != want {
		t.Errorf("want: '%s', got: '%s'", want, got)
	}
}
