package server_test

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/server"
	"maps"
	"slices"
	"sort"
	"strings"
	"testing"
)

func TestSessionDataMap_Save(t *testing.T) {
	clear(server.SessionData.Data)
	first := uuid.New()
	second := uuid.New()
	third := uuid.New()
	server.SessionData = server.SessionDataMap{Data: map[uuid.UUID]int64{first: 1, second: 2, third: 3}}

	var got bytes.Buffer
	server.SessionData.Save(&got)

	split := slices.DeleteFunc(strings.Split(got.String(), "\n"), func(s string) bool { return s == "" })
	want := []string{first.String() + ",1", second.String() + ",2", third.String() + ",3"}
	sort.Strings(want)
	sort.Strings(split)
	if !slices.Equal(want, split) {
		t.Fatalf("got %q but want %q", split, want)
	}
}

func TestSessionDataMap_Load(t *testing.T) {
	first := uuid.New()
	second := uuid.New()
	third := uuid.New()
	server.SessionData = server.SessionDataMap{Data: map[uuid.UUID]int64{first: 1, second: 2, third: 3}}

	var buf bytes.Buffer
	server.SessionData.Save(&buf)
	clear(server.SessionData.Data)
	server.SessionData.Load(&buf)

	if !maps.Equal(server.SessionData.Data, map[uuid.UUID]int64{first: 1, second: 2, third: 3}) {
		t.Fail()
	}
}
