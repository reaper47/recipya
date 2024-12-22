package app

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/disintegration/imaging"
	"github.com/gen2brain/webp"
	"github.com/reaper47/recipya/web"
	"image/jpeg"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func setup() {
	moveFileStructure()
	setupFDC()
	setupConfigFile()
	verifyMedia()
	verifyExtraSoftware()

	fmt.Println("Recipya is properly set up")
}

func moveFileStructure() {
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exeDir := filepath.Dir(exe)

	// Move configuration file
	configPathOld := filepath.Join(exeDir, "config.json")
	configPathNew := filepath.Join(filepath.Dir(DBBasePath), "config.json")

	_, err = os.Stat(configPathOld)
	if err == nil {
		temp, err := os.CreateTemp("", "*-config.json.bak")
		if err != nil {
			return
		}
		defer temp.Close()

		src, err := os.Open(configPathOld)
		if err != nil {
			fmt.Println("Could not open config.json: ", err)
		}

		_, err = io.Copy(temp, src)
		if err != nil {
			fmt.Println("Could copy config.json to temporary file: ", err)
		} else {
			fmt.Println("Copied config.json to temporary file ", temp.Name())
		}
		_ = src.Close()

		err = os.Rename(configPathOld, configPathNew)
		if err != nil {
			fmt.Printf("Could not move configuration file from %s to %s: %q", configPathOld, configPathNew, err)
		} else {
			fmt.Printf("Moved configuration file to new folder from %s to %s", configPathOld, configPathNew)
		}
	}

	// Move data folders
	dirs := map[string]string{"backup": BackupPath, "database": DBBasePath, "images": ImagesDir}
	count := 0
	for dir, newPath := range dirs {
		oldPath := filepath.Join(exeDir, "data", dir)
		_, err = os.Stat(oldPath)
		if err == nil {
			err = moveFiles(oldPath, newPath)
			if err != nil {
				fmt.Printf("Move %s folder to new location: %q", dir, err)
				continue
			}

			err = os.RemoveAll(oldPath)
			if err != nil {
				fmt.Printf("Please delete the old %s folder (%s) manually: %q", dir, oldPath, err)
			}

			count++
			fmt.Printf("Moved %s (%s) to new location %s", dir, oldPath, newPath)
		}
	}

	if count == len(dirs) {
		_ = os.RemoveAll(filepath.Join(exeDir, "data"))
	}

	// Move placeholder.webp to images dir
	_, err = os.Stat(filepath.Join(ImagesDir, "placeholder.recipe.original.webp"))
	if os.IsNotExist(err) {
		openFile, err := web.StaticFS.Open("static/img/recipes/placeholder.webp")
		if err != nil {
			return
		}
		defer openFile.Close()

		fRecipe, err := os.Create(filepath.Join(ImagesDir, "placeholder.recipe.webp"))
		if err != nil {
			panic(err)
		}
		defer fRecipe.Close()

		io.Copy(fRecipe, openFile)

		openFileCookbook, err := web.StaticFS.Open("static/img/cookbooks-new/placeholder.webp")
		if err != nil {
			panic(err)
		}
		defer openFileCookbook.Close()

		fCookbook, err := os.Create(filepath.Join(ImagesDir, "placeholder.cookbook.webp"))
		if err != nil {
			panic(err)
		}
		defer fCookbook.Close()

		io.Copy(fCookbook, openFileCookbook)
	}
}

func moveFiles(srcDir, destDir string) error {
	dir, err := os.Open(srcDir)
	if err != nil {
		return err
	}
	defer dir.Close()

	files, err := dir.Readdirnames(-1)
	if err != nil {
		return err
	}

	err = os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		return err
	}

	for _, f := range files {
		src := filepath.Join(srcDir, f)
		dest := filepath.Join(destDir, f)
		err = os.Rename(src, dest)
		if err != nil {
			return err
		}
	}

	// TODO: v1.3.0 - Move Images and Videos to Media folder. Media folder will have two subfolders: Images and Videos.
	return nil
}

