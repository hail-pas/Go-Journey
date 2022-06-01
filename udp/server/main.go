package main

import (
	"github.com/rs/zerolog/log"
	"net"
)

func main() {
	listen, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 30000})

	if err != nil {
		log.Error().Msgf("listen failed, err: %v", err)
		return
	}

	defer func(listen *net.UDPConn) {
		err := listen.Close()
		if err != nil {
			log.Error().Msgf("close listen failed, err: %v", err)
			return
		}
	}(listen)

	for {
		var data [1024]byte
		n, addr, err := listen.ReadFromUDP(data[:])
		if err != nil {
			log.Error().Msgf("read udp failed, err: %v", err)
			continue
		}
		log.Info().Msgf("data:%v addr:%v count:%v\n", string(data[:n]), addr, n)

		_, err = listen.WriteToUDP(data[:n], addr)

		if err != nil {
			log.Error().Msgf("write to udp failed, err: ", err)
			continue
		}
	}

}
