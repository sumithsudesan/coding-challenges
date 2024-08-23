package main

import (
	"log"

	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/config"
	"github.com/fiskaly/coding-challenges/signing-service-challenge-go/persistence"
)

//How to use?
//1.  export STORAGE_TYPE="in-memory"
//    export PORT=":8080"
//2.  go run main.go

// Main function - starting point
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("[ERROR] Failed to load configuration: ", err)
	}

	// Create storage based on the type
	store := persistence.NewStorage(cfg.StorageType)

	// Create sevrver- port and storage
	svr := api.NewServer(cfg.Port, store)
	log.Println("Server starting .....")
	// Running server
	if err := svr.Run(); err != nil {
		log.Fatal("[ERROR] Could not start server on port: ", cfg.Port)
	}
}
