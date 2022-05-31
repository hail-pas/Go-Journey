package main

import (
	"bufio"
	"github.com/rs/zerolog/log"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:20000")

	if err != nil {
		log.Info().Msgf("create new conn failed, err: ", err)
		return
	}

	defer func(conn net.Conn) {
		err = conn.Close()
		if err != nil {
			log.Info().Msgf("close conn failed, err", err)
			return
		}
	}(conn)
	inputReader := bufio.NewReader(os.Stdin)

	for {
		input, _ := inputReader.ReadString('\n')
		inputInfo := strings.Trim(input, "\r\n")

		if strings.ToUpper(inputInfo) == "Q" {
			log.Info().Msgf("Bye-bye")
			return
		}

		_, err := conn.Write([]byte(inputInfo))

		if err != nil {
			return
		}

		buf := make([]byte, 512)

		n, err := conn.Read(buf)

		if err != nil {
			log.Info().Msgf("receive from server failed, err: ", err)
			return
		}
		log.Info().Msgf("received from server, data: %v", string(buf[:n]))

	}

}
