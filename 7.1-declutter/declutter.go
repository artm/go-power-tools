package declutter

import (
	"io/fs"
	"time"

	durstr "github.com/xhit/go-str2duration/v2"
)

type cleaner struct {
	fsys         fs.FS
	modThreshold time.Time
}

type option func(*cleaner) error

func NewCleaner(options ...option) (*cleaner, error) {
	c := &cleaner{
		modThreshold: time.Now().AddDate(0, 0, -1),
	}
	for _, opt := range options {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func WithFS(fs fs.FS) option {
	return func(c *cleaner) error {
		c.fsys = fs
		return nil
	}
}

func WithAge(ageString string) option {
	return func(c *cleaner) error {
		age, _ := durstr.ParseDuration(ageString)
		c.modThreshold = time.Now().Add(-age)
		return nil
	}
}

func (c *cleaner) FindOldFiles() ([]string, error) {
	oldFiles := make([]string, 0)

	err := fs.WalkDir(c.fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		if d.IsDir() || info.ModTime().After(c.modThreshold) {
			return nil
		}
		oldFiles = append(oldFiles, p)
		return nil
	})

	return oldFiles, err
}
