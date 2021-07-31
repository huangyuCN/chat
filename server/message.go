package main

import (
	"chat/codec"
	"encoding/json"
	"strings"
	"time"
)

type Message struct {
	From    string   //发送者
	Time    int64    //发送时间
	Text    string   //消息内容
	isGm    bool     //是否为GM命令
	gmOrder []string //GM命令参数
}

// NewMessage 新消息
func NewMessage(from string, text string) *Message {
	return &Message{
		from,
		time.Now().Unix(),
		text,
		false,
		nil,
	}
}

// ToJson 将消息对象转换成json格式
func (msg *Message) ToJson() ([]byte, error) {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return codec.Encode(string(bytes))
}

// ToMessage 将json格式转化成消息对象
func ToMessage(data []byte) *Message {
	var msg Message
	str := string(data)
	msg.Text = Filter(str)
	msg.Time = time.Now().Unix()
	IsGm := IsGm(str)
	msg.isGm = IsGm
	if IsGm && len(str) > 0 {
		msg.gmOrder = make([]string, 0)
		strList := strings.Split(str, " ")
		for _, s := range strList {
			s = strings.Trim(s, " ")
			if len(s) > 0 {
				msg.gmOrder = append(msg.gmOrder, s)
			}
		}
	}
	return &msg
}
