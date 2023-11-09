package app

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/briandowns/spinner"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func setup() {
	reset := "\033[0m"
	greenText := func(s string) string {
		return "\033[32m" + s + reset
	}
	redText := func(s string) string {
		return "\033[31m" + s + reset
	}

	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dir := filepath.Dir(exe)

	_, err = os.Stat(filepath.Join(dir, "fdc.db"))
	if errors.Is(err, os.ErrNotExist) {
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Prefix = "Fetching the FDC database... "
		s.FinalMSG = "Fetching the FDC database... " + greenText("Success") + "\n"
		s.Start()
		err := downloadFile("fdc.db.zip", "https://raw.githubusercontent.com/reaper47/recipya/main/deploy/fdc.db.zip")
		if err != nil {
			fmt.Printf("\n"+redText("Error downloading FDC database")+": %s\n", err)
			fmt.Println("Application setup will terminate")
			os.Exit(1)
		}
		s.Stop()
	} else {
		fmt.Println(greenText("✓") + " FDC database")
	}

	_, err = os.Stat(filepath.Join(dir, "config.json"))
	if err != nil {
		fmt.Print("Creating the configuration file... ")
		err := createConfigFile()
		if err != nil {
			fmt.Printf("\n"+redText("Error creating config file")+": %s\n", err)
			fmt.Println("Application setup will terminate")
			os.Exit(1)
		}
		fmt.Println(greenText("Success"))
	} else {
		fmt.Println(greenText("✓") + " Configuration file")
	}

	fmt.Println("Recipya is properly set up.")
}

func createConfigFile() error {
	if isRunningInDocker() {
		return nil
	}

	var c ConfigFile
	r := bufio.NewReader(os.Stdin)
	fmt.Println()
	c.Email.From = promptUser(r, "What is the email address of your SendGrid account?", "")
	c.Email.SendGridAPIKey = promptUser(r, "What is your SendGrid API key?", "")
	isProduction := promptUser(r, "Is your application in production? [y/N]", "N")
	c.IsProduction = isProduction == "y"
	url := promptUser(r, "What is the app's URL? (default, http://0.0.0.0)", "http://0.0.0.0")
	if runtime.GOOS == "windows" && strings.Contains(url, "0.0.0.0") {
		url = strings.Replace(url, "0.0.0.0", "127.0.0.1", 1)
	}
	c.URL = url

	portStr := promptUser(r, "What is the port? (default, 8078)", "8078")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 8078
	}
	if port < 1024 {
		return fmt.Errorf("port must be greater than 1024, got %d", port)
	}
	c.Port = port

	j, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile("config.json", j, os.ModePerm)
}

func downloadFile(path, url string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("file not found at %q", url)
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return err
	}

	z, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer z.Close()

	dest := "fdc.db"
	destFile, err := os.OpenFile(filepath.Join(".", dest), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, z.File[0].Mode())
	if err != nil {
		return err
	}
	defer destFile.Close()

	zippedFile, err := z.File[0].Open()
	if err != nil {
		return err
	}
	defer zippedFile.Close()

	_, err = io.Copy(destFile, zippedFile)
	if err != nil {
		return err
	}
	return os.Remove(path)
}

func isRunningInDocker() bool {
	_, err := os.Stat("/.dockerenv")
	return err == nil
}

func promptUser(r *bufio.Reader, question string, def string) string {
	for {
		fmt.Print("\t" + question + " -> ")
		input, _ := r.ReadString('\n')
		input = strings.TrimSpace(input)

		if input != "" {
			return input
		}

		if input == "" && def != "" {
			return def
		}

		fmt.Println()
	}
}
