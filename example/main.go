package main

import (
	"net/http"

	"github.com/nat-brown/go-mimic/example/app"
	"github.com/nat-brown/go-mimic/example/client"
	"github.com/nat-brown/go-mimic/example/dependencies"
	"github.com/nat-brown/go-mimic/example/logger"
)

func main() {
	go http.ListenAndServe(":3000", dependencies.New())
	http.ListenAndServe(":8080", app.New(client.New(), logger.New()))
}
