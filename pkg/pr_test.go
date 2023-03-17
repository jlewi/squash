package pkg

import (
	"github.com/go-logr/zapr"
	"github.com/jlewi/hydros/pkg/files"
	"github.com/jlewi/hydros/pkg/github"
	"github.com/jlewi/hydros/pkg/github/ghrepo"
	"github.com/jlewi/hydros/pkg/util"
	"go.uber.org/zap"
	"io"
	"os"
	"path/filepath"
	"testing"
)

const (
	appID      = int64(266158)
	privateKey = "secrets/hydros-bot.2022-11-27.private-key.pem"
)

func getTransportManager() (*github.TransportManager, error) {
	log := zapr.NewLogger(zap.L())

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	hydrosKeyFile := filepath.Join(home, privateKey)

	f := &files.Factory{}
	h, err := f.Get(privateKey)
	if err != nil {
		return nil, err
	}
	r, err := h.NewReader(hydrosKeyFile)
	if err != nil {
		return nil, err
	}
	secretByte, err := io.ReadAll(r)

	if err != nil {
		return nil, err
	}

	return github.NewTransportManager(appID, secretByte, log)
}

// Test_GetPRMergePoint is an integration test.
func Test_GetPRMergePoint(t *testing.T) {
	// This test verifies that we can check out a repository to a clean directory
	util.SetupLogger("debug", true)

	trManager, err := getTransportManager()
	if err != nil {
		t.Fatalf("Error creating RepoHelper; %v", err)
	}

	org := "jlewi"
	repo := "roboweb"
	tr, err := trManager.Get(org, repo)

	args := &github.RepoHelperArgs{
		BaseRepo:   ghrepo.New(org, repo),
		GhTr:       tr,
		Name:       "notset",
		Email:      "notset@acme.com",
		BaseBranch: "main",
		BranchName: "jlewi/js",
	}

	rh, err := github.NewGithubRepoHelper(args)
	if err != nil {
		t.Fatalf("Error creating RepoHelper; %v", err)
	}

	if err := GetPRBranchPoint(rh, 1); err != nil {
		t.Fatalf("Error checking out PR branch; %v", err)
	}
}
