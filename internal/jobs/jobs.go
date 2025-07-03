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
	"strings"
	"time"
)

// ScheduleCronJobs schedules cron jobs for the web app. It starts the following jobs:
//
// - Clean Media: Removes unreferenced images and videos from the data folder to save space.
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
		images, videos := repo.Media()

		numImages, numBytesImages := cleanMedia(os.DirFS(app.ImagesDir), images, func(file string) error {
			_ = os.Remove(filepath.Join(app.ThumbnailsDir, file))
			return os.Remove(filepath.Join(app.ImagesDir, file))
		})

		numVideos, numBytesVideos := cleanMedia(os.DirFS(app.VideosDir), videos, func(file string) error {
			return os.Remove(filepath.Join(app.VideosDir, file))
		})

		var s string
		if numBytesImages > 0 || numBytesVideos > 0 {
			s = "(" + strconv.FormatFloat(float64(numBytesImages+numBytesVideos)/(1<<20), 'f', 2, 64) + " MB)"
		}

		slog.Info("Ran CleanMedia job", "numImagesRemoved", numImages, "numVideosRemoved", numVideos, "spaceReclaimed", s)
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

func cleanMedia(dir fs.FS, used []string, rmFileFunc func(path string) error) (numFilesDeleted, numBytesDeleted int64) {
	sort.Strings(used)

	_ = fs.WalkDir(dir, ".", func(path string, d fs.DirEntry, _ error) error {
		if path == "." || d.IsDir() || strings.HasPrefix(d.Name(), "placeholder.") {
			return nil
		}

		_, found := slices.BinarySearch(used, d.Name())
		if !found {
			info, err := d.Info()
			if err != nil {
				slog.Error("Clean media dir walk", "error", err)
				return err
			}

			err = rmFileFunc(path)
			if err != nil {
				slog.Error("Clean media walk", "error", err, "path", path)
				return err
			}

			numFilesDeleted++
			numBytesDeleted += info.Size()
		}
		return nil
	})
	return
}
