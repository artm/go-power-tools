package count_test

import (
	"bytes"
	"count"
	"testing"
)

func TestLines(t *testing.T) {
	t.Parallel()
	inputBuf := bytes.NewBufferString("1\n2\n3")
	c, err := count.NewCounter(
		count.WithInput(inputBuf),
	)
	if err != nil {
		t.Fatal(err)
	}
	want := 3
	got := c.Lines()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestMatches(t *testing.T) {
	t.Parallel()
	inputBuf := bytes.NewBufferString("a\nA\nA\nb\nA")
	c, err := count.NewCounter(
		count.WithInput(inputBuf),
		count.WithPattern("A"),
	)
	if err != nil {
		t.Fatal(err)
	}
	want := 3
	got := c.Lines()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestInsensitiveMatches(t *testing.T) {
	t.Parallel()
	inputBuf := bytes.NewBufferString("a\nA\nA\nb\nA\nc\nC")
	c, err := count.NewCounter(
		count.WithInput(inputBuf),
		count.WithPattern("a"),
		count.IgnoreCase(),
	)
	if err != nil {
		t.Fatal(err)
	}
	want := 4
	got := c.Lines()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
