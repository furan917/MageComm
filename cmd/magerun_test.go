package cmd

import (
    "magecomm/messages/publisher"
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
)

type testMessagePublisher struct {
    Queue            string
    MessageBody      string
    AddCorrelationID string
}

func (t *testMessagePublisher) Publish(messageBody string, queue string, AddCorrelationID string) (string, error) {
    t.Queue = queue
    t.MessageBody = messageBody
    t.AddCorrelationID = AddCorrelationID
    return "UniqueCorrelationID", nil
}

func TestMagerunCmd(t *testing.T) {
    testPublisher := &testMessagePublisher{}
    publisher.SetMessagePublisher(testPublisher)

    testRootCmd := CreateTestRootCmd()
    testRootCmd.AddCommand(MagerunCmd)
    testArgs := []string{"magerun", "cache:clean"}
    testRootCmd.SetArgs(testArgs)
    _ = testRootCmd.Execute()

    assert.Equal(t, MageRunQueue, testPublisher.Queue)
    assert.Equal(t, "cache:clean", testPublisher.MessageBody)

    publisher.SetMessagePublisher(&publisher.DefaultMessagePublisher{})
}

func TestMagerunCmdBlocksBannedCmd(t *testing.T) {
    testPublisher := &testMessagePublisher{}
    publisher.SetMessagePublisher(testPublisher)

    testRootCmd := CreateTestRootCmd()
    testRootCmd.AddCommand(MagerunCmd)
    testArgs := []string{"magerun", "module:enable", "Vendor_Module"}
    testRootCmd.SetArgs(testArgs)
    err := testRootCmd.Execute()

    assert.True(t, strings.Contains(err.Error(), "the command 'module:enable' is not allowed"))

    publisher.SetMessagePublisher(&publisher.DefaultMessagePublisher{})
}
