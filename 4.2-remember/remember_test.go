package remember_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"remember"
	"testing"
)

func TestRemember(t *testing.T) {
	t.Parallel()
	dir, err := ioutil.TempDir("/tmp", "remember_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	tests := []struct {
		args []string
		want string
	}{
		{
			args: []string{},
			want: "",
		},
		{
			args: []string{"buy", "milk"},
			want: "",
		},
		{
			args: []string{},
			want: "buy milk\n",
		},
		{
			args: []string{"call mom"},
			want: "",
		},
		{
			args: []string{},
			want: "buy milk\ncall mom\n",
		},
	}
	for i, test := range tests {
		mockWriter := &bytes.Buffer{}
		memPath := path.Join(dir, "remember.json")

		err := remember.Run(
			remember.WithWriter(io.Writer(mockWriter)),
			remember.WithMemPath(memPath),
			remember.WithArgs(test.args),
		)
		if err != nil {
			t.Errorf("test %d: %#v", i, err)
		}
		got := mockWriter.String()
		if got != test.want {
			t.Errorf("test %d: wanted: %#v but got: %#v",
				i, test.want, got)
		}
	}
}
