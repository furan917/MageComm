package cmd

// this is a basic cat command that will extract a file from an archive and print it to stdout

import (
	"github.com/spf13/cobra"
	"magecomm/archive"
	"magecomm/logger"
)

// CatDeployCmd extracts a file from the latest archived deploy and prints its contents to stdout
var CatDeployCmd = &cobra.Command{
	Use:   "cat-deploy [filepath]",
	Short: "Cat a file from the latest deploy",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		err := archive.CatFileFromDeploy(filePath)
		if err != nil {
			logger.Fatal(err)
		}
		return
	},
}
