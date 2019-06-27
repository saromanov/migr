package app_test

import (
	"os"
	"testing"

	"github.com/saromanov/migr/pkg/app"
)

const (
	basicPath = "../../testdata/basic"
)

var appTest *app.App

func init() {
	os.Setenv("MIGR_PATH", basicPath)
	appTest = app.New("postgres", "pinger", "pinger", "pinger", "pinger", 5432)
}

func TestRun(t *testing.T) {
	err := appTest.Run(basicPath)
	if err != nil {
		t.Errorf("unable to execute run command: %v", err)
	}
}
