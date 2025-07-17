package unpackalg

import (
	"io"
	"os"
)

func ProcessFile(reader io.Reader, filePath string, mode os.FileMode) error {
	writer, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = io.Copy(writer, reader)
	return err
}
