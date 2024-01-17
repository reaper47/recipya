package app

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/briandowns/spinner"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func setup(exeDir string) {
	oldPaths := []string{
		filepath.Join(exeDir, RecipyaDB),
		filepath.Join(exeDir, RecipyaDB+"-shm"),
		filepath.Join(exeDir, RecipyaDB+"-wal"),
		filepath.Join(exeDir, FdcDB),
	}
	for _, path := range oldPaths {
		_, err := os.Stat(path)
		if errors.Is(err, os.ErrNotExist) {
			continue
		} else if err != nil {
			panic(err)
		}

		newPath := filepath.Join(DBBasePath, filepath.Base(path))
		err = os.Rename(path, newPath)
		if err != nil {
			fmt.Printf("Error moving %q to %q: %q\n", path, newPath, err)
			continue
		}
		fmt.Printf("Moved %q to %q\n", filepath.Base(path), filepath.Dir(newPath))
	}

	setupFDC()
	setupConfigFile(exeDir)

	fmt.Println("Recipya is properly set up.")
}

func setupFDC() {
	_, err := os.Stat(filepath.Join(DBBasePath, "fdc.db"))
	if errors.Is(err, os.ErrNotExist) {
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Prefix = "Fetching the FDC database... "
		s.FinalMSG = "Fetching the FDC database... " + greenText("Success") + "\n"
		s.Start()
		err = downloadFile(filepath.Join(DBBasePath, "fdc.db.zip"), "https://raw.githubusercontent.com/reaper47/recipya/main/deploy/fdc.db.zip")
		if err != nil {
			fmt.Printf("\n"+redText("Error downloading FDC database")+": %s\n", err)
			fmt.Println("Application setup will terminate")
			os.Exit(1)
		}
		s.Stop()
		_ = os.Remove(filepath.Join(DBBasePath, "fdc.db.zip"))
	} else {
		fmt.Println(greenText("OK") + " FDC database")
	}
}

func downloadFile(path, url string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res == nil {
		return errors.New("download file response is nil")
	}

	body := res.Body
	if body == nil {
		return errors.New("download file response body is nil")
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("file not found at %q", url)
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, body)
	if err != nil {
		return err
	}

	z, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = z.Close()
	}()

	destFile, err := os.OpenFile(filepath.Join(filepath.Dir(path), "fdc.db"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, z.File[0].Mode())
	if err != nil {
		return err
	}
	defer func() {
		_ = destFile.Close()
	}()

	zippedFile, err := z.File[0].Open()
	if err != nil {
		return err
	}
	defer func() {
		_ = zippedFile.Close()
	}()

	_, err = io.Copy(destFile, zippedFile)
	return err
}

func setupConfigFile(exeDir string) {
	if isRunningInDocker() {
		isEnvOk := true
		xenv := []string{"RECIPYA_SERVER_PORT", "RECIPYA_SERVER_URL"}
		for _, env := range xenv {
			if os.Getenv(env) == "" {
				isEnvOk = false
				fmt.Println("Missing required env variable:", env)
			}
		}

		if !isEnvOk {
			fmt.Println("Application setup will terminate")
			os.Exit(1)
		}

		return
	}

	configFilePath := filepath.Join(exeDir, "config.json")
	_, err := os.Stat(configFilePath)
	if err != nil {
		fmt.Print("Creating the configuration file... ")
		err = createConfigFile(configFilePath)
		if err != nil {
			fmt.Printf("\n"+redText("Error creating config file")+": %s\n", err)
			fmt.Println("Application setup will terminate")
			os.Exit(1)
		}
		fmt.Println(greenText("Success"))
	} else {
		fmt.Println(greenText("OK") + " Configuration file")
	}
}

func createConfigFile(path string) error {
	if isRunningInDocker() {
		return nil
	}

	var c ConfigFile
	r := bufio.NewReader(os.Stdin)
	fmt.Println()

	hasSendGrid := promptUser(r, "Do you have a SendGrid account? If not, important emails will not be sent [Y/n]", "y")
	if isYes(hasSendGrid) {
		c.Email.From = promptUser(r, "\tWhat is the email address of your SendGrid account?", "")
		c.Email.SendGridAPIKey = promptUser(r, "\tWhat is your SendGrid API key?", "")
	}

	hasVisionAPI := promptUser(r, "Do you have an Azure AI Vision account? If not, OCR features will be disabled. [Y/n]", "y")
	if isYes(hasVisionAPI) {
		c.Integrations.AzureComputerVision.ResourceKey = promptUser(r, "\tWhat is your resource key?", "")
		c.Integrations.AzureComputerVision.VisionEndpoint = promptUser(r, "\tWhat is your vision API endpoint?", "")
	}

	isInProd := promptUser(r, "Is your application in production? [y/N]", "N")
	c.Server.IsProduction = isYes(isInProd)

	url := promptUser(r, "What is the app's URL? (default, http://0.0.0.0)", "http://0.0.0.0")
	if runtime.GOOS == "windows" && strings.Contains(url, "0.0.0.0") {
		url = strings.Replace(url, "0.0.0.0", "127.0.0.1", 1)
	}
	c.Server.URL = url

	if !isYes(isInProd) {
		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			panic(err)
		}
		defer func() {
			_ = listener.Close()
		}()
		c.Server.Port = listener.Addr().(*net.TCPAddr).Port
	}

	j, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(path, j, os.ModePerm)
}

func isYes(s string) bool {
	return strings.HasPrefix(strings.ToLower(s), "y")
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
