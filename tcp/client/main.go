package main

import (
	"Go-Journey/tcp/pkg"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:30000")
	if err != nil {
		fmt.Println("dial failed, err", err)
		return
	}
	defer conn.Close()
	for i := 0; i < 20; i++ {
		headerInfo := fmt.Sprintf("这是一个头信息-%d", i)
		body := pkg.Body{
			Version: fmt.Sprintf("v0.%02d", i),
			Topic:   "test",
			Author:  "phoenix",
			Content: "这是一个测试数据包体",
		}
		data, _ := pkg.GenerateRawData(headerInfo, body)
		conn.Write(data)
	}
}
