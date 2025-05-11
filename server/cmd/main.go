package main

import "github.com/owais/boli/server/pkg/server"

func main() {
	srv := &server.Server{}
	err := srv.Start()
	if err != nil {
		panic(err)
	}
}
