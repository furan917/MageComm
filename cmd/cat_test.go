package cmd

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCat(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	testArgs := []string{"cat", "../test_fixtures/cat_cmd/ArchiveTest.tar", "archivetestfile.html"}
	testRootCmd := CreateTestRootCmd()
	testRootCmd.AddCommand(CatCmd)
	testRootCmd.SetArgs(testArgs)
	err := testRootCmd.Execute()
	if err != nil {
		return
	}

	_ = w.Close()
	os.Stdout = oldStdout
	output, _ := io.ReadAll(r)

	assert.True(t, strings.Contains(string(output), "hello world, I am a cat'd file from an archive"))
}

//func TestCatDeploy(t *testing.T) {
//	oldStdout := os.Stdout
//	r, w, _ := os.Pipe()
//	os.Stdout = w
//
//	// Set the deploy archive folder to the test fixtures folder
//	viper.Set(config_manager.CommandConfigDeployArchiveFolder, "/home/francis/magecomm/test_fixtures/cat_deploy_cmd/")
//	testArgs := []string{"cat-deploy", "archivetestfile.html"}
//
//	testRootCmd := CreateTestRootCmd()
//	CatDeployCmd.Args = cobra.ExactArgs(1)
//	testRootCmd.AddCommand(CatDeployCmd)
//	testRootCmd.SetArgs(testArgs)
//	err := testRootCmd.Execute()
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	_ = w.Close()
//	os.Stdout = oldStdout
//	output, _ := io.ReadAll(r)
//
//	assert.True(t, strings.Contains(string(output), "hello world, I am a cat'd file from an archive"))
//}

func TestCatGzip(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	testArgs := []string{"cat", "../test_fixtures/cat_cmd/ArchiveTest.tar.gz", "archivetestfile.html"}
	testRootCmd := CreateTestRootCmd()
	testRootCmd.AddCommand(CatCmd)
	testRootCmd.SetArgs(testArgs)
	err := testRootCmd.Execute()
	if err != nil {
		return
	}

	_ = w.Close()
	os.Stdout = oldStdout
	output, _ := io.ReadAll(r)

	assert.True(t, strings.Contains(string(output), "hello world, I am a cat'd file from an archive"))
}

func TestCatZip(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	testArgs := []string{"cat", "../test_fixtures/cat_cmd/ArchiveTest.zip", "archivetestfile.html"}
	testRootCmd := CreateTestRootCmd()
	testRootCmd.AddCommand(CatCmd)
	testRootCmd.SetArgs(testArgs)
	err := testRootCmd.Execute()
	if err != nil {
		return
	}

	_ = w.Close()
	os.Stdout = oldStdout
	output, _ := io.ReadAll(r)

	assert.True(t, strings.Contains(string(output), "hello world, I am a cat'd file from an archive"))
}

func TestCatB2zip(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	testArgs := []string{"cat", "../test_fixtures/cat_cmd/ArchiveTest.tar.bz2", "archivetestfile.html"}
	testRootCmd := CreateTestRootCmd()
	testRootCmd.AddCommand(CatCmd)
	testRootCmd.SetArgs(testArgs)
	err := testRootCmd.Execute()
	if err != nil {
		return
	}

	_ = w.Close()
	os.Stdout = oldStdout
	output, _ := io.ReadAll(r)

	assert.True(t, strings.Contains(string(output), "hello world, I am a cat'd file from an archive"))
}
