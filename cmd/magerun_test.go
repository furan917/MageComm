package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testMessagePublisher struct {
	Queue       string
	MessageBody string
}

func (t *testMessagePublisher) Publish(queue string, messageBody string) {
	t.Queue = queue
	t.MessageBody = messageBody
}

func TestMagerunCmd(t *testing.T) {
	testPublisher := &testMessagePublisher{}
	SetMessagePublisher(testPublisher)

	testRootCmd := CreateTestRootCmd()
	testRootCmd.AddCommand(MagerunCmd)
	testArgs := []string{"magerun", "cache:clean"}
	testRootCmd.SetArgs(testArgs)
	err := testRootCmd.Execute()

	assert.NoError(t, err)
	assert.Equal(t, MageRunQueue, testPublisher.Queue)
	assert.Equal(t, "cache:clean", testPublisher.MessageBody)

	SetMessagePublisher(&defaultMessagePublisher{})
}

func TestMagerunCmdBlocksBannedCmd(t *testing.T) {
	testPublisher := &testMessagePublisher{}
	SetMessagePublisher(testPublisher)

	testRootCmd := CreateTestRootCmd()
	testRootCmd.AddCommand(MagerunCmd)
	testArgs := []string{"magerun", "module:enable", "Vendor_Module"}
	testRootCmd.SetArgs(testArgs)
	err := testRootCmd.Execute()

	assert.True(t, strings.Contains(err.Error(), "the command 'module:enable' is not allowed"))

	SetMessagePublisher(&defaultMessagePublisher{})
}