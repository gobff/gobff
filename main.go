package main

import (
	"github.com/gondalf/gondalf/server"
	"log"
)

func main() {
	s := server.New()
	s.MustLoadConfigFile("config.yml")
	if err := s.Run(); err != nil {
		log.Fatalln(err)
	}
}
