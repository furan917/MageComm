package notifictions

type Notifier interface {
	Notify(message string) error
}
