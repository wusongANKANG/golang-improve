package structsinterfaces

import "fmt"

type Notifier interface {
	Notify() string
}

type User struct {
	Name  string
	Email string
}

func (u User) Notify() string {
	return fmt.Sprintf("notify %s via %s", u.Name, u.Email)
}

func (u *User) Rename(name string) {
	u.Name = name
}

type Admin struct {
	User
	Level int
}

func Promote(u User, level int) Admin {
	return Admin{
		User:  u,
		Level: level,
	}
}

func SendAll(notifiers []Notifier) []string {
	messages := make([]string, 0, len(notifiers))
	for _, notifier := range notifiers {
		messages = append(messages, notifier.Notify())
	}

	return messages
}
