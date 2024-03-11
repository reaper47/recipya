package main

import (
	"context"
	"github.com/shirou/gopsutil/v3/process"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func main() {
	file, err := os.Create("update.log")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	log.SetOutput(file)

	if runtime.GOOS != "windows" {
		log.Fatalln("This program may only be run on Windows.")
	}

	log.Println("Update process started")
	waitForRecipyaToStop()

	exe, err := os.Executable()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Got updater executable path: ", exe)

	path := filepath.Join(filepath.Dir(exe), "recipya.exe")
	log.Printf("Renaming %q to %q", path+".new", path)
	err = os.Rename(path+".new", path)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Waiting 1s for rename to be done")
	time.Sleep(1 * time.Second)

	log.Println("Starting 'recipya.exe serve'")
	err = exec.Command(filepath.Join(filepath.Dir(exe), "recipya.exe"), "serve").Run()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Software updated successfully")
}

func waitForRecipyaToStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	processes, err := process.Processes()
	if err != nil {
		log.Fatalln(err)
	}

	var proc *process.Process = nil
	for _, p := range processes {
		name, _ := p.Name()
		if strings.Contains(name, "recipya") {
			proc = p
			break
		}
	}

	if proc != nil {
		n, _ := proc.Name()
		log.Printf("Recipya process %q found with pid %d", n, proc.Pid)

		for {
			isRunning, err := proc.IsRunningWithContext(ctx)
			if err != nil {
				log.Fatalln("Error getting process running status: ", err)
			}

			if !isRunning {
				log.Println("Process stopped running")
				break
			}
		}
	}

	time.Sleep(1 * time.Second)
}
