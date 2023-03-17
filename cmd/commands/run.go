package commands

import (
	"fmt"
	"github.com/jlewi/hydros/pkg/util"
	"github.com/spf13/cobra"
)

// NewRunCmd creates a command to summarize the commits on a PR.
func NewRunCmd() *cobra.Command {
	var level string
	var jsonLog bool

	cmd := &cobra.Command{
		Short: "Summarize the commit logs of a PR",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			util.SetupLogger(level, !jsonLog)
		},
		Run: func(cmd *cobra.Command, args []string) {
			err := func() error {

				return nil
			}()
			if err != nil {
				fmt.Printf("Failed to summarize PR; error:\n%+v", err)
			}
		},
	}

	//cmd.Flags().StringVarP(&secret, "private-key", "", "", "The uri containing the secret for the GitHub App to Authenticate as. Supported schema file, gcpSecretManager")
	//cmd.Flags().IntVarP(&githubAppID, "appId", "", hydros.HydrosGitHubAppID, "GitHubAppId.")
	//cmd.Flags().StringVarP(&org, "org", "o", "PrimerAI", "The GitHub org to obtain the token for")
	//cmd.Flags().StringVarP(&repo, "repo", "r", "", "The repo obtain the token for")
	//cmd.Flags().StringVarP(&envFile, "env-file", "f", "", "The file to right the github token to")
	//
	//util.IgnoreError(cmd.MarkFlagRequired("private-key"))

	cmd.PersistentFlags().StringVarP(&level, "level", "", "info", "The logging level.")
	cmd.PersistentFlags().BoolVarP(&jsonLog, "json-logs", "", false, "Enable json logging.")
	return cmd
}
