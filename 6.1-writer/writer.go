package writer

import "os"

func WriteToFile(path string, data []byte) error {
	err := os.WriteFile(path, data, 0600)
	if err != nil {
		return err
	}
	return os.Chmod(path, 0600)
}

func WriteZeros(path string, count int) error {
	data := make([]byte, count)
	return WriteToFile(path, data)
}
