package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"

	"github.com/rs/zerolog/log"

	"wow/internal/pow"
	"wow/pkg/protocol"
)

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		log.Fatal().Msg("empty addr")
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatal().Err(err).Msg("resolve tcp")
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("dial tcp")
	}

	defer conn.Close()

	var task protocol.Task
	err = gob.NewDecoder(conn).Decode(&task)
	if err != nil {
		panic(err)
	}

	sol := pow.Calculate(task.Nonce, task.Difficulty)
	err = gob.NewEncoder(conn).Encode(protocol.Solution{
		Answer: sol,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("encode solution")
	}
	var msg protocol.Message
	err = gob.NewDecoder(conn).Decode(&msg)
	if err != nil {
		log.Fatal().Err(err).Msg("decode msg")
	}
	fmt.Println(msg.Quote)
}
