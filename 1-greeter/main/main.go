package main

import (
	"greeter/greeter"
	"os"
)

func main() {
	greeter.Begreet(os.Stdout, os.Stdin)
}
