package cmd

// this is a basic cat command that will extract a file from an archive and print it to stdout

import (
	"github.com/spf13/cobra"
	"magecomm/archive"
	"magecomm/logger"
)

// CatCmd extracts a file from an archive and prints its contents to stdout
var CatCmd = &cobra.Command{
	Use:   "cat [archive] [filepath]",
	Short: "Cat a file from an archive",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		archivePath := args[0]
		filePath := args[1]

		err := archive.CatFileFromArchive(archivePath, filePath)
		if err != nil {
			logger.Fatal(err)
		}
		return
	},
}
