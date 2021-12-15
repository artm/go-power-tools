package writer

import (
	"fmt"
	"os"
)

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

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
