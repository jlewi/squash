package main

import (
	"fmt"
	"github.com/jlewi/squash/cmd/commands"
	"os"
)

func newRootCmd() *cobra.Command {
	var level string
	var jsonLog bool
	rootCmd := &cobra.Command{
		Short: "agent",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			_, err := logging.InitLogger(level, !jsonLog)
			if err != nil {
				panic(err)
			}
		},
	}

	rootCmd.PersistentFlags().StringVarP(&level, "level", "", "info", "The logging level.")
	rootCmd.PersistentFlags().BoolVarP(&jsonLog, "json-logs", "", false, "Enable json logging.")
	return rootCmd
}

func main() {
	rootCmd := commands.NewRunCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Command failed with error: %+v", err)
		os.Exit(1)
	}
}
