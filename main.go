package main

import (
	"fmt"
	"net"

	apiserver "github.com/o-my-god/observer/api/server"
	"github.com/o-my-god/observer/api/server/router/hello"
	"github.com/sirupsen/logrus"
)

func main() {

	l, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		logrus.Fatal("Failed to listen %s", "localhost:8080")
	}

	cfg := &apiserver.Config{}

	apiserver := apiserver.New(cfg)

	apiserver.Accept("localhost:8080", l)

	apiserver.InitRouter(hello.NewRouter())

	serveAPIWait := make(chan error)
	go apiserver.Wait(serveAPIWait)

	errAPI := <-serveAPIWait

	fmt.Printf("errAPI: %v", errAPI)

}
