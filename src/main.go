package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// Wrap discordgo.Session to simplify functions call
type DiscordSession struct {
	*discordgo.Session
}

// Bot Discord's token
var botToken string

// Use to determine when user wants to be notified
type TimeOfNotification struct {
	Morning bool
	Midday  bool
	Evening bool
}

// Struct of user's data
type Dev struct {
	Name                     string
	UserId                   string
	DesiredNotificationTimes TimeOfNotification
}

// Slice of Devs for test TODO Remove
var devs = []Dev{
	{
		Name:   "UserName",
		UserId: "uuid",
		DesiredNotificationTimes: TimeOfNotification{
			Morning: true,
			Midday:  true,
			Evening: false,
		},
	},
}

// Create discord session using provided token and handle errors
func createDiscordSession(token string) *DiscordSession {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating discord session", err)
	}
	return &DiscordSession{session}
}

// Open websocket for session
func openDiscordConnection(session *DiscordSession) {
	err := session.Open()
	if err != nil {
		log.Fatal("Error opening connection", err)
	}
}

// Create a Channel to DM the user with provided userID
func (session *DiscordSession) createUserChannel(userId string) *discordgo.Channel {
	channel, err := session.UserChannelCreate(userId)
	if err != nil {
		log.Println("Error creating DM channel:", err)
		return nil
	}
	return channel
}

// For current session, send provided message to specified channel
func (session *DiscordSession) sendMessage(channel *discordgo.Channel, message string) {
	_, err := session.ChannelMessageSend(channel.ID, message)
	if err != nil {
		log.Println("Error sending message:", err)
	}
}

// TODO Send reminder to provided Devs accordingly to DesiredNotificationTimes
func (session *DiscordSession) sendReminder(devs []Dev) {
	for _, dev := range devs {
		channel := session.createUserChannel(dev.UserId)
		message := fmt.Sprintf("Hey %v ! Let's eat a pizza tonight !", dev.Name)
		session.sendMessage(channel, message)
	}
}

// Load environment variables from .env file in current directory
func loadDotEnv() {
	log.Println("Load .env file")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	log.Println("Loading successful")
}

func main() {
	log.Println("Feed your bot starting")

	// Load environment variables from .env file
	loadDotEnv()

	// Set token variable with discord app token
	botToken = os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("Empty Bot Token in env")
	}

	// Create new discord session with provided token
	discord := createDiscordSession(botToken)
	// Open websocket connection and ensure the connection will be disconnected
	openDiscordConnection(discord)
	defer discord.Close()

	log.Println("Feed your dev is running")
	// Send reminder to the slice of Devs passed as argument
	discord.sendReminder(devs)
}
