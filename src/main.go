package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// Wrap discordgo.Session to simplify functions call
type DiscordSession struct {
	*discordgo.Session
}

type UserId string

type NotificationTime string

// Struct of user's data
type Dev struct {
	Name                     string
	UserId                   UserId
	DesiredNotificationTimes []NotificationTime
}

// Reference time layout "Mon Jan 2 15:04:05 MST 2006"
// var h,_ = time.Parse(timeLayout,test)
const timeLayout = "15:04"

// Bot Discord's token
var botToken string

// ticker set to 1 minute
var tickerMinute *time.Ticker = time.NewTicker(1 * time.Minute)

// Slice of Devs for test TODO Remove
var devs = []Dev{
	{
		Name:                     "UserName",
		UserId:                   "UserID",
		DesiredNotificationTimes: []NotificationTime{"08:15", "22:30"},
	},
}

// Test MAP TODO Remove
var testMapTime = map[NotificationTime][]UserId{
	"08:15": {devs[0].UserId},
	"22:30": {devs[0].UserId},
}

// Return the slice containing the id of user to notify for the provided time TODO change call to DB
func getIdToNotify(time NotificationTime) []UserId {
	return testMapTime[time]
}

// Ensure input string is formated as "15:04" so it can be used as NotificationTime
func validateNotificationTime(input string) (NotificationTime, error) {
	isFormated := regexp.MustCompile(`^\d{2}:\d{2}$`).MatchString(input)
	_, err := time.Parse(timeLayout, input)
	if err != nil {
		log.Println("Invalid time: ", err)
	}
	if isFormated {
		return NotificationTime(input), nil
	}
	return "", fmt.Errorf("Incorrect formating: %s ", input)
}

// Validate the discord User ID (length between 17 and 18, only digits)
func validateUserId(input string) (UserId, error) {
	validLength := false
	if len(input) == 17 || len(input) == 18 {
		validLength = true
	}
	onlyDigits := regexp.MustCompile(`^\d+$`).MatchString(input)
	if validLength && onlyDigits {
		userId := UserId(input)
		return userId, nil
	}
	return "", fmt.Errorf("invalid UserId: %s ", input)
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

// TODO personalize to message
func (session *DiscordSession) notify(id UserId) {
	channel := session.createUserChannel(string(id))
	message := fmt.Sprintf("Hey ! Let's eat a pizza tonight !")
	session.sendMessage(channel, message)
}

// Get list UserId where desired notification time match current time and notify them
func (session *DiscordSession) notifyUsersNow() {
	now := NotificationTime(time.Now().Format(timeLayout))
	idToNotify := getIdToNotify(now)
	for i := range idToNotify {
		session.notify(idToNotify[i])
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

	// Handle killing of the process
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Load environment variables from .env file
	loadDotEnv()

	// Set token variable with discord app token
	botToken = os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("Empty Bot Token in env")
	}

	// Create new discordSession session with provided token
	discordSession := createDiscordSession(botToken)
	// Open websocket connection and ensure the connection will be disconnected
	openDiscordConnection(discordSession)
	defer func() {
		discordSession.Close()
		log.Println("Discord session closed")
	}()

	log.Println("Feed your dev is running")
	// Send reminder to the slice of Devs passed as argument
	go func() {
		for {
			select {
			case <-tickerMinute.C:
				discordSession.notifyUsersNow()
			}
		}
	}()
	<-quit // wait until signal is received
	log.Println("Shutting down")
}
