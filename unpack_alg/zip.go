package unpackalg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alexmullins/zip"
)

func Unzip(path string, dest string, password string) error {
	archive, err := zip.OpenReader(path)
	if err != nil {
		fmt.Println("Error during unzip:", err)
		return err
	}
	defer archive.Close()

	for _, file := range archive.Reader.File {
		if file.IsEncrypted() {
			if password != "" {
				file.SetPassword(password)
			} else {
				return fmt.Errorf("archive is encrypted, provide a password via -p <password>")
			}
		}

		if file.FileInfo().IsDir() {
			filePath := filepath.Join(dest, file.Name)
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				fmt.Println("Error creating directory:", err)
				return UnzipSystem(path, dest, password)
			}
			continue
		}

		reader, err := file.Open()
		if err != nil {
			fmt.Println("Error opening file:", err)
			if strings.Contains(err.Error(), "decrypt") {
				fmt.Println("Decryption error: wrong password or unsupported encryption method")
			}
			return UnzipSystem(path, dest, password)
		}

		filePath := filepath.Join(dest, file.Name)
		dirPath := filepath.Dir(filePath)
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			reader.Close()
			fmt.Println("Error creating directory:", err)
			return UnzipSystem(path, dest, password)
		}

		err = ProcessFile(reader, filePath, file.Mode())
		reader.Close()
		if err != nil {
			fmt.Println("Error processing file:", err)
			return UnzipSystem(path, dest, password)
		}
	}
	return nil
}

func UnzipSystem(path string, dest string, password string) error {
	_, err := exec.LookPath("unzip")
	if err != nil {
		return fmt.Errorf("system unzip utility not found")
	}

	fmt.Println("Easy-unpacker failed to decrypt. Trying system unzip...")

	var args []string

	//TODO: implement user notification if file with same name as unpacking already exists

	if password != "" {
		args = []string{"-P", password, path, "-d", dest}
	} else {
		args = []string{"-q", "-n", path, "-d", dest}
	}

	cmd := exec.Command("unzip", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("system unzip failed: %v\n%s", err, string(output))
	}
	return nil
}
