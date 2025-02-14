package main

import (
	"bearer-issuer/api"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Version 2.0.0")
	godotenv.Load(".env")

	go api.StartAPI()

	select {}
}
