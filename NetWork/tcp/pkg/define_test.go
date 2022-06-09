package pkg

import (
	"fmt"
	"testing"
)

func TestNewTcpPackage(t *testing.T) {
	headerInfo := "这是一个头信息"
	body := Body{
		Version: "v0.01",
		Topic:   "test",
		Author:  "phoenix",
		Content: "这是一个测试数据包体",
	}

	data, err := GenerateRawData(headerInfo, body)
	if err != nil {
		fmt.Println("failed", err)
		return
	}
	fmt.Println(data)
	tcpPackage := NewTcpPackage(data)
	fmt.Println(tcpPackage)
	fmt.Println(tcpPackage.Body)
}
