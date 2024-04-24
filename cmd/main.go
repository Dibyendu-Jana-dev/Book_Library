package main

import (
	"github.com/dibyendu/Authentication-Authorization/pkg/app"
	"github.com/joho/godotenv"
	"log"
)
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file: " + err.Error())
	}
	app.StartApp()
}
