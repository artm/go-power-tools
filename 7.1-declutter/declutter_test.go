package declutter_test

import (
	"declutter"
	"sort"
	"testing"
	"testing/fstest"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestFindOldFiles(t *testing.T) {
	t.Parallel()

	oldTime := time.Now().Add(-31 * 24 * time.Hour)
	newTime := time.Now().Add(-29 * 24 * time.Hour)

	fsys := fstest.MapFS{
		"old-file":         {ModTime: oldTime},
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
	got, err := c.FindOldFiles()
	if err != nil {
		t.Error(err)
	}

	want := []string{"old-file", "another-old-file", "folder/old-file"}
	sort.Strings(got)
	sort.Strings(want)
	if !cmp.Equal(got, want) {
		t.Errorf(
			"unexpected FildOldFiles() results, +wanted -got:\n%s\n",
			cmp.Diff(want, got),
		)
	}
}
