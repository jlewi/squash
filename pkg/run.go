package pkg

import (
	"context"
	"fmt"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"strings"
	"text/template"
	"time"
)

const (
	promptTemplate = `Below is a list of commit messages that have not been squashed. Please squash them into the commit above.
Please construct a commit message that summarizes the changes in the commits below. Remove spurious messages 
like "fix typo", "fix lint", "latest", etc... 
The output should start with a summary of the changes that tries to capture the primary purpose of the changes.
You can use markdown to format the message.

Here are the commit messages:
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
	t, err := template.New("prompt").Parse(promptTemplate)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	if err := t.Execute(&sb, &PromptArgs{Messages: messages}); err != nil {
		return "", err
	}

	return sb.String(), nil
}

// SummarizeLogMessages returns a summary of the log messages on the branch when compared to the base branch.
// path should be the path to a GitHub repository. baseBranch should be the branch from which it was forked.
func SummarizeLogMessages(client gpt3.Client, path string, baseBranch string) (string, error) {
	log := zapr.NewLogger(zap.L())
	r, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{})

	if err != nil {
		return "", err
	}

	headHash, err := r.ResolveRevision("HEAD")

	log.Info("HEAD hash", "hash", headHash)
	if err != nil {
		return "", err
	}
	headCommit, err := r.CommitObject(*headHash)
	if err != nil {
		return "", err
	}

	baseHash, err := r.ResolveRevision(plumbing.Revision(baseBranch))
	if err != nil {
		return "", err
	}
	baseCommit, err := r.CommitObject(*baseHash)
	log.Info("base branch", "hash", baseHash)
	if err != nil {
		return "", err
	}
	commits, err := headCommit.MergeBase(baseCommit)
	if err != nil {
		return "", err
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
		return "", err
	}

	// Reverse the order of messages
	messages = reverseList(messages)

	prompt, err := buildPrompt(messages)
	if err != nil {
		return "", err
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancelFunc()

	// TODO(jeremy): We should really try ChatGPT3 but we have to update the GoClient to support it
	engine := "text-davinci-003"
	resp, err := client.CompletionWithEngine(ctx, engine, gpt3.CompletionRequest{
		Prompt:      []string{prompt},
		MaxTokens:   gpt3.IntPtr(1000),
		Temperature: gpt3.Float32Ptr(0.7),
		TopP:        gpt3.Float32Ptr(1),
		// We currently don't set these to match the playground.
		Stop: []string{},
		Echo: false,
	})
	if err != nil {
		return "", err
	}

	if len(resp.Choices) != 1 {
		return "", fmt.Errorf("expected 1 choice but got %d", len(resp.Choices))
	}
	summary := resp.Choices[0].Text
	log.Info("Summarized log messages", "summary", summary, "prompt", prompt)
	return summary, nil
}
