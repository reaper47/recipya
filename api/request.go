package api

import (
	"net/http"
	"sync"

	"github.com/reaper47/recipya/model"
)

// Client is the global variable for the HTTPClient interface.
var Client HTTPClient

// HTTPClient is an interface for the http package.
//
// It makes testing easier because we can mock.
type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}

// GetAsync is an asynchronous function that calls the specified URL and sends the response back to the channel.
func GetAsync(url string, c chan *model.HttpResponse, wg *sync.WaitGroup) {
	defer (*wg).Done()
	res, err := Client.Get(url)
	if err != nil {
		res.Body.Close()
	}
	c <- &model.HttpResponse{Url: url, Response: res, Err: err}
}

func init() {
	Client = &http.Client{}
}
