package server_test

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/server"
	"maps"
	"testing"
)

func TestSessionDataMap_Save(t *testing.T) {
	clear(server.SessionData)
	first := uuid.New()
	second := uuid.New()
	third := uuid.New()
	server.SessionData = server.SessionDataMap{first: 1, second: 2, third: 3}

	var got bytes.Buffer
	server.SessionData.Save(&got)

	want := fmt.Sprintf("%s,1\n%s,2\n%s,3\n", first, second, third)
	if got.String() != want {
		t.Fatalf("got %q but want %q", got, want)
	}
}

func TestSessionDataMap_Load(t *testing.T) {
	first := uuid.New()
	second := uuid.New()
	third := uuid.New()
	server.SessionData = server.SessionDataMap{first: 1, second: 2, third: 3}

	var buf bytes.Buffer
	server.SessionData.Save(&buf)
	clear(server.SessionData)
	server.SessionData.Load(&buf)

	if !maps.Equal(server.SessionData, server.SessionDataMap{first: 1, second: 2, third: 3}) {
		t.Fail()
	}
}
