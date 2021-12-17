package writer

import (
	"os"
)

func RunCLI() {
	z, err := NewZeroer(
		FromArgs(os.Args[1:]),
	)
	checkErr(err)
	err = z.Write()
	checkErr(err)
}
