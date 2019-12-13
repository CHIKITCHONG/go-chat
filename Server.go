package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

/*
需求：
1.
单聊
群聊
2.
上线、下线通知
*/
var (
	clientsMap = make(map[string]net.Conn)
)

func SHandleError(err error, why string) {
	if err != nil {
		fmt.Println(why, err)
		os.Exit(1)
	}
}

func ioWithConn(conn net.Conn) {
	clientAddr := conn.RemoteAddr().String()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != io.EOF {
			SHandleError(err, "conn.Read")
		}
		if n > 0 {
			msg := string(buffer[:n])
			fmt.Println(clientAddr, msg)

			strs := strings.Split(msg, "#")
			if len(strs) > 1 {
				targetAddr := strs[0]
				targetMsg := strs[1]

				if targetAddr == "all" {
					// 群发消息
					for _, conn := range clientsMap {
						conn.Write([]byte(clientAddr + ":" + targetMsg))
					}
				} else {
					// 点对点消息
					for addr, conn := range clientsMap {
						if addr == targetAddr {
							conn.Write([]byte(clientAddr + ":" + targetMsg))
							break
						}
					}
				}
			} else {
				conn.Write([]byte("已阅：" + msg))
			}
		}
	}
}

func main() {
	// 建立服务端监听
	listener, e := net.Listen("tcp", "127.0.0.1:8888")
	SHandleError(e, "net.Listen")
	defer func() {
		for _, conn := range clientsMap {
			conn.Write([]byte("服务器进入维护"))
		}
		listener.Close()
	}()

	for {
		// 循环接入所有协程
		conn, e := listener.Accept()
		SHandleError(e, "listen.Accept")

		// 将每一个链接放入 map
		clientsMap[conn.RemoteAddr().String()] = conn

		// 在单独的协程中与每一个具体的协程聊天
		go ioWithConn(conn)
	}

	// 设置优雅退出逻辑

}
