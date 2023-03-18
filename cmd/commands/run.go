package commands

import (
	"fmt"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/go-errors/errors"
	"github.com/jlewi/hydros/pkg/util"
	"github.com/jlewi/squash/pkg"
	"github.com/spf13/cobra"
)

// NewRunCmd creates a command to summarize the commits on a PR.
func NewRunCmd() *cobra.Command {
	var level string
	var jsonLog bool

	var path string
	var base string

	cmd := &cobra.Command{
		Short: "Summarize the commit logs of a PR",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			util.SetupLogger(level, !jsonLog)
		},
		Run: func(cmd *cobra.Command, args []string) {
			err := func() error {
				apiKey := pkg.GetAPIKey()
				if apiKey == "" {
					return errors.New("Could not locate an OPENAI API key not set")
				}
				client := gpt3.NewClient(string(apiKey))
				summary, err := pkg.SummarizeLogMessages(client, path, base)
				if err != nil {
					return err
				}
				fmt.Printf("Summary:\n%s\n", summary)
				return nil
			}()
			if err != nil {
				fmt.Printf("Failed to summarize PR; error:\n%+v", err)
			}
		},
	}

	cmd.Flags().StringVarP(&path, "path", "", "", "The path of the repository to summarize the commits of. HEAD should be pointing to the branch to get the commits of")
	cmd.Flags().StringVarP(&base, "base", "", "origin/main", "The base branch to compare against. It should have a common ancestor with the current branch. It will be used to compute the fork point.")
	cmd.MarkFlagRequired("path")
	cmd.PersistentFlags().StringVarP(&level, "level", "", "info", "The logging level.")
	cmd.PersistentFlags().BoolVarP(&jsonLog, "json-logs", "", false, "Enable json logging.")
	return cmd
}
