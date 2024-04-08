package server

import (
	"bufio"
	"bytes"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"strconv"
	"sync"
)

// SessionData maps a UUID to a user id. It's used to track who is logged in session-wise.
var SessionData SessionDataMap

// SessionDataMap is a type alias to map UUIDs to int64s.
type SessionDataMap struct {
	Data  map[uuid.UUID]int64
	mutex sync.Mutex
}

// Save saves the SessionDataMap to the writer in the CSV format.
func (s *SessionDataMap) Save(w io.Writer) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for k, v := range s.Data {
		_, _ = w.Write([]byte(k.String() + "," + strconv.FormatInt(v, 10) + "\n"))
	}
}

// Load populates the SessionDataMap from the reader.
func (s *SessionDataMap) Load(r io.Reader) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		parts := bytes.Split(sc.Bytes(), []byte(","))
		if len(parts) != 2 {
			continue
		}

		k, err := uuid.ParseBytes(parts[0])
		if err != nil {
			continue
		}

		v, err := strconv.ParseInt(string(parts[1]), 10, 64)
		if err != nil {
			continue
		}

		s.Data[k] = v
	}

	slog.Info("User sessions restored")
}
