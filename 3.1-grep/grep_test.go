package grep_test

import (
	"bytes"
	"grep"
	"io"
	"testing"
)

const (
	testInput = `foo
bar
foobar
kung-foo
schak
baz
`
	want = `foo
foobar
kung-foo
`
)

func TestGrep(t *testing.T) {
	mockReader := bytes.NewBufferString(testInput)
	mockWriter := &bytes.Buffer{}
	searcher := grep.NewSearcher(
		grep.WithWriter(io.Writer(mockWriter)),
		grep.WithReader(io.Reader(mockReader)),
	)
	searcher.Search("foo")
	got := mockWriter.String()
	if got != want {
		t.Errorf("want: %#v, got: %#v", want, got)
	}
}
