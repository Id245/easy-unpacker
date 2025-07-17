package unpackalg

import (
	"os"
	"path/filepath"

	"github.com/bodgit/sevenzip"
)

func Un7zip(path string, dest string) error {
	archive, err := os.Open(path)
	if err != nil {
		return err
	}
	defer archive.Close()

	fileInfo, err := archive.Stat()
	if err != nil {
		return err
	}
	archSize := fileInfo.Size()

	reader, err := sevenzip.NewReader(archive, archSize)
	if err != nil {
		return err
	}

	for _, file := range reader.File {
		if file.FileInfo().IsDir() {
			filePath := filepath.Join(dest, file.Name)
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		rFile, err := file.Open()
		if err != nil {
			return err
		}

		filePath := filepath.Join(dest, file.Name)
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			rFile.Close()
			return err
		}

		err = ProcessFile(rFile, filePath, os.FileMode(file.Mode()))
		rFile.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
