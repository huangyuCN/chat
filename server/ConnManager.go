package main

import (
	"bufio"
	"chat/codec"
	"fmt"
	"net"
	"strings"
)

func NewConnManager(conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()
	defer func() {
		fmt.Println("断开连接：" + ipStr)
		//移除连接
		conn.Close()
	}()
	reader := bufio.NewReader(conn)
	var user *user
	for {
		//		Message, err := reader.ReadString('\n')
		bytes, err := codec.Decode(reader)
		if err != nil {
			fmt.Println("解析消息错误", err)
			if user != nil {
				if user.room != nil {
					user.leaveRoom()
				}
				UserManager.Unregister(user.name)
			}
			return
		}
		//fmt.Println(manager.conn.RemoteAddr().String() + ":" + string(Message))
		message := ToMessage(bytes)
		if user != nil {
			message.From = user.name
		}
		if message.isGm == true && message.gmOrder[0] == Help {
			var builder strings.Builder
			builder.WriteString("/help = show all orders\n")
			builder.WriteString("/register [name] = register a new user with name\n")
			builder.WriteString("/rooms = show all chat rooms\n")
			builder.WriteString("/createRoom [name] = create a chat room with name\n")
			builder.WriteString("/leaveRoom = leave the chat room\n")
			builder.WriteString("/joinRoom [name] = join the chat room by name\n")
			builder.WriteString("/closeRoom = close the room, if the room create by yourself\n")
			builder.WriteString("/users = show all users registered\n")
			builder.WriteString("/popular [n] = show the popular word in recent n seconds, n less than 61\n")
			builder.WriteString("/stats [username] = show the online time of user whose name is username\n")
			msg := NewMessage("system", builder.String())
			bytes, err := msg.ToJson()
			if err != nil {
				fmt.Println("解析消息错误", err)
			} else {
				conn.Write(bytes)
			}
			continue
		}

		//如果没有注册，需要先注册
		if user == nil && (message.isGm == false || message.gmOrder[0] != Register) {
			msg := NewMessage("system", NeedRegister)
			bytes, err := msg.ToJson()
			if err != nil {
				fmt.Println("解析消息错误", err)
			} else {
				conn.Write(bytes)
			}
			continue
		}
		//注册一个新用户，并且返回所有的聊天室列表
		if user == nil && message.isGm == true && message.gmOrder[0] == Register {
			if len(message.gmOrder) < 2 {
				msg := NewMessage("system", "order error, please user /register name")
				bytes, _ := msg.ToJson()
				conn.Write(bytes)
			} else {
				var regErr error
				regErr, user = UserManager.Register(message.gmOrder[1], conn)
				if regErr != nil {
					fmt.Println("register error:", regErr)
					msg := NewMessage("system", "register failed:"+regErr.Error())
					bytes, _ := msg.ToJson()
					conn.Write(bytes)
				}
			}
			continue
		}
		// 必须要加入房间之后才能聊天
		if user.room == nil && message.isGm == false {
			msg := NewMessage("system", NeedJoinRoom)
			bytes, err := msg.ToJson()
			if err != nil {
				fmt.Println("解析消息错误", err)
			} else {
				conn.Write(bytes)
			}
			continue
		}
		//消息路由
		user.MessageRouter(message)
	}
}
