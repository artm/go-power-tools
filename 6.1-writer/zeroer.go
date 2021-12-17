package writer

import (
	"flag"
	"fmt"
	"os"
)

type zeroer struct {
	path    string
	size    int
	retries int
}

type option func(*zeroer) error

func WriteZeros(path string, size int) error {
	z, err := NewZeroer(
		WithPath(path),
		WithSize(size),
	)
	if err != nil {
		return err
	}
	return z.Write()
}

func NewZeroer(options ...option) (*zeroer, error) {
	z := &zeroer{}
	for _, opt := range options {
		err := opt(z)
		if err != nil {
			return nil, err
		}
	}
	return z, nil
}

func WithPath(path string) option {
	return func(z *zeroer) error {
		z.path = path
		return nil
	}
}

func WithSize(size int) option {
	return func(z *zeroer) error {
		z.size = size
		return nil
	}
}

func WithRetries(retries int) option {
	return func(z *zeroer) error {
		z.retries = retries
		return nil
	}
}

func FromArgs(args []string) option {
	return func(z *zeroer) error {
		fset := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		fset.IntVar(&z.size, "size", 1000, "number of zeros to write")
		fset.IntVar(&z.retries, "retries", 5, "retry this many times on write errors")
		err := fset.Parse(args)
		checkErr(err)
		rest := fset.Args()
		if len(rest) != 1 {
			checkErr(fmt.Errorf("single filename argument is required"))
		}
		z.path = rest[0]
		return nil
	}
}

func (z *zeroer) Write() error {
	f, err := os.Create(z.path)
	if err != nil {
		return err
	}
	defer f.Close()

	buffer := make([]byte, bufferSize)
	size := z.size
	for size > 0 {
		chunkSize := bufferSize
		if size < bufferSize {
			chunkSize = size
		}
		_, err = f.Write(buffer[:chunkSize])
		if err != nil {
			return err
		}
		size -= chunkSize
	}
	return nil
}
