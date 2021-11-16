package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/reaper47/recipya/cmd"
)

func main() {
	cmd.Execute()
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
}
