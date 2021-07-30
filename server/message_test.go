package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestMessage_ToJson(t *testing.T) {
	var builder strings.Builder
	builder.WriteString("/help = show all orders\n")
	builder.WriteString("/register [name] = register a new user with name\n")
	builder.WriteString("/rooms = show all chat rooms\n")
	builder.WriteString("/creatRoom [name] = create a chat room with name\n")
	builder.WriteString("/leaveRoom = leave the chat room\n")
	builder.WriteString("/joinRoom [name] = join the chat room by name\n")
	builder.WriteString("/closeRoom = close the room, if the room create by yourself\n")
	builder.WriteString("/users = show all users registered\n")
	builder.WriteString("/popular [n] = show the popular word in recent n seconds, n less than 61\n")
	builder.WriteString("/stats [username] = show the online time of user whose name is username\n")
	msg := NewMessage("system", builder.String())
	//msg := NewMessage("system", "ok")
	bytes, err := msg.ToJson()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("json:", string(bytes))
}

func TestTrim(t *testing.T) {
	str := "          "
	str = strings.Trim(str, " ")
	if str == " " {
		t.Fatal(str)
	}
	fmt.Println("str:", str, len(str))
}
