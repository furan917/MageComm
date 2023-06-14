package handler

import (
	"fmt"
	"magecomm/logger"
)

type MagerunOutputHandler struct {
}

func (handler *MagerunOutputHandler) ProcessMessage(messageBody string, correlationID string) error {
	if messageBody == "" {
		return fmt.Errorf("message body is empty")
	}
	fmt.Println(messageBody)
	logger.Infof("Message body: %s", messageBody)

	return nil
}
