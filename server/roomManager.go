package main

import (
	"errors"
	"sync"
)

type roomManager struct {
	rooms       map[string]*room //所有现有的聊天室
	messagesLen uint             //聊天室保存的历史消息长度
	lock        *sync.Mutex      //保证线程安全
}

// NewRoomManager 创建聊天室管理器
func NewRoomManager(messagesLen uint) *roomManager {
	return &roomManager{
		rooms:       make(map[string]*room),
		messagesLen: messagesLen,
		lock:        new(sync.Mutex),
	}
}

// Rooms 查询现在所有的聊天室
func (rm *roomManager) Rooms() map[string]*room {
	return rm.rooms
}

// NewRoom 创建一个新的聊天室
func (rm *roomManager) NewRoom(name string, owner *user) (error, *room) {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	if _, ok := rm.rooms[name]; ok {
		return errors.New(RoomNameRepeated), nil
	}
	room := NewRoom(name, owner, rm.messagesLen, rm)
	rm.rooms[name] = room
	return nil, room
}

// CloseRoom 关闭聊天室
func (rm *roomManager) CloseRoom(owner *user, room *room) error {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	err := room.Close(owner)
	if err != nil {
		return err
	}
	delete(rm.rooms, room.name)
	return nil
}

// DeleteRoom 删除聊天室
func (rm *roomManager) DeleteRoom(name string) {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	delete(rm.rooms, name)
}

// Join 加入聊天室
func (rm *roomManager) Join(roomName string, user *user) (error, *room) {
	if room, ok := rm.rooms[roomName]; ok {
		err := room.Join(user)
		if err != nil {
			return err, nil
		}
		return nil, room
	} else {
		return errors.New(RoomNotfound), nil
	}
}
