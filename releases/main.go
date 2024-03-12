package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var licenseFile string

func main() {
	packageName := flag.String("package", "", "The package name")
	tag := flag.String("tag", "", "The release tag")
	flag.Parse()

	if *packageName == "" || *tag == "" {
		fmt.Println("usage: main -package <package-name> -tag <tag>")
		os.Exit(1)
	}

	licenseFile = filepath.Join(".", "LICENSE")

	buildRelease(*packageName, *tag)
}

func buildRelease(packageName, tag string) {
	parts := filepath.SplitList(packageName)
	packageName = parts[len(parts)-1]

	platforms := []string{
		"darwin/amd64",
		"darwin/arm64",
		"linux/386",
		"linux/amd64",
		"linux/arm",
		"linux/arm64",
		"linux/riscv64",
		"linux/s390x",
		"windows/amd64",
		"windows/arm64",
	}

	for _, platform := range platforms {
		fmt.Printf("Building %s...\n", platform)
		build(platform, packageName, tag)
	}

	err := os.RemoveAll("builds")
	if err != nil {
		fmt.Printf("Deleting the builds folder failed: %q.\nAborting the script...\n", err)
		os.Exit(1)
	}
}

func build(platform, packageName, tag string) {
	goos, goarch, _ := strings.Cut(platform, "/")
	outputName := fmt.Sprintf("builds/%s-%s-%s", packageName, goos, goarch)

	binary := "builds/recipya"
	if goos == "windows" {
		binary += ".exe"
	}

	cmd := exec.Command("go", "build", "-ldflags=-s -w", "-o", binary, packageName)
	cmd.Env = append(os.Environ(), fmt.Sprintf("GOOS=%s", goos), fmt.Sprintf("GOARCH=%s", goarch))
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Running the build command failed: %q.\nAborting the script...\n", err)
		os.Exit(1)
	}

	_, err = os.Stat(binary)
	if err == nil {
		err = os.MkdirAll(filepath.Join(".", "releases", tag), os.ModePerm)
		if err != nil {
			fmt.Printf("Creating the tag's directory failed: %q.\nAborting the script...\n", err)
			os.Exit(1)
		}

		files := []string{binary, licenseFile}
		if goos == "windows" {
			updater := "builds/updater.exe"
			cmd = exec.Command("go", "-C", "./updater", "build", "-ldflags=-s -w", "-o", "../"+updater, "main.go")
			cmd.Env = append(os.Environ(), fmt.Sprintf("GOOS=%s", goos), fmt.Sprintf("GOARCH=%s", goarch))
			err = cmd.Run()
			if err != nil {
				fmt.Printf("Running the build command failed: %q.\nAborting the script...\n", err)
				os.Exit(1)
			}

			startScript, err := os.CreateTemp("", "start-script-*.bat")
			if err != nil {
				fmt.Printf("Failed to create temporary file: %q.\nAborting the script...\n", err)
				os.Exit(1)
			}
			defer os.Remove(startScript.Name())

			_, err = startScript.WriteString("cd /C \"%~dp0\"\r\ncall recipya.exe serve\r\n")
			if err != nil {
				fmt.Printf("Failed to write to temporary file: %q.\nAborting the script execution...\n", err)
				os.Exit(1)
			}
			startScript.Close()

			files = append(files, updater, startScript.Name())
		}

		dest := filepath.Join(".", "releases", tag, filepath.Join(filepath.Base(outputName)+".zip"))
		err = zipFiles(dest, files...)
		if err != nil {
			fmt.Printf("Zip failed for %s: %q.\nAborting the script...\n", platform, err)
			os.Exit(1)
		}
	}
}

func zipFiles(destFile string, files ...string) error {
	zipFile, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	w := zip.NewWriter(zipFile)
	defer w.Close()

	for _, file := range files {
		src, err := os.Open(file)
		if err != nil {
			return err
		}

		info, err := src.Stat()
		if err != nil {
			_ = src.Close()
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			_ = src.Close()
			return err
		}

		header.Name = filepath.Base(file)
		header.Method = zip.Deflate

		writer, err := w.CreateHeader(header)
		if err != nil {
			_ = src.Close()
			return err
		}

		_, err = io.Copy(writer, src)
		if err != nil {
			_ = src.Close()
			return err
		}

		_ = src.Close()
	}

	return w.Close()
}
