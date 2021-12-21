package declutter_test

import (
	"declutter"
	"testing"
	"testing/fstest"
)

func TestDeclutter(t *testing.T) {
	t.Parallel()

	fs := fstest.MapFS{
		"olf-file":         {},
		"another-old-file": {},
		"folder/old-file":  {},
		"folder/new-file":  {},
		"new-file":         {},
	}

	c, err := declutter.NewCleaner(
		declutter.WithFS(fs),
		declutter.WithAge("30d"),
	)
	if err != nil {
		t.Error(err)
	}
	err = c.Declutter()
	if err != nil {
		t.Error(err)
	}
}
