package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"wow/internal/difficulty"
	"wow/internal/handler"

	"github.com/rs/zerolog/log"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal().Msg("empty port")
	}
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal().Err(err).Msg("listen port")
	}
	defer ln.Close()

	h := handler.NewHandler(difficulty.New(100))
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		log.Info().Msgf("got signal %s", (<-sigCh).String())
		cancel()
	}()
	wg := sync.WaitGroup{}
	for {
		select {
		case <-ctx.Done():
			log.Error().Err(ctx.Err()).Msg("context done")
			return
		default:
			conn, err := ln.Accept()
			if err != nil {
				log.Error().Err(err).Msg("accept")
				continue
			}

			go func() {
				wg.Add(1)
				e := h.Handle(conn)
				wg.Done()
				if e != nil {
					log.Error().Err(e).Msg("handle")
				}
			}()
		}
	}
	wg.Wait()
}
