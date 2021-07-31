package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

//user 玩家对象
type user struct {
	name       string        //名字
	login      int64         //登陆时间
	room       *room         //当前所在聊天室
	messageIn  chan *Message //接收消息的频道
	conn       *net.TCPConn  //tcp链接
	lock       *sync.Mutex   //线程安全
	userManger *userManager  //玩家管理器
}

// NewUser 创建一个新的玩家对象
func NewUser(name string, conn *net.TCPConn, userManger *userManager) *user {
	user := &user{
		name:       name,
		login:      time.Now().Unix(),
		messageIn:  make(chan *Message, 10),
		conn:       conn,
		lock:       new(sync.Mutex),
		userManger: userManger,
	}
	go user.Listen()
	if conn != nil {
		user.rooms()
	}
	return user
}

// Disconnect 玩家断开链接
func (u *user) Disconnect() {
	u.lock.Lock()
	defer u.lock.Unlock()
	// 离开聊天室
	if u.room != nil {
		u.room.UserLeave(u.name)
	}
	//关闭消息通道
	close(u.messageIn)
	//注销玩家信息
	u.userManger.Unregister(u.name)
	fmt.Println("user [" + u.name + "] leave server")
}

// Broadcast 广播消息
func (u *user) Broadcast(message *Message) {
	u.messageIn <- message
}

// RoomClosed 离开聊天室
func (u *user) RoomClosed() {
	msg := NewMessage("system", "chat room closed")
	u.Broadcast(msg)
	u.room = nil
}

// Listen 监听广播消息
func (u *user) Listen() {
	for {
		select {
		case msg, ok := <-u.messageIn:
			if !ok {
				return
			}
			bytes, err := msg.ToJson()
			if err != nil {
				fmt.Println(MessageEncodeError, err.Error())
			} else if u.conn == nil {
				fmt.Println("send message error: conn is nil")
			} else {
				_, err = u.conn.Write(bytes)
				if err != nil {
					fmt.Println("send message error:", err.Error())
				}
			}
		}
	}
}

// Send 发送群消息
func (u *user) Send(message *Message) {
	u.room.messageIn <- message
}

// MessageRouter 消息路由，处理所有的GM命令和非GM命令消息
func (u *user) MessageRouter(message *Message) {
	u.lock.Lock()
	defer u.lock.Unlock()
	if !message.isGm { //非GM命令
		PopularManager.text <- message.Text //统计单词频次
		u.Send(message)
	} else { //GM命令
		switch message.gmOrder[0] {
		case Register:
			msg := NewMessage("system", AlreadyRegistered)
			bytes, _ := msg.ToJson()
			u.conn.Write(bytes)
		case Rooms:
			u.rooms()
		case CreateRoom:
			u.createRoom(message)
		case LeaveRoom:
			u.leaveRoom()
		case JoinRoom:
			u.joinRoom(message)
		case CloseRoom:
			u.closeRoom()
		case Users:
			u.users()
		case Popular:
			u.popular(message)
		case Stats:
			u.stats(message)
		default:
			msg := NewMessage("system", UnknownOrder)
			bytes, _ := msg.ToJson()
			u.conn.Write(bytes)
		}
	}
}

func (u *user) rooms() {
	var builder strings.Builder
	for k, _ := range RoomManager.rooms {
		builder.WriteString(k)
		builder.WriteString("\n")
	}
	msg := NewMessage("system", builder.String())
	bytes, _ := msg.ToJson()
	u.conn.Write(bytes)
}

func (u *user) createRoom(message *Message) {
	if len(message.gmOrder) < 2 {
		msg := NewMessage("system", "command error, please input /createRoom name")
		bytes, _ := msg.ToJson()
		u.conn.Write(bytes)
		return
	}
	err, room := RoomManager.NewRoom(message.gmOrder[1], u)
	str := "chat room create success"
	if err != nil {
		str = err.Error()
	} else {
		u.room = room
	}
	msg := NewMessage("system", str)
	bytes, _ := msg.ToJson()
	u.conn.Write(bytes)
}

func (u *user) leaveRoom() {
	u.room.UserLeave(u.name)
	u.room = nil
	msg := NewMessage("system", "leave room success")
	bytes, _ := msg.ToJson()
	u.conn.Write(bytes)
}

func (u *user) joinRoom(message *Message) {
	if len(message.gmOrder) < 2 {
		msg := NewMessage("system", "command error, please input /joinRoom name")
		bytes, _ := msg.ToJson()
		u.conn.Write(bytes)
		return
	}
	err, room := RoomManager.Join(message.gmOrder[1], u)
	if err != nil {
		str := err.Error()
		msg := NewMessage("system", str)
		bytes, _ := msg.ToJson()
		u.conn.Write(bytes)
	} else {
		u.room = room
		for _, msg := range room.messages {
			bytes, _ := msg.ToJson()
			u.conn.Write(bytes)
		}
	}
}

func (u *user) closeRoom() {
	str := "chat room close success"
	if u.room != nil {
		err := u.room.Close(u)
		if err != nil {
			str = err.Error()
		} else {
			u.room = nil
		}
	}
	msg := NewMessage("system", str)
	bytes, _ := msg.ToJson()
	u.conn.Write(bytes)
}

func (u *user) users() {
	var builder strings.Builder
	for k, _ := range u.userManger.users {
		builder.WriteString(k)
		builder.WriteString("\n")
	}
	msg := NewMessage("system", builder.String())
	bytes, _ := msg.ToJson()
	u.conn.Write(bytes)
}

func (u *user) popular(message *Message) {
	if len(message.gmOrder) < 2 {
		msg := NewMessage("system", "command error, please input /popular n")
		bytes, _ := msg.ToJson()
		u.conn.Write(bytes)
		return
	}
	seconds := message.gmOrder[1]
	i, err := strconv.ParseInt(seconds, 10, 64)
	var text string
	if err != nil {
		text = ParamError
	} else {
		text = PopularManager.Count(i)
	}
	msg := NewMessage("system", text)
	bytes, _ := msg.ToJson()
	u.conn.Write(bytes)
}

func (u *user) stats(message *Message) {
	if len(message.gmOrder) < 2 {
		msg := NewMessage("system", "command error, please input /stats name")
		bytes, _ := msg.ToJson()
		u.conn.Write(bytes)
		return
	}
	now := time.Now().Unix()
	user, find := u.userManger.users[message.gmOrder[1]]
	if !find {
		msg := NewMessage("system", "user not found")
		bytes, _ := msg.ToJson()
		u.conn.Write(bytes)
		return
	}
	delta := now - user.login
	t := SecondsToDayStr(delta)
	msg := NewMessage("system", t)
	bytes, _ := msg.ToJson()
	u.conn.Write(bytes)
}
