package declutter_test

import (
	"declutter"
	"io/fs"
	"testing"
	"testing/fstest"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestDeclutter(t *testing.T) {
	t.Parallel()

	oldTime := time.Now().Add(-31 * 24 * time.Hour)
	newTime := time.Now().Add(-29 * 24 * time.Hour)

	fsys := fstest.MapFS{
		"olf-file":         {ModTime: oldTime},
		"another-old-file": {ModTime: oldTime},
		"folder/old-file":  {ModTime: oldTime},
		"folder/new-file":  {ModTime: newTime},
		"new-file":         {ModTime: newTime},
	}

	c, err := declutter.NewCleaner(
		declutter.WithFS(fsys),
		declutter.WithAge("30d"),
	)
	if err != nil {
		t.Error(err)
	}
	err = c.Declutter()
	if err != nil {
		t.Error(err)
	}

	remaining := make([]string, 0)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			t.Error(err)
		}
		if !d.IsDir() {
			remaining = append(remaining, p)
		}
		return nil
	})
	want := []string{"folder/new-file", "new-file"}
	if !cmp.Equal(remaining, want) {
		t.Errorf("unexpected remaining +wanted -got:\n%s\n", cmp.Diff(want, remaining))
	}
}
