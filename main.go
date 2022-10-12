package main

import (
	"github.com/gondalf/gondalf/resource/http"
	"github.com/gondalf/gondalf/server"
	"log"
)

func main() {
	s := server.New()
	s.MustLoadConfigFile("config.yml")
	if err := s.RegisterResourceFactory("http", http.FactoryFunc); err != nil {
		log.Fatalln(err)
	}
	if err := s.Run(); err != nil {
		log.Fatalln(err)
	}
}
