package main

import (
	"spy-cats/internal/config"
)

func main() {
	config.MustLoad()

	log := SetupLogger()

	SetupAndRunServer()
	log.Info("server started")
}
