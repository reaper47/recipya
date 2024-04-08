package main

import (
	"context"
	"github.com/shirou/gopsutil/v3/process"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var logger *slog.Logger

func main() {
	file, err := os.Create("update.log")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	logger = slog.New(slog.NewTextHandler(file, nil))

	if runtime.GOOS != "windows" {
		logger.Error("This program may only be run on Windows.")
		os.Exit(1)
	}

	logger.Info("Update process started")
	waitForRecipyaToStop()

	exe, err := os.Executable()
	if err != nil {
		slog.Error("Failed to get executable path", "error", err)
		os.Exit(1)
	}
	logger.Info("Got updater executable path", "path", exe)

	path := filepath.Join(filepath.Dir(exe), "recipya.exe")
	logger.Info("Renaming files", "from", path+".new", "to", path)
	err = os.Rename(path+".new", path)
	if err != nil {
		logger.Error("Failed to rename files", "error", err)
		os.Exit(1)
	}

	logger.Info("Waiting 1s for rename to be done")
	time.Sleep(1 * time.Second)

	logger.Info("Starting 'recipya.exe serve'")
	cmd := exec.Command(filepath.Join(filepath.Dir(exe), "recipya.exe"), "serve")
	err = cmd.Run()
	if err != nil {
		logger.Error("Failed to run command", "command", cmd.Args, "error", err)
		os.Exit(1)
	}

	logger.Info("Software updated successfully")
}

func waitForRecipyaToStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	processes, err := process.Processes()
	if err != nil {
		logger.Error("Failed to get processes", "error", err)
		os.Exit(1)
	}

	var (
		proc    *process.Process
		isFound = false
	)

	for _, p := range processes {
		name, _ := p.Name()
		if strings.Contains(name, "recipya") {
			proc = p
			isFound = true
			break
		}
	}

	if isFound {
		n, _ := proc.Name()
		logger.Info("Recipya process found", "name", n, "pid", proc.Pid)

		for {
			isRunning, err := proc.IsRunningWithContext(ctx)
			if err != nil {
				logger.Error("Failed to get process with running status", "pid", proc.Pid, "error", err)
				os.Exit(1)
			}

			if !isRunning {
				logger.Info("Process stopped running")
				break
			}
		}
	}

	time.Sleep(1 * time.Second)
}
