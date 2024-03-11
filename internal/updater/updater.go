package updater

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"github.com/blang/semver"
	"github.com/google/go-github/v59/github"
	"golang.org/x/exp/slices"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// ErrNoUpdate is the error for when no update is available.
var ErrNoUpdate = errors.New("already latest version")

// IUpdater is the updater's interface
type IUpdater interface {
	// IsLatest checks whether there is a software update.
	IsLatest(current semver.Version) (bool, *github.RepositoryRelease, error)

	// Update updates the application to the latest version.
	Update(current semver.Version) error
}

// Updater is the main application's updater.
type Updater struct{}

// IsLatest checks whether there is a software update.
func (u *Updater) IsLatest(current semver.Version) (bool, *github.RepositoryRelease, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	gh := github.NewClient(http.DefaultClient)

	rel, res, err := gh.Repositories.GetLatestRelease(ctx, "reaper47", "recipya")
	if err != nil {
		return false, nil, err
	}

	if res.StatusCode != http.StatusOK {
		return false, nil, fmt.Errorf("got status code %d instead of 200 when fetching latest releases on GitHub", res.StatusCode)
	}

	version, err := semver.Parse(strings.Replace(rel.GetTagName(), "v", "", 1))
	if err != nil {
		return false, nil, err
	}

	return version.LTE(current), rel, nil
}

// Update updates the application to the latest version.
func (u *Updater) Update(current semver.Version) error {
	// Check if latest
	isLatest, rel, err := u.IsLatest(current)
	if err != nil {
		return err
	}

	if isLatest {
		return ErrNoUpdate
	}

	// Find asset
	name := fmt.Sprintf("recipya-%s-%s", runtime.GOOS, runtime.GOARCH)
	i := slices.IndexFunc(rel.Assets, func(asset *github.ReleaseAsset) bool { return *asset.Name == name+".zip" })
	if i == -1 {
		return fmt.Errorf("could not find asset %q", name)
	}

	// Download asset
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	gh := github.NewClient(http.DefaultClient)
	rc, _, err := gh.Repositories.DownloadReleaseAsset(ctx, "reaper47", "recipya", rel.Assets[i].GetID(), gh.Client())
	if err != nil {
		return err
	}
	defer func() {
		_ = rc.Close()
	}()

	f, err := os.CreateTemp("", "*")
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Remove(f.Name())
	}()

	_, err = io.Copy(f, rc)
	if err != nil {
		_ = f.Close()
		return err
	}
	_ = f.Close()

	tempDir, err := os.MkdirTemp("", "*")
	if err != nil {
		return err
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	err = unzip(f.Name(), tempDir)
	if err != nil {
		return err
	}

	// Install
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	exeBak := exe + ".bak"
	defer func() {
		_ = os.Remove(exeBak)
	}()

	path := filepath.Join(tempDir, name)
	if runtime.GOOS == "windows" {
		err = copyFile(exe, exeBak)
		if err != nil {
			return err
		}

		err = copyFile(path+".exe", exe+".new")
		if err != nil {
			return err
		}
	} else {
		err = os.Rename(exe, exeBak)
		if err != nil {
			return err
		}

		err = os.Rename(path, exe)
		if err != nil {
			return err
		}

		err = os.Chmod(exe, 0775)
		if err != nil {
			return err
		}
	}

	return nil
}

func unzip(src, dest string) error {
	zr, err := zip.OpenReader(src)
	if err != nil {
		return err
	}

	for _, f := range zr.File {
		rcFile, err := f.Open()
		if err != nil {
			return err
		}

		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			err = os.MkdirAll(path, f.Mode())
			if err != nil {
				return err
			}
			continue
		}

		dest, err := os.Create(path)
		if err != nil {
			return err
		}

		_, err = io.Copy(dest, rcFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	f, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, f, 0o644)
}
