package cmd

import (
	"github.com/spf13/cobra"
)

// CreateTestRootCmd creates a root command for testing purposes
func CreateTestRootCmd() *cobra.Command {
	testRootCmd := &cobra.Command{
		Use:   "magecomm",
		Short: "MageComm CLI is a command line tool for managing Magento applications",
	}

	return testRootCmd
}
