package main

import (
	"log"
	"main/cmd/api/docs"
	"main/pkg/config"
	di "main/pkg/di"
)

func main() {
	docs.SwaggerInfo.Title = "GameVerse"

	docs.SwaggerInfo.Version = "1.0"

	docs.SwaggerInfo.Host = "localhost:1245"
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	config, configErr := config.LoadConfig()
	if configErr != nil {
		log.Fatal("cannot load config: ", configErr)
	}

	server, diErr := di.InitializeAPI(config)
	if diErr != nil {
		log.Fatal("cannot start server: ", diErr)
	} else {
		server.Start()
	}
}
