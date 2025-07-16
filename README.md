# Easy Unpacker

Easy Unpacker is a lightweight command-line tool written in Go that extracts various archive formats through a simple, unified interface.

## Features

- Support for multiple archive formats:
  - ZIP (`.zip`)
  - TAR.GZ (`.tar.gz`, `.tgz`)
  - RAR (`.rar`) 
  - 7-Zip (`.7z`)
- Preserves directory structure during extraction
- Maintains original file permissions
- Clear error reporting
- Automatic directory creation

> **Note:** The project is actively being developed with plans to support additional archive formats in the future.

## Installation

### Prerequisites

- Go 1.16 or later

### Option 1: Using Go Modules (Recommended)

```bash
# Clone the repository
git clone https://github.com/yourusername/easy-unpacker.git
cd easy-unpacker

# Initialize Go modules
go mod init easy-unpacker

# Download dependencies and build
go mod tidy
go build
```

### Option 2: Manual Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/easy-unpacker.git
cd easy-unpacker

# Install dependencies explicitly
go get github.com/nwaples/rardecode
go get github.com/bodgit/sevenzip

# Build the executable
go build
```

## Usage

The Easy Unpacker can be used with positional arguments:

```bash
./easy-unpacker <path-to-archive> <destination-directory>
```

### Parameters

- Path to the archive file (required)
- Destination directory for extracted files (required)

If the destination directory doesn't exist, it will be automatically created.

### Examples

Extract a ZIP archive:
```bash
./easy-unpacker documents.zip ./extracted_docs
```

Extract a TAR.GZ archive:
```bash
./easy-unpacker backup.tar.gz ./restored_backup
```

