package paths

import (
	"log"
	"os"
	"path/filepath"
)

// Data returns the path of the data folder.
func Data() string {
	return filepath.Join(exe(), "data")
}

// Images returns the path of the data images' folder.
func Images() string {
	return filepath.Join(Data(), "img")
}

func exe() string {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Dir(exe)
}
