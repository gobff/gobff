package main

import (
	"github.com/carlosrodriguesf/gobff/internal/server"
)

func main() {
	s := server.New()
	s.LoadConfigFile()
	s.Run()
}
