// client
package main

import (
	"bufio" //便于读写的包buffer io
	"chat/codec"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "localhost:8080")
	conn, _ := net.DialTCP("tcp", nil, tcpAddr)
	defer conn.Close()
	fmt.Println("tcp connected...")
	go receiveMessage(conn)
	//控制台聊天功能
	for {
		var msg string
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			msg = scanner.Text()
		}
		if msg == "/quit" {
			break
		}
		b, err := codec.Encode(msg)
		if err != nil {
			fmt.Println("message encode error:", err)
		}
		conn.Write(b)
	}
}

//接收消息
func receiveMessage(conn *net.TCPConn) {
	reader := bufio.NewReader(conn)
	for {
		bytes, err := codec.Decode(reader)
		err, msg := ToInfo(bytes)
		if err != nil {
			fmt.Println("decode message error:", err)
			break
		}
		tm := time.Unix(msg.Time, 0)
		fmt.Printf("\n")
		fmt.Printf("-------------From:%s Time:%s -------------\n", msg.From, tm.Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("%s\n", msg.Text)
	}
}
