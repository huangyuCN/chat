package main

import "testing"

func newRoom() *room {
	UserManager = NewUserManager()
	RoomManager = NewRoomManager(HistoryMessageLen)
	err, owner := UserManager.Register("username", nil)
	if err != nil {
		panic(err)
	}
	room := NewRoom("test", owner, 50, RoomManager)
	return room
}

func TestRoom_Join(t *testing.T) {
	room := newRoom()
	newUser := NewUser("username1", nil, nil)
	err := room.Join(newUser)
	if err != nil || len(room.users) != 2 {
		t.Fatal("error")
	}
}

func TestRoom_UserLeave(t *testing.T) {
	room := newRoom()
	newUser := NewUser("username1", nil, nil)
	err := room.Join(newUser)
	if err != nil || len(room.users) != 2 {
		t.Fatal("error")
	}
	room.UserLeave("username")
	if len(room.users) != 1 {
		t.Fatal("error")
	}
}

func TestRoom_Close(t *testing.T) {
	room := newRoom()
	newUser := NewUser("username1", nil, nil)
	err := room.Join(newUser)
	if err != nil || len(room.users) != 2 {
		t.Fatal("error")
	}
	owner := room.users["username"]
	err = room.Close(owner)
	if err != nil {
		t.Fatal("error")
	}
}
