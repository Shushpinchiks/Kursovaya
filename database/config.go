package database

import (
	"fmt"
	"github.com/joho/godotenv"
)

func InitEnv() {
	err := godotenv.Load("config.env")
	if err != nil {
		fmt.Println("Error loading config.env file")
	}
}
