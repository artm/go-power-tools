package counter_test

import (
	"bytes"
	"counter"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/aquilax/truncate"
)

func TestCounterNext(t *testing.T) {
	t.Parallel()
	counter := counter.NewCounter()
	got := counter.Next()
	want := 0
	if got != want {
		t.Errorf("want: %#v, got: %#v", want, got)
	}
	got = counter.Next()
	want = 1
	if got != want {
		t.Errorf("want: %#v, got: %#v", want, got)
	}
	got = counter.Next()
	want = 2
	if got != want {
		t.Errorf("want: %#v, got: %#v", want, got)
	}
}

func TestCounterRun(t *testing.T) {
	t.Parallel()
	mockWriter := &bytes.Buffer{}
	counter := counter.Counter{
		Writer: io.Writer(mockWriter),
		Stop:   make(chan bool),
	}
	go counter.Run()
	time.Sleep(1 * time.Millisecond)
	counter.Stop <- false
	got := mockWriter.String()
	want := "0\n1\n2\n3\n"
	if !strings.HasPrefix(got, want) {
		got = truncate.Truncate(got, 20, truncate.DEFAULT_OMISSION, truncate.PositionMiddle)
		t.Errorf("expected Run() to output %#v etc; got %#v", want, got)
	}
}
