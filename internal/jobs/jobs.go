package jobs

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/reaper47/recipya/internal/config"
)

// ScheduleCronJobs schedules cron jobs for the web app. It starts the following jobs:
//
// - Clean Images: Removes unreferenced images from the data/img folder.
func ScheduleCronJobs() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).MonthLastDay().Do(cleanImages)
	s.StartAsync()
}

func cleanImages() {
	dir := "./data/img"

	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(err)
	}

	// TODO: Use https://pkg.go.dev/golang.org/x/exp/slices#IndexFunc when upgrading to Go 1.18
	dbImages := config.App().Repo.Images()

	m := make(map[string]int64)
	for _, v := range fs {
		m[v.Name()] = v.Size()
	}

	for _, v := range dbImages {
		_, ok := m[v]
		if ok {
			delete(m, v)
		}
	}

	var size int64
	for k, v := range m {
		err := os.Remove(dir + "/" + k)
		if err == nil {
			size += v
		}
	}

	var s string
	if size > 0 {
		s = "(" + strconv.FormatFloat(float64(size)/(1<<20), 'f', 2, 64) + " MB)"
	}
	log.Println("CleanImages: Removed " + strconv.Itoa(len(m)) + " unreferenced images " + s)
}
