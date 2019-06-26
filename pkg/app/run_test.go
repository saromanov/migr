package app

import (
	"testing"
)


const (
	basicPath = "../testdata"
)
func TestRun(t *testing.T) {
	err := Run(basicPath)
	if err != nil {
		t.Errorf("unable to execute run command")
	}
}