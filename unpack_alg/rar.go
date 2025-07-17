package unpackalg

import (
	"io"
	"os"
	"path/filepath"

	"github.com/nwaples/rardecode"
)

func Unrar(path string, dest string) error {
	archive, err := os.Open(path)
	if err != nil {
		return err
	}
	defer archive.Close()

	reader, err := rardecode.NewReader(archive, "")
	if err != nil {
		return err
	}

	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		filePath := filepath.Join(dest, header.Name)
		if header.IsDir {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		err = ProcessFile(reader, filePath, os.FileMode(header.Mode()))
		if err != nil {
			return err
		}
	}
	return nil
}
