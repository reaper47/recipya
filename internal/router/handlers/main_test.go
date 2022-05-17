package handlers_test

import (
	"net/http"
	"os"
	"os/exec"
	"testing"

	approvals "github.com/approvals/go-approval-tests"
	"github.com/approvals/go-approval-tests/reporters"
	"github.com/reaper47/recipya/internal/auth"
	"github.com/reaper47/recipya/internal/constants"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/repository"
)

var testCookie *http.Cookie

type goland struct{}

func NewGoland() reporters.Reporter {
	return &goland{}
}

func (s *goland) Report(approved, received string) bool {
	cmd := exec.Command("/opt/GoLand-2022.1.1/bin/goland.sh", approved, received)
	_ = cmd.Start()
	return true
}

func TestMain(t *testing.M) {
	sid, _ := auth.CreateToken("test")
	testCookie = &http.Cookie{Name: constants.CookieSession, Value: sid}
	repository.Sessions["test"] = models.Session{}

	approvals.UseReporter(NewGoland())
	approvals.UseFolder("testdata")

	os.Exit(t.Run())
}
