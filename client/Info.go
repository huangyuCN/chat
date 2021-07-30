package main

import (
	"encoding/json"
)

type Info struct {
	From string //发送者
	Time int64  //发送时间
	Text string //消息内容
}

// ToInfo 将json格式转化成消息对象
func ToInfo(data []byte) (error, *Info) {
	var msg Info
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return err, nil
	}
	return nil, &msg
}
