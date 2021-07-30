// chat project main.go
package main

import (
	"fmt"
	"net"
)

//用来记录所有的客户端连接
//var ConnMap map[string]*net.TCPConn

var UserManager *userManager
var RoomManager *roomManager
var HistoryMessageLen uint = 50
var SensitiveWords []string
var PopularManager *popular

func main() {
	SensitiveWords = LoadSensitiveWords()
	UserManager = NewUserManager()
	RoomManager = NewRoomManager(HistoryMessageLen)
	PopularManager = NewPopular()

	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "localhost:8080")
	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)
	defer tcpListener.Close()
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			continue
		}

		fmt.Println("有个一个连接接入：" + tcpConn.RemoteAddr().String())
		//开启新的线程来接收消息和发送消息
		go NewConnManager(tcpConn)
	}
}

//func main() {
//	var tcpAddr *net.TCPAddr
//	ConnMap = make(map[string]*net.TCPConn)
//	tcpAddr, _ = net.ResolveTCPAddr("tcp", "localhost:8080")
//	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)
//	defer tcpListener.Close()
//
//	for {
//		tcpConn, err := tcpListener.AcceptTCP()
//		if err != nil {
//			continue
//		}
//
//		fmt.Println("有个一个连接接入：" + tcpConn.RemoteAddr().String())
//		//新连接加入map
//		ConnMap[tcpConn.RemoteAddr().String()] = tcpConn
//		//开启新的线程来接收消息和发送消息
//		go tcpPipe(tcpConn)
//	}
//}
//
//func tcpPipe(conn *net.TCPConn) {
//	ipStr := conn.RemoteAddr().String()
//	defer func() {
//		fmt.Println("断开连接：" + ipStr)
//		//移除连接
//		delete(ConnMap, ipStr)
//		conn.Close()
//	}()
//	reader := bufio.NewReader(conn)
//	for {
//		//		Message, err := reader.ReadString('\n')
//		Message, err := codec.Decode(reader)
//		if err != nil {
//			return
//		}
//		fmt.Println(conn.RemoteAddr().String() + ":" + string(Message))
//		//广播信息
//		broadcastMesage("转发消息：" + string(Message))
//	}
//}
//
//func broadcastMesage(Message string) {
//	fmt.Println("当前连接个数：", len(ConnMap))
//	//	b := []byte(Message)
//	b, err := codec.Encode(Message)
//	if err != nil {
//		fmt.Println("编码错误", err.Error())
//	}
//	//遍历所有连接
//	for _, conn := range ConnMap {
//		conn.Write(b)
//	}
//}
