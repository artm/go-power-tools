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
