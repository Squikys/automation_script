# Automated Software Installation Utility

## Overview

This Go script automates the process of downloading, installing, and configuring software packages, specifically designed for setting up a development environment with specific tools and dependencies.

## Features

- Download files from Mega.nz using megatools
- Extract ZIP archives
- Install executables silently
- Install Python
- Install Python packages
- Add Python to system PATH

## Prerequisites

- Go programming environment
- Windows operating system
- megatools (for Mega.nz downloads)
- Internet connection

## Installed Components

The script automatically installs:
- ixBrowser (version 2.2.34)
- Python 3.13.0
- Python packages:
  - requests
  - asyncio
  - ixbrowser_local_api
  - selenium
  - pyppeteer

## Configuration

### File Paths
Modify these paths in the script before running:
- `installerPath`: Path to Python installer
- `pythonDir`: Python installation directory
- Mega.nz download URL

### Customization

You can customize the installation by modifying:
- `packages` array in `installPythonPackages()`
- Silent installation arguments in `installExe()`
- Python installation method

## Functions

- `extract()`: Extracts files from a ZIP archive
- `installExe()`: Installs executables with multiple silent installation strategies
- `installPythonPackages()`: Installs specified Python packages
- `downloadWithMegatools()`: Downloads files from Mega.nz
- `installPython()`: Installs Python with specific configuration flags
- `addToPath()`: Adds Python directory to system PATH

## Usage

1. Ensure all required files are in the correct directories
2. Build the Go script
3. Run the executable

```bash
go build
./script
```

## Error Handling

The script includes comprehensive error handling for:
- File operations
- Executable installations
- Package installations
- PATH modifications

## Notes

- Requires administrative privileges
- Designed for Windows environments
- Always test in a controlled environment first

## Disclaimer

Use this script responsibly. Ensure you have the right to install and use the specified software.