package main

import (
	"errors"
	"sync"
	"time"
)

type room struct {
	name        string           //聊天室名字
	users       map[string]*user //聊天成员
	createTime  int64            //创建时间
	owner       *user            //创建者
	messageIn   chan *Message    //接收消息的信道
	messages    []*Message       //历史消息
	messagesLen uint             //历史消息长度
	roomManager *roomManager     //聊天室管理器
	lock        *sync.Mutex      //保证线程安全
}

//NewRoom 创建一个新的聊天室
func NewRoom(name string, owner *user, messageLen uint, roomManager *roomManager) *room {
	users := make(map[string]*user)
	users[owner.name] = owner
	room := &room{
		name:        name,
		users:       users,
		createTime:  time.Now().Unix(),
		owner:       owner,
		messageIn:   make(chan *Message, 10),
		messages:    make([]*Message, messageLen, messageLen),
		roomManager: roomManager,
		lock:        new(sync.Mutex),
	}
	go room.listen()
	return room
}

// Close 关闭一个聊天室
func (room *room) Close(owner *user) error {
	room.lock.Lock()
	defer room.lock.Unlock()
	if owner != room.owner {
		return errors.New("permission denied")
	}
	close(room.messageIn) //关闭信道
	for _, u := range room.users {
		u.RoomClosed()
	}
	room.roomManager.DeleteRoom(room.name)
	return nil
}

// Broadcast 发送消息
func (room *room) Broadcast(message *Message) {
	for _, user := range room.users {
		user.Broadcast(message)
	}
}

// listen 监听消息
func (room *room) listen() {
	for {
		select {
		case msg, ok := <-room.messageIn:
			if !ok {
				return
			}
			room.Broadcast(msg)
		}
	}
}

// Join 用户加入聊天室
func (room *room) Join(user *user) error {
	room.lock.Lock()
	defer room.lock.Unlock()
	if _, ok := room.users[user.name]; ok {
		return errors.New("already join")
	}
	room.users[user.name] = user
	msg := NewMessage("system", user.name+" join in the chat")
	room.messageIn <- msg
	return nil
}

// UserLeave 玩家主动离开聊天室
func (room *room) UserLeave(userName string) {
	room.lock.Lock()
	defer room.lock.Unlock()
	msg := NewMessage("system", userName+" leave the chat room")
	room.messageIn <- msg
	delete(room.users, userName)
	if len(room.users) == 0 {
		room.roomManager.DeleteRoom(room.name)
	}
	return
}
