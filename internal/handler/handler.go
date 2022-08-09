package handler

import (
	"encoding/gob"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"

	"wow/internal/difficulty"
	"wow/internal/pow"
	"wow/internal/quotes"
	"wow/pkg/protocol"
)

type Handler struct {
	diff *difficulty.Difficulty
}

func NewHandler(diff *difficulty.Difficulty) *Handler {
	rand.Seed(time.Now().Unix())
	return &Handler{diff: diff}
}

// Handle handles tcp requests
func (h *Handler) Handle(conn net.Conn) error {
	err := conn.SetDeadline(time.Now().Add(time.Minute))
	if err != nil {
		return fmt.Errorf("set deadline: %w", err)
	}
	h.diff.NewConn()
	defer h.diff.ConnOver()
	defer conn.Close()
	diff := h.diff.Difficulty()
	time.Sleep(time.Second)
	nonce := []byte(conn.RemoteAddr().String() + strconv.Itoa(int(time.Now().UnixNano())))
	err = gob.NewEncoder(conn).Encode(protocol.Task{
		Difficulty: diff,
		Nonce:      nonce,
	})
	if err != nil {
		return fmt.Errorf("encode init message: %w", err)
	}

	var resp protocol.Solution
	err = gob.NewDecoder(conn).Decode(&resp)
	if err != nil {
		return fmt.Errorf("decode message: %w", err)
	}
	var quote string
	if !pow.IsValid(nonce, resp.Answer, diff) {
		quote = "bad answer"
	} else {
		quote = quotes.Quotes[rand.Intn(len(quotes.Quotes))]
	}

	err = gob.NewEncoder(conn).Encode(protocol.Message{Quote: quote})
	if err != nil {
		return fmt.Errorf("encode qoute: %w", err)
	}
	return nil
}
