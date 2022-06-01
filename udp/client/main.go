package main

import (
	"github.com/rs/zerolog/log"
	"net"
)

func main() {
	socket, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 30000})

	if err != nil {
		log.Error().Msgf("connect to server failed, err: ", err)
		return
	}

	defer socket.Close()

	sendData := []byte("Hello server")

	_, err = socket.Write(sendData)

	if err != nil {
		log.Error().Msgf("send data failed, err: %v", err)
		return
	}

	data := make([]byte, 4096)
	n, remoteAddr, err := socket.ReadFromUDP(data)
	if err != nil {
		log.Error().Msgf("receive data failed， err: ", err)
		return
	}
	log.Info().Msgf("recv:%v addr:%v count:%v\n", string(data[:n]), remoteAddr, n)
}
