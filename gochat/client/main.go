package main

import (
	"net"
	"flag"
	"fmt"
)

var nick string
var port string

func init() {
	flag.StringVar(&port, "port", "8080", "socket port")
	flag.Parse()
}

func main() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", "127.0.0.1", port))
	if err != nil {
		fmt.Println("conn fail...")
	}
	defer conn.Close()
	fmt.Println("connect server successed \n")

	//给自己取一个昵称吧
	fmt.Println("Make a nickname:")
	fmt.Scanf("%s", &nick)
	fmt.Println("hello : ", nick)
	//加入聊天室
	conn.Write([]byte("nick|" + nick))

	//接受其他人的信息
	go Handle(conn)

	//在聊天室聊天
	var msg string
	for {
		msg = ""
		fmt.Scan(&msg)
		conn.Write([]byte("say|" + nick + "|" + msg))
		if msg == "quit" {
			conn.Write([]byte("quit|" + nick))
			break
		}
	}
}

func Handle(conn net.Conn) {
	for {
		data := make([]byte, 255)
		msg_read_num, err := conn.Read(data)
		if msg_read_num == 0 || err != nil {
			break
		}
		fmt.Println(string(data[0:msg_read_num]))
	}
}
