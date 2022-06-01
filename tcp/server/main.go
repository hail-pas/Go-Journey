package main

import (
	"Go-Journey/tcp/pkg"
	"bufio"
	"bytes"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net"
	"os"
)

var rawFlow = make(chan []byte, 4096)
var cleanedPackage = make(chan *pkg.TcpPackage, 4096)
var signals = make(chan os.Signal, 1)
var pkgBuffer []byte

func process(conn net.Conn) {
	fmt.Println("process")
	defer func() {
		conn.Close()
	}()

	var buffers bytes.Buffer
	var buf = make([]byte, 4096, 4096)
	reader := bufio.NewReader(conn)
	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("read from client failed, err:", err)
			break
		}
		buffers.Write(buf[:n])
		rawFlow <- buffers.Bytes()
		buffers.Reset()
	}
}

func main() {

	listen, err := net.Listen("tcp", "127.0.0.1:30000")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	defer func() {
		listen.Close()
	}()
	log.Info().Msgf("listening ")
	go printPkg()
	go splitMessage()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		go process(conn)
	}
}

func printPkg() {
	for {
		select {
		case tcpPackage := <-cleanedPackage:
			log.Info().Msgf("%v", tcpPackage.Body)
		}
	}
}

func splitMessage() {
	log.Info().Msgf("splitMessage")
	var data []byte
	for {
		select {
		case segment := <-rawFlow:
			pkgBuffer = append(pkgBuffer, segment...)
			for {
				if len(pkgBuffer) < 8 {
					break
				}
				pkgLength := pkg.ConvertByte2Int(pkgBuffer[:8])
				bufferLength := uint64(len(pkgBuffer))
				if bufferLength >= pkgLength {
					data = pkgBuffer[:pkgLength]
					pkgBuffer = pkgBuffer[pkgLength:]
					cleanedPackage <- pkg.NewTcpPackage(data)
				} else {
					break
				}
			}
		}
	}
}
