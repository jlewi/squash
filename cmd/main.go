package main

import (
	"fmt"
	"github.com/jlewi/squash/cmd/commands"
	"os"
)

func main() {
	rootCmd := commands.NewRunCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Command failed with error: %+v", err)
		os.Exit(1)
	}
}
