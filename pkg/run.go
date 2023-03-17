package pkg

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

func Run(path string, baseBranch string) error {
	log := zapr.NewLogger(zap.L())
	r, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{})

	if err != nil {
		return err
	}

	headHash, err := r.ResolveRevision("HEAD")

	log.Info("HEAD hash", "hash", headHash)
	if err != nil {
		return err
	}
	headCommit, err := r.CommitObject(*headHash)
	if err != nil {
		return err
	}

	baseHash, err := r.ResolveRevision(plumbing.Revision(baseBranch))
	if err != nil {
		return err
	}
	baseCommit, err := r.CommitObject(*baseHash)
	log.Info("base branch", "hash", baseHash)
	if err != nil {
		return err
	}
	commits, err := headCommit.MergeBase(baseCommit)
	if err != nil {
		return err
	}
	if len(commits) != 1 {
		log.Info("Warning expected 1 commit but got many")
	}

	log.Info("Merge base", "hash", commits[0].Hash)
	
	return nil
}
