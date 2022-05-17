package jobs

import (
	"testing"
	"testing/fstest"
)

func TestJobsCleanImages(t *testing.T) {
	cases := []struct {
		name                       string
		dir                        fstest.MapFS
		usedFiles                  []string
		wantNumFiles, wantNumBytes int64
	}{
		{
			name: "delete all specified files",
			dir: fstest.MapFS{
				"one":   {Data: []byte("1")},
				"two":   {Data: []byte("22")},
				"three": {Data: []byte("333")},
			},
			usedFiles:    []string{"one", "three"},
			wantNumBytes: 2,
			wantNumFiles: 1,
		},
		{
			name:         "do nothing if no files in folder",
			dir:          fstest.MapFS{},
			usedFiles:    []string{"one", "three"},
			wantNumBytes: 0,
			wantNumFiles: 0,
		},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			gotNumFiles, gotNumBytes := cleanImages(test.dir, test.usedFiles, func(path string) error {
				return nil
			})

			if gotNumFiles != test.wantNumFiles {
				t.Errorf("got %d deleted files but want %d", gotNumFiles, test.wantNumFiles)
			}
			if gotNumBytes != test.wantNumBytes {
				t.Errorf("got %d deleted bytes but want %d", gotNumBytes, test.wantNumBytes)
			}
		})
	}
}
