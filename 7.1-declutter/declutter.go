package declutter

import "io/fs"

type cleaner struct {
	fs fs.FS
}

type option func(*cleaner) error

func NewCleaner(options ...option) (cleaner, error) {
	cleaner := cleaner{}
	return cleaner, nil
}

func WithFS(fs fs.FS) option {
	return func(c *cleaner) error {
		c.fs = fs
		return nil
	}
}

func WithAge(age string) option {
	return func(c *cleaner) error {
		return nil
	}
}

func (c *cleaner) Declutter() error {
	return nil
}
