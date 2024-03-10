package updater

import (
	"archive/zip"
	"bytes"
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

func IsLatest(current semver.Version) (bool, *github.RepositoryRelease, error) {
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

func Update(current semver.Version) error {
	// Check if latest
	isLatest, rel, err := IsLatest(current)
	if err != nil {
		return err
	}

	if isLatest {
		return errors.New("already latest version")
	}

	// Find asset
	name := fmt.Sprintf("recipya-%s-%s", runtime.GOOS, runtime.GOARCH)
	i := slices.IndexFunc(rel.Assets, func(asset *github.ReleaseAsset) bool { return *asset.Name == name })
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

	tempDir, err := os.MkdirTemp("", "*")
	if err != nil {
		return err
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	err = unzip(rc, tempDir)
	if err != nil {
		return err
	}

	return nil
}

func unzip(rc io.ReadCloser, dest string) error {
	data, err := io.ReadAll(rc)
	if err != nil {
		return err
	}

	zr, err := zip.NewReader(bytes.NewReader(data), 9)
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
