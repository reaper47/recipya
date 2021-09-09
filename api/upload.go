package api

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
)

// UploadFile retrieves a file from a request and returns its buffer.
func UploadFileBuffer(r *http.Request, maxMB int64, key string) (*bytes.Buffer, error) {
	err := r.ParseMultipartForm(maxMB << 20)
	if err != nil {
		return nil, errors.New(
			"the file's size must be" + strconv.FormatInt(maxMB, 10) + "MB maximum",
		)
	}

	file, _, err := r.FormFile(key)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, err
	}
	return buf, nil
}
