package pkg

import (
	"fmt"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/google/go-cmp/cmp"
	"github.com/jlewi/hydros/pkg/util"
	"os"
	"testing"
)

func Test_build_prompt(t *testing.T) {
	type testCase struct {
		in       []string
		expected string
	}

	cases := []testCase{
		{
			in: []string{"foo", "bar"},
			expected: `Below is a list of commit messages that have not been squashed. Please squash them into the commit above.
Please construct a commit message that summarizes the changes in the commits below. Remove spurious messages 
like "fix typo", "fix lint", "latest", etc... You can use markdown to format the message e.g. to use
lists.

---
foo

---
bar

`,
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("Case %d", i), func(t *testing.T) {
			actual, err := buildPrompt(c.in)

			fmt.Printf("actual:\n%v", actual)
			if err != nil {
				t.Fatalf("Error building prompt; %v", err)
			}

			if d := cmp.Diff(c.expected, actual); d != "" {
				t.Fatalf("Unexpected diff; %v", d)
			}
		})
	}
}

func Test_SummarizeLogMessages(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skip("Test_SummarizeLogMessages doesn't run in GHA because it requires a git repository and branch")
	}
	// An integration test.
	util.SetupLogger("debug", true)

	apiKey := GetAPIKey()
	if apiKey == "" {
		t.Fatalf("Could not locate an OPENAI API key not set")
	}
	client := gpt3.NewClient(string(apiKey))

	summary, err := SummarizeLogMessages(client, "/users/jlewi/git_squash", "origin/main")
	if err != nil {
		t.Fatalf("Error running; %v", err)
	}

	if summary == "" {
		t.Fatalf("Expected summary to be non-empty")
	}

	t.Logf("Summary:\n%v", summary)
}
