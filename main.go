package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func extract() {

	reader, err := zip.OpenReader("downloads/scripts.zip")
	if err != nil {
		fmt.Printf("Error opening ZIP file: %v\n", err)
		return
	}
	defer reader.Close()

	err = os.MkdirAll("scripts", 0755)
	if err != nil {
		fmt.Printf("Error creating destination directory: %v\n", err)
		return
	}

	totalFiles := len(reader.File)
	extractedFiles := 0

	for _, file := range reader.File {
		extractedFiles++
		fmt.Printf("\rExtracting files... %d/%d", extractedFiles, totalFiles)

		destPath := filepath.Join("scripts", file.Name)
		if !strings.HasPrefix(destPath, filepath.Clean("scripts")+string(os.PathSeparator)) {
			fmt.Printf("\nError: illegal file path: %s\n", file.Name)
			continue
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(destPath, file.Mode())
			continue
		}

		err = os.MkdirAll(filepath.Dir(destPath), 0755)
		if err != nil {
			fmt.Printf("\nError creating directory: %v\n", err)
			continue
		}

		destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			fmt.Printf("\nError creating file: %v\n", err)
			continue
		}

		srcFile, err := file.Open()
		if err != nil {
			fmt.Printf("\nError opening file in ZIP: %v\n", err)
			destFile.Close()
			continue
		}

		_, err = io.Copy(destFile, srcFile)
		srcFile.Close()
		destFile.Close()

		if err != nil {
			fmt.Printf("\nError extracting file: %v\n", err)
			continue
		}

		err = os.Chtimes(destPath, file.Modified, file.Modified)
		if err != nil {
			fmt.Printf("\nError setting file times: %v\n", err)
		}
	}

	fmt.Printf("\nExtraction completed to: %s\n", "scripts")

	fmt.Printf("Total files extracted: %d\n", extractedFiles)
}

func installExe(exePath string) error {
	absPath, err := filepath.Abs(exePath)
	if err != nil {
		return fmt.Errorf("error getting absolute path: %v", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", absPath)
	}

	fmt.Printf("Starting installation of: %s\n", absPath)
	startTime := time.Now()

	silentArgs := []string{
		"/S",          // Standard silent
		"/SILENT",     // Inno Setup
		"/VERYSILENT", // Inno Setup
		"/quiet",      // Microsoft installers
		"/Q",          // InstallShield
		"-q",          // Generic
		"--quiet",     // Generic
	}

	if strings.ToLower(filepath.Ext(absPath)) == ".msi" {
		cmd := exec.Command("msiexec", "/i", absPath, "/norestart")
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow:    true,
			CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
	} else {
		for _, arg := range silentArgs {
			cmd := exec.Command(absPath, arg)
			cmd.SysProcAttr = &syscall.SysProcAttr{
				HideWindow:    true,
				CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
			}
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err = cmd.Run()
			if err == nil {
				break
			}
		}
	}

	if err != nil {
		cmd := exec.Command(absPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
	}

	duration := time.Since(startTime)
	if err != nil {
		return fmt.Errorf("installation failed after %v: %v", duration, err)
	}

	fmt.Printf("Installation completed successfully in %v\n", duration)
	return nil
}

func installPythonPackages() error {
	packages := []string{"requests", "asyncio", "ixbrowser_local_api", "selenium", "pyppeteer"}
	for _, pkg := range packages {
		cmd := exec.Command("pip", "install", "--upgrade", pkg)
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to install package %s: %v", pkg, err)
		}
	}
	return nil
}

func downloadWithMegatools(megaURL, downloadDir string) error {
	if err := os.MkdirAll(downloadDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create download directory: %v", err)
	}

	cmd := exec.Command("./megatools.exe", "dl", "--path", downloadDir, megaURL)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("megadl command failed: %v", err)
	}

	return nil
}
func downloadFromMega(url string) {
	downloadDir := "downloads"

	fmt.Printf("Downloading from Mega.nz: %s\n", url)
	if err := downloadWithMegatools(url, downloadDir); err != nil {
		fmt.Printf("Error downloading file: %v\n", err)
		return
	}

	fmt.Printf("File successfully downloaded to directory: %s\n", downloadDir)
}

func installPython(installerPath string) error {
	cmd := exec.Command(installerPath, "InstallAllUsers=1", "PrependPath=1")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func addToPath(pythonDir string) error {
	pathEnv := os.Getenv("PATH")

	if strings.Contains(pathEnv, pythonDir) {
		fmt.Println("Python path is already in the system PATH.")
		return nil
	}

	newPath := pathEnv + ";" + pythonDir
	return os.Setenv("PATH", newPath)
}

func pythonInstall() {
	installerPath := "python-3.13.0-amd64.exe"

	fmt.Println("Starting Python installation...")
	if err := installPython(installerPath); err != nil {
		fmt.Printf("Error during installation: %v\n", err)
		return
	}

	pythonDir := "C:\\Users\\%USERNAME%\\AppData\\Local\\Programs\\Python\\Python3X"

	if err := addToPath(pythonDir); err != nil {
		fmt.Printf("Error adding Python to PATH: %v\n", err)
		return
	}

	fmt.Println("Python installed and PATH updated successfully.")
}

func installPython2(installerPath string) error {
	absPath, err := filepath.Abs(installerPath)
	if err != nil {
		return fmt.Errorf("error getting absolute path: %v", err)
	}

	cmd := exec.Command(absPath, "InstallAllUsers=1", "PrependPath=1")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func addToPath2(pythonDir string) error {
	pathEnv := os.Getenv("PATH")
	if strings.Contains(pathEnv, pythonDir) {
		fmt.Println("Python path is already in the system PATH.")
		return nil
	}

	newPath := pathEnv + ";" + pythonDir
	return os.Setenv("PATH", newPath)
}

func pythonInstall2() {
	installerPath := "python-3.13.0-amd64.exe"

	fmt.Println("Starting Python installation...")
	if err := installPython(installerPath); err != nil {
		fmt.Printf("Error during installation: %v\n", err)
		return
	}

	pythonDir := "C:\\Users\\%USERNAME%\\AppData\\Local\\Programs\\Python\\Python3X"
	if err := addToPath(pythonDir); err != nil {
		fmt.Printf("Error adding Python to PATH: %v\n", err)
		return
	}

	fmt.Println("Python installed and PATH updated successfully.")
}
func main() {

	installExe("ixBrowser_Setup_2_2_34.exe")
	pythonInstall()
	installPythonPackages()

}
