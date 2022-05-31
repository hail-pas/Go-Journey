package main

import (
	"bufio"
	"github.com/rs/zerolog/log"
	"net"
)

func process(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Info().Msgf("conn close failed, err: ", err)
			return
		}
	}(conn)

	for {
		reader := bufio.NewReader(conn)
		var buf = make([]byte, 128, 128)

		n, err := reader.Read(buf)

		if err != nil {
			log.Info().Msgf("read from client failed, err:", err)
			break
		}
		receiveStr := string(buf[:n])
		log.Info().Msgf("receive from client, data:", receiveStr)
		_, err = conn.Write([]byte(receiveStr))
		if err != nil {
			log.Error().Msgf("send to client failed, err:", err)
			break
		}
	}

}

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:20000")
	if err != nil {
		log.Info().Msgf("listen failed, err:", err)
		return
	}
	log.Info().Msgf("listening on port: 20000 ...")
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Info().Msgf("accept failed, err: ", err)
			continue
		}
		go process(conn)
	}
}
