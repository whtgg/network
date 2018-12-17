/**
 * Created by Goland
 * User: lishangyuzi
 *
 * 实现一个server作中转群发消息，多个客户端聊天
 */
package main

import (
	"flag"
	"net"
	"fmt"
	"strings"
)

var port string
var ConnMap map[string]net.Conn = make(map[string]net.Conn)

func init() {
	flag.StringVar(&port, "port", "8080", "socket port")
	flag.Parse()
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", "127.0.0.1", port))
	if err != nil {
		fmt.Println("server start errot")
	}

	defer listener.Close()
	fmt.Println("server is waiting ...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("conn fail ...")
		}
		fmt.Println("connect client successed")

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	//分次读取  平衡读写
	for {
		data := make([]byte, 255)
		msg_read_num, err := conn.Read(data)
		if msg_read_num == 0 || err != nil {
			continue
		}

		//解析协议
		//client '\' 前表示区分交互类型
		msg_str_arr := strings.Split(string(data[0:msg_read_num]), "|")

		switch msg_str_arr[0] {
		case "nick":
			fmt.Println(conn.RemoteAddr(), "-->", msg_str_arr[1])
			//通知新用户加入
			for k, v := range ConnMap {
				if k != msg_str_arr[1] {
					v.Write([]byte("[" + msg_str_arr[1] + "]:join..."))
				}
			}
			ConnMap[msg_str_arr[1]] = conn
		case "say":
			for k, v := range ConnMap {
				if k != msg_str_arr[1] {
					fmt.Println("Send "+msg_str_arr[2]+" to ", k)
					v.Write([]byte("[" + msg_str_arr[1] + "]: " + msg_str_arr[2]))
				}
			}
		case "quit":
			for k, v := range ConnMap {
				if k != msg_str_arr[1] {
					v.Write([]byte("[" + msg_str_arr[1] + "]: quit"))
				}
			}
			delete(ConnMap, msg_str_arr[1])
		}
	}
}
