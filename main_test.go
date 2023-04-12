package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestMainCli(t *testing.T) {
	// Save original args and replace with test args
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Redirect output to a buffer
	outputBuffer := new(bytes.Buffer)
	cobra.MousetrapHelpText = ""

	// Set test args and run main()
	testArgs := []string{"magecomm", "help"}
	os.Args = testArgs
	RootCmd.SetOut(outputBuffer)
	main()

	// Verify that the output contains the expected help text
	output := outputBuffer.String()
	assert.True(t, strings.Contains(output, "MageComm CLI is a command line tool for managing Magento applications"))
}
