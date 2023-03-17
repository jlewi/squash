package pkg

import (
	"github.com/jlewi/hydros/pkg/util"
	"testing"
)

func Test_run(t *testing.T) {
	// An integration test.
	util.SetupLogger("debug", true)
	err := Run("/users/jlewi/git_squash", "origin/main")
	if err != nil {
		t.Fatalf("Error running; %v", err)
	}
}
