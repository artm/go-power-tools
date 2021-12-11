package wc_test

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"testing"
	"wc"

	"github.com/google/go-cmp/cmp"
)

var testCases = []struct {
	args  string
	stdin string
}{
	{"", "abra\ncadabra"},
	{"-c", "abra\ncadabra"},
	{"-m", "abra\ncadabra"},
	{"-l", "abra\ncadabra"},
	{"-L", "abra\ncadabra"},
	{"-w", "abra\ncadabra"},
	{"-", "abra\ncadabra"},
	{"-c -", "abra\ncadabra"},
	{"-m -", "abra\ncadabra"},
	{"-l -", "abra\ncadabra"},
	{"-L -", "abra\ncadabra"},
	{"-w -", "abra\ncadabra"},
	{"testdata/three_lines.txt", ""},
	{"-c testdata/three_lines.txt", ""},
	{"-m testdata/three_lines.txt", ""},
	{"-l testdata/three_lines.txt", ""},
	{"-L testdata/three_lines.txt", ""},
	{"-w testdata/three_lines.txt", ""},
	{"-w -c testdata/three_lines.txt", ""},
	{"-w -c -l testdata/three_lines.txt", ""},
	{"-w -c -l -m testdata/three_lines.txt", ""},
	{"-w -c -l -m -L testdata/three_lines.txt", ""},
	{"testdata/three_lines.txt testdata/lorem.txt", ""},
	{"-c testdata/three_lines.txt testdata/lorem.txt", ""},
	{"-m testdata/three_lines.txt testdata/lorem.txt", ""},
	{"-l testdata/three_lines.txt testdata/lorem.txt", ""},
	{"-L testdata/three_lines.txt testdata/lorem.txt", ""},
	{"-w testdata/three_lines.txt testdata/lorem.txt", ""},
	{"-w -c testdata/three_lines.txt testdata/lorem.txt", ""},
	{"-w -c -l testdata/three_lines.txt testdata/lorem.txt", ""},
	{"-w -c -l -m testdata/three_lines.txt testdata/lorem.txt", ""},
	{"-w -c -l -m -L testdata/three_lines.txt testdata/lorem.txt", ""},
	{"testdata/three_lines.txt -", "abra\ncadabra"},
	{"-c testdata/three_lines.txt -", "abra\ncadabra"},
	{"-m testdata/three_lines.txt -", "abra\ncadabra"},
	{"-l testdata/three_lines.txt -", "abra\ncadabra"},
	{"-L testdata/three_lines.txt -", "abra\ncadabra"},
	{"-w testdata/three_lines.txt -", "abra\ncadabra"},
	{"-w -c testdata/three_lines.txt -", "abra\ncadabra"},
	{"-w -c -l testdata/three_lines.txt -", "abra\ncadabra"},
	{"-w -c -l -m testdata/three_lines.txt -", "abra\ncadabra"},
	{"-w -c -l -m -L testdata/three_lines.txt -", "abra\ncadabra"},
}

func TestWc(t *testing.T) {
	for _, testCase := range testCases {
		fakeInput := strings.NewReader(testCase.stdin)
		fakeOutput := bytes.NewBuffer([]byte{})
		var args []string
		if testCase.args == "" {
			args = make([]string, 0)
		} else {
			args = strings.Split(testCase.args, " ")
		}
		want := runWc(t, testCase.stdin, args...)
		wc, err := wc.NewWc(
			wc.WithInput(io.Reader(fakeInput)),
			wc.WithArgs(args),
			wc.WithOutput(io.Writer(fakeOutput)),
		)
		if err != nil {
			t.Error(err)
		}
		err = wc.Count()
		if err != nil {
			t.Error(err)
		}
		got := fakeOutput.String()
		if err != nil {
			t.Error(err)
		}
		if got != want {
			t.Errorf(textDiff(testCase.args, want, got))
		}
	}
}

func textDiff(title, want, got string) string {
	wants := strings.Split(strings.TrimRight(want, "\n"), "\n")
	gots := strings.Split(strings.TrimRight(got, "\n"), "\n")
	return fmt.Sprintf(
		"%s (-want +got):\n%s",
		title, cmp.Diff(wants, gots),
	)
}

func runWc(t *testing.T, input string, args ...string) string {
	cmd := exec.Command("wc", args...)
	if input != "" {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			t.Error(err)
		}
		fmt.Fprint(stdin, input)
		err = stdin.Close()
		if err != nil {
			t.Error(err)
		}
	}
	output, err := cmd.Output()
	if err != nil {
		t.Errorf("wc args: %#v: %s", args, err)
	}
	return string(output)
}
