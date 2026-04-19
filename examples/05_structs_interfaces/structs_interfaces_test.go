package structsinterfaces

import "testing"

func TestRename(t *testing.T) {
	user := User{Name: "alice", Email: "alice@example.com"}
	user.Rename("bob")

	if user.Name != "bob" {
		t.Fatalf("user.Name = %q, want %q", user.Name, "bob")
	}
}

func TestSendAll(t *testing.T) {
	user := User{Name: "alice", Email: "alice@example.com"}
	admin := Promote(User{Name: "root", Email: "root@example.com"}, 10)

	messages := SendAll([]Notifier{user, admin})

	if len(messages) != 2 {
		t.Fatalf("len(messages) = %d, want 2", len(messages))
	}

	if messages[0] != "notify alice via alice@example.com" {
		t.Fatalf("messages[0] = %q", messages[0])
	}
}
