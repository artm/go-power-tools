package writer

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type zeroer struct {
	output  io.WriteCloser
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

func WithOutput(output io.WriteCloser) option {
	return func(z *zeroer) error {
		z.output = output
		return nil
	}
}

func WithPath(path string) option {
	return func(z *zeroer) error {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		z.output = io.WriteCloser(f)
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
		switch len(rest) {
		case 1:
			return WithPath(rest[0])(z)
		case 0:
			if z.output == nil {
				checkErr(fmt.Errorf("a filename argument is required"))
			}
		default:
			checkErr(fmt.Errorf("a single filename argument is expected"))
		}
		return nil
	}
}

func (z *zeroer) Write() error {
	defer z.output.Close()

	buffer := make([]byte, bufferSize)
	size := z.size
	for size > 0 {
		chunkSize := bufferSize
		if size < bufferSize {
			chunkSize = size
		}
		_, err := z.output.Write(buffer[:chunkSize])
		if err != nil {
			return err
		}
		size -= chunkSize
	}
	return nil
}
