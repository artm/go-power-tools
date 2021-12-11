package wc_test

import (
	"bytes"
	"io"
	"os/exec"
	"strings"
	"testing"
	"wc"
)

var testCases = []string{
	"testdata/three_lines.txt",
	"-c testdata/three_lines.txt",
	"-m testdata/three_lines.txt",
	"-l testdata/three_lines.txt",
	"-L testdata/three_lines.txt",
	"-w testdata/three_lines.txt",
	"-w -c testdata/three_lines.txt",
	"-w -c -l testdata/three_lines.txt",
	"-w -c -l -m testdata/three_lines.txt",
	"-w -c -l -m -L testdata/three_lines.txt",
}

func TestWc(t *testing.T) {
	for _, testCase := range testCases {
		fakeWriter := bytes.NewBuffer([]byte{})
		args := strings.Split(testCase, " ")
		want := runWc(t, args...)
		wc, err := wc.NewWc(
			wc.WithArgs(args),
			wc.WithOutput(io.Writer(fakeWriter)),
		)
		if err != nil {
			t.Error(err)
		}
		err = wc.Count()
		if err != nil {
			t.Error(err)
		}
		got := fakeWriter.String()
		if err != nil {
			t.Error(err)
		}
		if got != want {
			t.Errorf(
				"args: %#v wanted: %#v got: %#v",
				args, want, got,
			)
		}
	}
}

func runWc(t *testing.T, args ...string) string {
	output, err := exec.Command("wc", args...).Output()
	if err != nil {
		t.Error(err)
	}
	return string(output)
}
