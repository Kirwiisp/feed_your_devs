package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var token string

func main() {
	log.Println("Feed your bot starting")

	log.Println("Load .env file")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Loading successful")

	token = os.Getenv("BOT_TOKEN")
	discord, _ := discordgo.New("Bot " + token)
	test := struct {
		test  *discordgo.Session
		title string
	}{
		discord,
		"test",
	}
	fmt.Println(token[:10])
	fmt.Println(test.title)
}
