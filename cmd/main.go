package main

import (
	"github.com/askrella/go-mockserver/internal/server"
	"log"
)

func main() {
	log.Println("Starting Go-MockServer v1.0.0 by Askrella")
	server.InitializeServer()
}
