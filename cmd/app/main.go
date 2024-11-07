package main

import (
	"githib.com/s4bb4t/leadgen/internal/config"
	"githib.com/s4bb4t/leadgen/internal/logger/sl"
	"githib.com/s4bb4t/leadgen/internal/server"
	"githib.com/s4bb4t/leadgen/internal/storage"
)

func main() {
	cfg := config.MustLoad()
	repo := storage.MustLoad(cfg)
	log := sl.InitLogger()

	server.Run(log, repo)
}
