package main

import (
	"github.com/gondalf/gondalf/server"
	"github.com/gondalf/gondalf/source/http"
	"log"
)

func registerSourceFactories(s server.Server) {
	if err := s.RegisterSourceFactory("http", http.FactoryFunc); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	s := server.New()
	s.MustLoadConfigFile("config.yml")
	registerSourceFactories(s)
	if err := s.Run(); err != nil {
		log.Fatalln(err)
	}
}
