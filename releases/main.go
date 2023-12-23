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

var (
	licenseFile string
	configFile  string
)

func main() {
	packageName := flag.String("package", "", "The package name")
	tag := flag.String("tag", "", "The release tag")
	flag.Parse()

	if *packageName == "" || *tag == "" {
		fmt.Println("usage: main -package <package-name> -tag <tag>")
		os.Exit(1)
	}

	licenseFile = filepath.Join(".", "LICENSE")
	configFile = filepath.Join(".", "deploy", "config.example.json")

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
		fmt.Printf("Deleting the builds folder failed: %q.\nAborting the script execution...\n", err)
		os.Exit(1)
	}
}

func build(platform, packageName, tag string) {
	goos, goarch, _ := strings.Cut(platform, "/")
	outputName := fmt.Sprintf("builds/%s-%s-%s", packageName, goos, goarch)
	if goos == "windows" {
		outputName += ".exe"
	}

	cmd := exec.Command("go", "build", "-ldflags=-s -w", "-o", outputName, packageName)
	cmd.Env = append(os.Environ(), fmt.Sprintf("GOOS=%s", goos), fmt.Sprintf("GOARCH=%s", goarch))
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Running the build command failed %#v: %q.\nAborting the script execution...\n", cmd.Args, err)
		os.Exit(1)
	}

	_, err = os.Stat(outputName)
	if err == nil {
		err = os.MkdirAll(filepath.Join(".", "releases", tag), os.ModePerm)
		if err != nil {
			fmt.Printf("Creating the tag's directory failed: %q.\nAborting the script execution...\n", err)
			os.Exit(1)
		}

		destFile := filepath.Join(".", "releases", tag, filepath.Join(filepath.Base(outputName)+".zip"))
		destFile = strings.Replace(destFile, ".exe", "", 1)
		err = zipFiles(destFile, outputName, licenseFile, configFile)
		if err != nil {
			fmt.Printf("Zip failed for %s: %q.\nAborting the script execution...\n", platform, err)
			os.Exit(1)
		}
	}
}

func zipFiles(destFile string, files ...string) error {
	zipFile, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer func() {
		_ = zipFile.Close()
	}()

	w := zip.NewWriter(zipFile)
	defer func() {
		_ = w.Close()
	}()

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
