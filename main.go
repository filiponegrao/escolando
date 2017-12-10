package main

import (
	"log"
	"os"
	"strconv"

	"github.com/filiponegrao/escolando/db"
	"github.com/filiponegrao/escolando/server"
)

// main ...
func main() {
	database := db.Connect()
	s := server.Setup(database)
	port := "8080"

	if p := os.Getenv("PORT"); p != "" {
		if _, err := strconv.Atoi(p); err == nil {
			port = p
		}
	}

	// Inicia configuracoes basicas
	if err := server.InitConfigurations(*database); err != nil {
		log.Println(err)
		return
	}

	s.Run(":" + port)

}
