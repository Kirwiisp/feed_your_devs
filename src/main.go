package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

type Dev struct {
	name   string
	userId string
}

var token string

var devs = []Dev{
	{
		name:   "UserName",
		userId: "uuid",
	},
}

func sendReminder(session *discordgo.Session) {
	for _, dev := range devs {
		channel, err := session.UserChannelCreate(dev.userId)
		if err != nil {
			log.Println("Error creating DM channel:", err)
			continue
		}

		message := fmt.Sprintf("Hey %v ! Let's eat a pizza tonight !", dev.name)

		_, err = session.ChannelMessageSend(channel.ID, message)
		if err != nil {
			log.Println("Error sending DM:", err)
		}
	}
}

func main() {
	log.Println("Feed your bot starting")

	log.Println("Load .env file")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Loading successful")

	token = os.Getenv("BOT_TOKEN")
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating discord session", err)
	}

	err = discord.Open()
	if err != nil {
		log.Fatal("Error opening connection", err)
	}
	defer discord.Close()

	log.Println("Feed your dev is running")
	sendReminder(discord)
}
