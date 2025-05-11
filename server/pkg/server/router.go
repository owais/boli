package server

import "fmt"

type commandHandler func(cmd *command) error

type Router struct {
	handlers map[string]commandHandler
}

func NewRouter(handlers map[string]commandHandler) *Router {
	return &Router{
		handlers: handlers,
	}
}

func (r *Router) Handle(cmd *command) error {
	handler, ok := r.handlers[cmd.Name]
	if !ok {
		return fmt.Errorf("no handler for command: %s", cmd.Name)
	}
	return handler(cmd)
}
