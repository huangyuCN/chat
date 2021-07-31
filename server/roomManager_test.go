package main

import "testing"

func TestNewRoomManager(t *testing.T) {
	RoomManager = NewRoomManager(HistoryMessageLen)
	if RoomManager.messagesLen != HistoryMessageLen || len(RoomManager.rooms) != 0 {
		t.Fatal("error")
	}
}

func TestRoomManager_NewRoom(t *testing.T) {
	RoomManager = NewRoomManager(HistoryMessageLen)
	UserManager = NewUserManager()
	err, owner := UserManager.Register("username", nil)
	if err != nil {
		t.Fatal("error")
	}
	err, room := RoomManager.NewRoom("room", owner)
	if err != nil {
		t.Fatal("error")
	}
	if room.owner != owner {
		t.Fatal("error")
	}
}

func TestRoomManager_Rooms(t *testing.T) {
	RoomManager = NewRoomManager(HistoryMessageLen)
	UserManager = NewUserManager()
	err, owner := UserManager.Register("username", nil)
	if err != nil {
		t.Fatal("error")
	}
	err, room := RoomManager.NewRoom("room", owner)
	if err != nil {
		t.Fatal("error")
	}
	if RoomManager.rooms["room"] != room {
		t.Fatal("error")
	}
}

func TestRoomManager_DeleteRoom(t *testing.T) {
	RoomManager = NewRoomManager(HistoryMessageLen)
	UserManager = NewUserManager()
	err, owner := UserManager.Register("username", nil)
	if err != nil {
		t.Fatal("error")
	}
	err, _ = RoomManager.NewRoom("room", owner)
	if err != nil {
		t.Fatal("error")
	}
	RoomManager.DeleteRoom("room")
	if RoomManager.rooms["room"] != nil {
		t.Fatal("error")
	}
}

func TestRoomManager_Join(t *testing.T) {
	RoomManager = NewRoomManager(HistoryMessageLen)
	UserManager = NewUserManager()
	err, owner := UserManager.Register("username", nil)
	if err != nil {
		t.Fatal("error")
	}
	err, _ = RoomManager.NewRoom("room", owner)
	if err != nil {
		t.Fatal("error")
	}
	err, user := UserManager.Register("username1", nil)
	if err != nil {
		t.Fatal("error")
	}
	err, _ = RoomManager.Join("room", user)
	if err != nil {
		t.Fatal("error")
	}
}
