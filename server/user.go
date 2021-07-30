package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

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
	user.rooms()
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
	fmt.Println("玩家[" + u.name + "]离开服务器")
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
				fmt.Println("编码错误", err.Error())
			} else {
				_, err = u.conn.Write(bytes)
				if err != nil {
					fmt.Println("发送消息错误", err.Error())
				}
			}
		}
	}
}

// Send 发送群消息
func (u *user) Send(message *Message) {
	u.room.messageIn <- message
}

// MessageRouter 消息路由
func (u *user) MessageRouter(message *Message) {
	u.lock.Lock()
	defer u.lock.Unlock()
	if !message.isGm {
		PopularManager.text <- message.Text //统计单词频次
		u.Send(message)
	} else {
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
			u.stats()
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
	if len(message.gmOrder) < 2{
		msg := NewMessage("system", "order error, please user /createRoom name")
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
	if len(message.gmOrder) < 2{
		msg := NewMessage("system", "order error, please user /joinRoom name")
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
	if len(message.gmOrder) < 2{
		msg := NewMessage("system", "order error, please user /popular n")
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

func (u *user) stats() {
	now := time.Now().Unix()
	delta := now - u.login
	t := SecondsToDayStr(delta)
	msg := NewMessage("system", t)
	bytes, _ := msg.ToJson()
	u.conn.Write(bytes)
}
