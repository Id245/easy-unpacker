# Easy Unpacker

Easy Unpacker is a lightweight command-line tool written in Go that extracts various archive formats through a simple, unified interface.

## Features

- Support for multiple archive formats:
  - ZIP (`.zip`)
  - TAR.GZ (`.tar.gz`, `.tgz`)
  - RAR (`.rar`) 
  - 7-Zip (`.7z`)
- Password-protected archive support
- Preserves directory structure during extraction
- Maintains original file permissions
- Clear error reporting
- Automatic directory creation
- Fallback to system utilities when needed

> **Note:** The project is actively being developed with plans to support additional archive formats in the future.

## Installation

### Prerequisites

- Go 1.16 or later
- System `unzip` utility (optional, for advanced ZIP extraction)

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
- `-p <password>` - Password for encrypted archives (optional)
- `-h` - Show help information

If the destination directory doesn't exist, it will be automatically created.

### Password-Protected Archives

For encrypted archives, use the `-p` flag followed by the password:

```bash
./easy-unpacker -p mypassword ./encrypted.zip ./extracted
```

Note: Password protection is currently supported for ZIP archives only.

### Fallback Mechanism

Easy Unpacker uses Go libraries for extraction by default, but will automatically fall back to system utilities in the following cases:

- When dealing with complex encrypted ZIP formats
- When the built-in libraries encounter extraction errors
- For archives with non-standard compression methods

This fallback mechanism requires the corresponding system utilities (`unzip` for ZIP files) to be installed on your system.

## Examples

Extract a standard ZIP archive:
```bash
./easy-unpacker archive.zip ./extracted_files
```

Extract a password-protected ZIP archive:
```bash
./easy-unpacker -p secretpassword secure_archive.zip ./extracted_files
```