func setupFDC() {
	_, err := os.Stat(filepath.Join(DBBasePath, "fdc.db"))
	if errors.Is(err, os.ErrNotExist) {
		if isRunningInDocker() {
			fmt.Println("Fetching the FDC database (62.6 MB)")
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Prefix = "Fetching the FDC database... "
		s.FinalMSG = "Fetching the FDC database... " + greenText("Success") + "\n"
		s.Start()
		err = downloadFile(filepath.Join(DBBasePath, "fdc.db.zip"), "fdc.db", "https://media.githubusercontent.com/media/reaper47/recipya/main/deploy/fdc.db.zip")
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

func downloadFile(path, filename, url string) error {
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

	destFile, err := os.OpenFile(filepath.Join(filepath.Dir(path), filename), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, z.File[0].Mode())
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
	return err
}

func setupConfigFile() {
	if isRunningInDocker() {
		isEnvOk := true
		xenv := []string{"RECIPYA_SERVER_PORT"}
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

		fmt.Println("Environment variables ok")
		return
	}

	configFilePath := filepath.Join(filepath.Dir(DBBasePath), "config.json")
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

	hasSendGrid := promptUser(r, "Do you have a SendGrid account? If not, important emails will not be sent [Y/n]", "n")
	if isYes(hasSendGrid) {
		c.Email.From = promptUser(r, "\tWhat is the email address of your SendGrid account?", "")
		c.Email.SendGridAPIKey = promptUser(r, "\tWhat is your SendGrid API key?", "")
	}

	hasOCR := promptUser(r, "Do you have an Azure AI Document Intelligence account? If not, OCR features will be disabled. [Y/n]", "n")
	if isYes(hasOCR) {
		c.Integrations.AzureDI.Key = promptUser(r, "\tWhat is your resource key?", "")
		c.Integrations.AzureDI.Endpoint = promptUser(r, "\tWhat is your endpoint?", "")
	}

	isAutologin := promptUser(r, "Do you wish to autologin? [y/N]", "N")
	c.Server.IsAutologin = isYes(isAutologin)

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
		defer listener.Close()
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

func verifyMedia() {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	msg := "Converting images to webp"
	s.Prefix = msg + "... "
	s.FinalMSG = greenText("OK ") + msg + "\n"
	s.Start()

	dir, err := os.ReadDir(ImagesDir)
	if err != nil {
		fmt.Printf(redText("Could not read images directory")+": %q\n", err)
		return
	}

	for _, entry := range dir {
		var (
			name = entry.Name()
			ext  = filepath.Ext(name)
		)

		if entry.IsDir() || ext != ".jpg" {
			continue
		}

		path := filepath.Join(ImagesDir, name)
		f, err := os.Open(path)
		if err != nil {
			fmt.Printf(redText("Could not open %q")+": %q\n", name, err)
			continue
		}

		origImg, err := jpeg.Decode(f)
		if err != nil {
			_ = f.Close()
			fmt.Printf(redText("Could not decode %q")+": %q\n", name, err)
			continue
		}
		_ = f.Close()

		origWidth := origImg.Bounds().Dx()
		origHeight := origImg.Bounds().Dy()

		// Thumbnail
		thumbnail, err := os.Create(filepath.Join(ThumbnailsDir, strings.TrimSuffix(name, ext)+ImageExt))
		if err != nil {
			fmt.Printf(redText("Could not create file")+": %q\n", name, err)
			continue
		}

		img := imaging.Resize(origImg, 384, int((float64(384)/float64(origWidth))*float64(origHeight)), imaging.NearestNeighbor)
		err = webp.Encode(thumbnail, img, webp.Options{Quality: 33})
		if err != nil {
			fmt.Printf(redText("Could not encode thumbnail %q to webp")+": %q\n", name, err)
			continue
		}
		_ = thumbnail.Close()

		// Normal
		webpFile, err := os.Create(filepath.Join(ImagesDir, strings.TrimSuffix(name, ext)+ImageExt))
		if err != nil {
			fmt.Printf(redText("Could not create file")+": %q\n", name, err)
			continue
		}

		err = webp.Encode(webpFile, origImg, webp.Options{Quality: 50})
		if err != nil {
			_ = webpFile.Close()
			fmt.Printf(redText("Could not encode %q to webp")+": %q\n", name, err)
			continue
		}

		_ = webpFile.Close()
		os.Remove(path)
	}
	s.Stop()
}

func verifyExtraSoftware() {
	err := exec.Command("ffmpeg", "-version").Run()
	Info.IsFFmpegInstalled = err == nil
	if !Info.IsFFmpegInstalled {
		msg := redText("X") + " Could not find ffmpeg"
		fmt.Printf(msg)
		switch runtime.GOOS {
		case "darwin":
			fmt.Println("\tPlease execute: brew install ffmpeg")
		case "linux":
			fmt.Println("\tPlease consult your package manager to install it.")
		case "windows":
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Prefix = msg + ". Attempting to install using winget..."
			s.FinalMSG = greenText("OK") + " FFmpeg installed\n"
			s.Start()

			err = exec.Command("winget", "install", "FFmpeg (Essentials Build)").Run()
			s.Stop()

			if err != nil {
				fmt.Printf(redText("X") + " Failed to install using winget. Please install manually: https://www.gyan.dev/ffmpeg/builds\n")
			} else {
				fmt.Println("Please reload your command prompt to refresh the environment variables.")
				os.Exit(0)
			}
		}
	} else {
		fmt.Println(greenText("OK") + " FFmpeg is installed")
	}
}
