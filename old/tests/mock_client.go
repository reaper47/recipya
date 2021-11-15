package tests

import "net/http"

// MockClient is the mock for the api.Client client.
type MockClient struct {
	GetFunc func(url string) (resp *http.Response, err error)
}

var (
	// GetGetFunc fetches the mock client's `GetFunc` function.
	GetGetFunc func(url string) (resp *http.Response, err error)
)

// Get calls GetGetFunc. It lets the user define the MockClient.GetFunc when testing.
func (m *MockClient) Get(url string) (resp *http.Response, err error) {
	return GetGetFunc(url)
}
