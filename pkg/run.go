package pkg

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"strings"
	"text/template"
)

const (
	prompt = `Below is a list of commit messages that have not been squashed. Please squash them into the commit above.
Please construct a commit message that summarizes the changes in the commits below. Remove spurious messages 
like "fix typo", "fix lint", "latest", etc... You can use markdown to format the message e.g. to use
lists.
{{range .Messages}}
---
{{ . }}
{{end}}
`
)

type PromptArgs struct {
	Messages []string
}

func reverseList(s []string) []string {
	n := len(s)
	ret := make([]string, n)
	for i := 0; i < n; i++ {
		ret[n-i-1] = s[i]
	}
	return ret
}

func buildPrompt(messages []string) (string, error) {
	t, err := template.New("prompt").Parse(prompt)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	if err := t.Execute(&sb, &PromptArgs{Messages: messages}); err != nil {
		return "", err
	}

	return sb.String(), nil
}
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

	messages := make([]string, 0, 20)
	if err := cIter.ForEach(func(c *object.Commit) error {
		if c.Hash == forkCommit.Hash {
			// We reached the common ancestor so we can stop.
			return storer.ErrStop
		}
		log.Info("Commit", "hash", c.Hash, "message", c.Message)
		messages = append(messages, c.Message)
		return nil
	}); err != nil {
		return err
	}

	// Reverse the order of messages
	messages = reverseList(messages)

	return nil
}
