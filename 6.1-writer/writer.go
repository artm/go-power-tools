package writer

import "os"

const bufferSize = 10000

func WriteToFile(path string, data []byte) error {
	err := os.WriteFile(path, data, 0600)
	if err != nil {
		return err
	}
	return os.Chmod(path, 0600)
}

func WriteZeros(path string, count int) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	buffer := make([]byte, bufferSize)
	for count > 0 {
		size := bufferSize
		if count < bufferSize {
			size = count
		}
		_, err = f.Write(buffer[:size])
		if err != nil {
			return err
		}
		count -= size
	}
	return nil
}
