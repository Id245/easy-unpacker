package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alexmullins/zip"
	"github.com/bodgit/sevenzip"
	"github.com/nwaples/rardecode"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Easy Unpacker - Extract various archive formats\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  %s <path-to-archive> <destination-directory>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Parameters:\n")
		fmt.Fprintf(os.Stderr, "  <path-to-archive>         Path to the archive file\n")
		fmt.Fprintf(os.Stderr, "  <destination-directory>   Directory for extraction\n")
		fmt.Fprintf(os.Stderr, "  -h                        Show help\n")
		fmt.Fprintf(os.Stderr, "\nSupported formats: .zip, .tar.gz, .tgz, .rar, .7z\n")
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s archive.zip ./extracted\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s backup.tar.gz ./backup\n", os.Args[0])
	}

	help := flag.Bool("h", false, "Show help")
	pass := flag.String("p", "", "Set password if required")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	args := flag.Args()

	if len(args) != 2 {
		fmt.Println("Error: Exactly two arguments required: archive path and destination directory")
		flag.Usage()
		return
	}

	path := args[0]
	dest := args[1]

	if !checkExists(path) {
		fmt.Println("No such file")
		return
	}

	ext := checkExtension(path)
	switch ext {
	case ".tar.gz", ".tgz":
		err := untar(path, dest)
		if err != nil {
			fmt.Println("Error during untar:", err)
			return
		}
	case ".zip":
		err := unzip(path, dest, *pass)
		if err != nil {
			fmt.Println("Error during unzip:", err)
			return
		}
	case ".rar":
		err := unrar(path, dest)
		if err != nil {
			fmt.Println("Error during unrar:", err)
			return
		}
	case ".7z":
		err := un7zip(path, dest)
		if err != nil {
			fmt.Println("Error during un7z:", err)
			return
		}
	default:
		fmt.Printf("Unsupported file format: %s\n", ext)
		return
	}
	fmt.Printf("Successfully extracted %s to %s\n", path, dest)
}

func checkExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func checkExtension(path string) string {
	ext := filepath.Ext(path)

	if (ext == ".gz") && (strings.HasSuffix(filepath.Base(path), ".tar.gz")) {
		ext = ".tar.gz"
	}
	return ext
}

func unzip(path string, dest string, password string) error {
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

		err = processFile(reader, filePath, file.Mode())
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

func untar(path string, dest string) error {
	archive, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer archive.Close()

	gzipReader, err := gzip.NewReader(archive)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		filePath := filepath.Join(dest, header.Name)
		if header.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		err = processFile(tarReader, filePath, os.FileMode(header.Mode))
		if err != nil {
			return err
		}
	}
	return nil
}

func unrar(path string, dest string) error {
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

		err = processFile(reader, filePath, os.FileMode(header.Mode()))
		if err != nil {
			return err
		}
	}
	return nil
}

func un7zip(path string, dest string) error {
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

		err = processFile(rFile, filePath, os.FileMode(file.Mode()))
		rFile.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func processFile(reader io.Reader, filePath string, mode os.FileMode) error {
	writer, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = io.Copy(writer, reader)
	return err
}
