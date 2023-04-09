package services

type MessageService interface {
	Connect() error
	Disconnect() error
	Refresh() error
}
