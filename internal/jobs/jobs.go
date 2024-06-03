package jobs

import (
	"github.com/go-co-op/gocron"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/services"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"time"
)

// ScheduleCronJobs schedules cron jobs for the web app. It starts the following jobs:
//
// - Clean Images: Removes unreferenced images from the data/images folder to save space.
//
// - Send queued emails
//
// - Backup data
//
// - Check for a new release
func ScheduleCronJobs(repo services.RepositoryService, files services.FilesService, email services.EmailService) {
	scheduler := gocron.NewScheduler(time.UTC)

	// Clean Images
	_, _ = scheduler.Every(1).MonthLastDay().Do(func() {
		rmFunc := func(file string) error {
			_ = os.Remove(filepath.Join(app.ThumbnailsDir, file))
			return os.Remove(filepath.Join(app.ImagesDir, file))
		}

		numFiles, numBytes := cleanImages(os.DirFS(app.ImagesDir), repo.Images(), rmFunc)

		var s string
		if numBytes > 0 {
			s = "(" + strconv.FormatFloat(float64(numBytes)/(1<<20), 'f', 2, 64) + " MB)"
		}

		slog.Info("Ran CleanImages job", "numImagesRemoved", numFiles, "spaceReclaimed", s)
	})

	// Send queued emails
	_, _ = scheduler.Every(1).Day().At("00:00").Do(func() {
		sent, remaining, err := email.SendQueue()
		slog.Info("Ran SendQueuedEmails job", "sent", sent, "remaining", remaining, "error", err)
	})

	// Backup data
	_, _ = scheduler.Every(3).Days().Do(func() {
		err := files.BackupGlobal()
		if err != nil {
			slog.Error("Global backup failed", "error", err)
			return
		}

		err = files.BackupUsersData(repo)
		if err != nil {
			slog.Error("User backups failed", "error", err)
			return
		}

		slog.Info("Backup successful")
	})

	// Check for a new release
	_, _ = scheduler.Every(3).Days().Do(func() {
		info, err := repo.CheckUpdate(files)
		if err != nil {
			slog.Error("Check for update failed", "error", err)
			return
		}

		app.Info.IsUpdateAvailable = info.IsUpdateAvailable
		app.Info.LastCheckedUpdateAt = info.LastCheckedUpdateAt
		app.Info.LastUpdatedAt = info.LastUpdatedAt

		slog.Info("Checked for an application update")
	})

	scheduler.StartAsync()
}

func cleanImages(dir fs.FS, usedImages []string, rmFileFunc func(path string) error) (numFilesDeleted, numBytesDeleted int64) {
	sort.Strings(usedImages)

	_ = fs.WalkDir(dir, ".", func(path string, d fs.DirEntry, _ error) error {
		if path == "." || d.IsDir() {
			return nil
		}

		_, found := slices.BinarySearch(usedImages, d.Name())
		if !found {
			info, err := d.Info()
			if err != nil {
				slog.Error("Clean images dir walk", "error", err)
				return err
			}

			err = rmFileFunc(path)
			if err != nil {
				slog.Error("Clean images walk", "error", err, "path", path)
				return err
			}

			numFilesDeleted++
			numBytesDeleted += info.Size()
		}
		return nil
	})
	return
}
