package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bodgit/sevenzip"
	"github.com/nwaples/rardecode"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Easy Unpacker - Extract various archive formats\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  %s -p <path-to-archive> -d <destination-directory>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s <path-to-archive> <destination-directory>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Parameters:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nSupported formats: .zip, .tar.gz, .tgz, .rar, .7z\n")
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -p archive.zip -d ./extracted\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s backup.tar.gz ./backup\n", os.Args[0])
	}

	pathFlag := flag.String("p", "", "Path to the archive file (required unless provided as first argument)")
	destFlag := flag.String("d", "", "Destination directory for extraction (required unless provided as second argument)")
	help := flag.Bool("h", false, "Show help")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	var path, dest string

	args := flag.Args()

	if *pathFlag != "" || *destFlag != "" {

		path = *pathFlag
		dest = *destFlag

		if path == "" || dest == "" {
			fmt.Println("Error: Both -p and -d parameters are required when using flags")
			flag.Usage()
			return
		}
	} else if len(args) == 2 {
		path = args[0]
		dest = args[1]
	} else {
		fmt.Println("Error: Provide either flags (-p and -d) or two positional arguments")
		flag.Usage()
		return
	}

	fmt.Println("Path:", path)

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
		err := unzip(path, dest)
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

func unzip(path string, dest string) error {
	archive, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer archive.Close()

	for _, file := range archive.Reader.File {
		if file.FileInfo().IsDir() {
			filePath := filepath.Join(dest, file.Name)
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		reader, err := file.Open()
		if err != nil {
			return err
		}

		filePath := filepath.Join(dest, file.Name)
		dirPath := filepath.Dir(filePath)
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			reader.Close()
			return err
		}

		err = processFile(reader, filePath, file.Mode())
		reader.Close()
		if err != nil {
			return err
		}
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
