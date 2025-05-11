package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/lesismal/nbio/nbhttp"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"go.uber.org/zap"
)

type Server struct {
	router *Router
}

func (s *Server) Start() error {
	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	defer logger.Sync()

	upgrader := s.newUpgrader(s)
	mux := &http.ServeMux{}
	mux.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		s.onWebsocketConnect(upgrader, w, r)
	})

	s.router = NewRouter(map[string]commandHandler{
		"game:new": newGameHandler,
	})

	engine := nbhttp.NewEngine(nbhttp.Config{
		Network:                 "tcp",
		Addrs:                   []string{"localhost:8080"},
		MaxLoad:                 1000000,
		ReleaseWebsocketPayload: true,
		Handler:                 mux,
	})

	if err = engine.Start(); err != nil {
		fmt.Printf("nbio.Start failed: %v\n", err)
		return err
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	return engine.Shutdown(ctx)
}

func (s *Server) onMessage(messageType websocket.MessageType, data []byte) error {
	cmd := &command{}
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("error decoding command: %w", err)
	}

	s.router.Handle(cmd)
	return nil
}

func (s *Server) onWebsocketConnect(upgrader *websocket.Upgrader, w http.ResponseWriter, r *http.Request) {
	// TODO:authenticate request
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Upgraded:", conn.RemoteAddr().String())
}

func (s *Server) newUpgrader(srv *Server) *websocket.Upgrader {
	u := websocket.NewUpgrader()

	u.OnOpen(func(c *websocket.Conn) {
		// echo
		fmt.Println("OnOpen:", c.RemoteAddr().String())
	})

	u.OnMessage(func(c *websocket.Conn, messageType websocket.MessageType, data []byte) {
		// echo
		// global message router
		// c.WriteMessage(messageType, data)
		err := srv.onMessage(messageType, data)
		if err != nil {
			fmt.Println("Error:", err)
		}
	})

	u.OnClose(func(c *websocket.Conn, err error) {
		fmt.Println("OnClose:", c.RemoteAddr().String(), err)
	})

	return u
}
