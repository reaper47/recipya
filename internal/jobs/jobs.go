package jobs

import (
	"archive/zip"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/services"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"time"
)

// ScheduleCronJobs schedules cron jobs for the web app. It starts the following jobs:
//
// - Clean Images: Removes unreferenced images from the data/img folder to save space.
//
// - Send Queued Emails
//
// - Backup data
func ScheduleCronJobs(repo services.RepositoryService, imagesDir string, email services.EmailService) {
	scheduler := gocron.NewScheduler(time.UTC)

	_, _ = scheduler.Every(1).MonthLastDay().Do(func() {
		rmFunc := func(file string) error {
			return os.Remove(filepath.Join(imagesDir, file))
		}
		numFiles, numBytes := cleanImages(os.DirFS(imagesDir), repo.Images(), rmFunc)

		var s string
		if numBytes > 0 {
			s = "(" + strconv.FormatFloat(float64(numBytes)/(1<<20), 'f', 2, 64) + " MB)"
		}
		log.Printf("CleanImages: Removed %d unreferenced images %s", numFiles, s)
	})

	_, _ = scheduler.Every(1).Day().At("00:00").Do(func() {
		sent, remaining, err := email.SendQueue()
		log.Printf("SendQueuedEmails: Sent %d | Remaining %d | Error: %q", sent, remaining, err)
	})

	_, _ = scheduler.Every(1).Sunday().Do(backupData)

	scheduler.StartAsync()
}

func backupData() {
	name := fmt.Sprintf("recipya-%s.bak", time.Now().Format(time.DateOnly))
	target := filepath.Join(app.BackupPath, name)
	zf, err := os.Create(target)
	if err != nil {
		log.Printf("Could not create backup %q", name)
		return
	}
	defer func() {
		_ = zf.Close()
	}()

	zw := zip.NewWriter(zf)
	defer zw.Close()

	source := filepath.Dir(app.DBBasePath)
	backupBase := filepath.Base(app.BackupPath)

	err = filepath.Walk(source, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		omit := []string{"backup", "fdc.db", app.RecipyaDB + "-wal", app.RecipyaDB + "-shm"}
		if filepath.Base(filepath.Dir(path)) == backupBase || slices.Contains(omit, info.Name()) {
			return nil
		}

		h, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		h.Method = zip.Deflate
		h.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			h.Name += "/"
		}

		w, err := zw.CreateHeader(h)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() {
			_ = f.Close()
		}()

		_, err = io.Copy(w, f)
		return err
	})
	if err != nil {
		log.Printf("Could not assemble backup %q", name)
		return
	}

	log.Printf("Backup successful: %q", target)
}

func cleanImages(dir fs.FS, usedImages []string, rmFileFunc func(path string) error) (numFilesDeleted, numBytesDeleted int64) {
	sort.Strings(usedImages)

	_ = fs.WalkDir(dir, ".", func(path string, d fs.DirEntry, err error) error {
		if path == "." {
			return nil
		}

		_, found := slices.BinarySearch(usedImages, d.Name())
		if !found {
			info, err := d.Info()
			if err != nil {
				log.Printf("clean images dir walk error: %s", err)
				return err
			}

			err = rmFileFunc(path)
			if err != nil {
				log.Printf("clean images walk '%s': %s", path, err)
				return err
			}

			numFilesDeleted++
			numBytesDeleted += info.Size()
		}
		return nil
	})
	return
}
