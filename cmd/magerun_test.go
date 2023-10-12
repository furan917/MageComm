package cmd

import (
	"magecomm/messages/publisher"
	"os"
	"strings"
	"testing"
	"time"

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

	// Create a pipe to simulate stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin = r

	testRootCmd := CreateTestRootCmd()
	testRootCmd.AddCommand(MagerunCmd)
	testArgs := []string{"magerun", "cache:clean"}
	testRootCmd.SetArgs(testArgs)

	done := make(chan struct{})
	go func() {
		_ = testRootCmd.Execute()
		close(done)
	}()

	// Simulate stop command for listening to return value
	time.Sleep(1 * time.Second)
	_, _ = w.Write([]byte("\n"))
	<-done

	assert.Equal(t, MageRunQueue, testPublisher.Queue)
	assert.Equal(t, "cache:clean", testPublisher.MessageBody)

	publisher.SetMessagePublisher(&publisher.DefaultMessagePublisher{})
}

func TestMagerunCmdBlocksBannedCmd(t *testing.T) {
	testPublisher := &testMessagePublisher{}
	publisher.SetMessagePublisher(testPublisher)

	// Create a pipe to simulate stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin = r

	testRootCmd := CreateTestRootCmd()
	testRootCmd.AddCommand(MagerunCmd)
	testArgs := []string{"magerun", "module:enable", "Vendor_Module"}
	testRootCmd.SetArgs(testArgs)

	done := make(chan struct{})
	var executeErr error
	go func() {
		executeErr = testRootCmd.Execute()
		close(done)
	}()

	// Simulate stop command for listening to return value
	time.Sleep(1 * time.Second)
	_, _ = w.Write([]byte("\n"))
	<-done

	assert.NotNil(t, executeErr)
	assert.True(t, strings.Contains(executeErr.Error(), "`module:enable` Command not allowed"))

	publisher.SetMessagePublisher(&publisher.DefaultMessagePublisher{})
}
