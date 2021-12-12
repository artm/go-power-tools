package writer

import (
	"flag"
	"fmt"
	"os"
)

type cli struct {
	path string
	size int
}

type option func(*cli) error

const bufferSize = 10000

func WriteToFile(path string, data []byte) error {
	err := os.WriteFile(path, data, 0600)
	if err != nil {
		return err
	}
	return os.Chmod(path, 0600)
}

func WriteZeros(path string, size int) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	buffer := make([]byte, bufferSize)
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

func RunCLI() {
	cli, err := NewCLI(
		FromArgs(os.Args[1:]),
	)
	checkErr(err)
	err = cli.Write()
	checkErr(err)
}

func NewCLI(options ...option) (*cli, error) {
	cli := &cli{}
	for _, opt := range options {
		err := opt(cli)
		if err != nil {
			return nil, err
		}
	}
	return cli, nil
}

func FromArgs(args []string) option {
	return func(cli *cli) error {
		fset := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		fset.IntVar(&cli.size, "size", 1000, "number of zeros to write")
		err := fset.Parse(args)
		checkErr(err)
		rest := fset.Args()
		if len(rest) != 1 {
			checkErr(fmt.Errorf("single filename argument is required"))
		}
		cli.path = rest[0]
		return nil
	}
}

func (cli *cli) Write() error {
	return WriteZeros(cli.path, cli.size)
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
