package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

var (
	chanQuit = make(chan bool, 0)
)

// 异常退出时，得不到 conn
func CHandleError(err error, why string) {

	if err != nil {
		fmt.Println(why, err)
		os.Exit(1)
	}
}

func main() {

	// 处理异常退出
	defer func() {
		err := recover()
		if err != nil {

		}
	}()

	// 拨号连接
	conn, e := net.Dial("tcp", "127.0.0.1:8888")
	CHandleError(e, "net.Dial")
	defer func() {
		// 通知服务端客户端下线了
		conn.Write([]byte("exit"))
		conn.Close()
	}()

	// 在一条独立的协程中接收输入，并发送消息
	go handleSend(conn)

	// 在一条独立的协程中接收服务端消息
	go handleReceive(conn)

	// 设置优雅退出
	<-chanQuit
}
// zhushi
func handleReceive(conn net.Conn) {
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != io.EOF {
			CHandleError(err, "conn.Read")
		}

		if n > 0 {
			msg := string(buffer[:n])
			fmt.Println(msg)
		}
	}
}

func handleSend(conn net.Conn) {

	reader := bufio.NewReader(os.Stdin)

	for {
		// 读取标准输入
		lineBytes, _, _ := reader.ReadLine()

		// 发送到服务端
		_, err := conn.Write(lineBytes)
		CHandleError(err, "conn.Write")

		// 因为 ctrl + c 强制退出,可能 main 中的 defer 不会执行
		if string(lineBytes) == "exit" {
			// 正常退出
			os.Exit(1)
		}
	}

}
