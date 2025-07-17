package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	unpackalg "easy-unpacker/unpack_alg"
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
		err := unpackalg.Untar(path, dest)
		if err != nil {
			fmt.Println("Error during untar:", err)
			return
		}
	case ".zip":
		err := unpackalg.Unzip(path, dest, *pass)
		if err != nil {
			fmt.Println("Error during unzip:", err)
			return
		}
	case ".rar":
		err := unpackalg.Unrar(path, dest)
		if err != nil {
			fmt.Println("Error during unrar:", err)
			return
		}
	case ".7z":
		err := unpackalg.Un7zip(path, dest)
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
