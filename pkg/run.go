package pkg

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
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

	forkCommit := commits[0]
	log.Info("Merge base", "hash", forkCommit)

	// Relevant issue for how to get all log messages between two commits.
	// https://github.com/go-git/go-git/issues/69

	cIter, err := r.Log(&git.LogOptions{
		From:  headCommit.Hash,
		Order: git.LogOrderCommitterTime,
	})

	if err := cIter.ForEach(func(c *object.Commit) error {
		if c.Hash == forkCommit.Hash {
			// We reached the common ancestor so we can stop.
			return storer.ErrStop
		}
		log.Info("Commit", "hash", c.Hash, "message", c.Message)
		return nil
	}); err != nil {
		return err
	}
	
	return nil
}
