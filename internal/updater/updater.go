package updater

import (
	"context"
	"fmt"
	"github.com/blang/semver"
	"github.com/google/go-github/v59/github"
	"golang.org/x/exp/slices"
	"net/http"
	"runtime"
	"strings"
	"time"
)

type Updater struct {
	current  semver.Version
	ghAssets []*github.ReleaseAsset
	ghClient *github.Client
}

func New(current semver.Version) *Updater {
	return &Updater{
		current:  current,
		ghClient: github.NewClient(http.DefaultClient),
	}
}

func (u *Updater) Check() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	rel, res, err := u.ghClient.Repositories.GetLatestRelease(ctx, "reaper47", "recipya")
	if err != nil {
		return false, err
	}

	if res.StatusCode != http.StatusOK {
		return false, fmt.Errorf("got status code %d instead of 200 when fetching latest releases on GitHub", res.StatusCode)
	}

	version, err := semver.Parse(strings.Replace(rel.GetTagName(), "v", "", 1))
	if err != nil {
		return false, err
	}

	u.ghAssets = rel.Assets
	return version.LTE(u.current), nil
}

func (u *Updater) Update() error {
	name := fmt.Sprintf("recipya-%s-%s", runtime.GOOS, runtime.GOARCH)
	i := slices.IndexFunc(u.ghAssets, func(asset *github.ReleaseAsset) bool { return *asset.Name == name })
	if i == -1 {
		return fmt.Errorf("could not find asset %q", name)
	}
	asset := u.ghAssets[i]

	return nil
}
