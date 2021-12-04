package greeter_test

import (
	"bytes"
	"greeter/greeter"
	"testing"
)

func TestAskNameAsks(t *testing.T) {
	fakeOut := &bytes.Buffer{}
	fakeIn := bytes.NewBufferString("artm")
	greeter.AskName(fakeOut, fakeIn)
	got := fakeOut.String()
	want := "What's your name? "
	if got != want {
		t.Errorf("want: '%s', got: '%s'", want, got)
	}
}

func TestAskNameReturnsAnswer(t *testing.T) {
	fakeOut := &bytes.Buffer{}
	fakeIn := bytes.NewBufferString("artm\n")
	want := "artm"
	got := greeter.AskName(fakeOut, fakeIn)
	if got != want {
		t.Errorf("want: '%s', got: '%s'", want, got)
	}
}

func TestGreetGreets(t *testing.T) {
	fakeOut := &bytes.Buffer{}
	greeter.Greet(fakeOut, "artm")
	got := fakeOut.String()
	want := "Hello, artm!\n"
	if got != want {
		t.Errorf("want: '%s', got: '%s'", want, got)
	}
}
