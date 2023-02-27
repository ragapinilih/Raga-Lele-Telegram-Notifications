package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const (
	telegramAPIURL = "https://api.telegram.org/bot%s/sendMessage"

	// Every 07.00 and 18.00
	ROUTINE = "Kasih Makan Lele, cek PH dan cek TDS ya 🐟🐟🐟. Kalau PH kurang dari 6 berikan EM4 10 tutup botol per meter kubik atau dolomit 200gr per meter kubik"
	// Every Week on Friday at 18.01
	FLOK_CHECK = "Cek Flok di masing-masing kolam. Jangan lupa matikan airasi minimal 5 menit sebelum mengambil sample. Kalau flok masih kurang berikan campuran untuk molase, bakteri dan tepung sesuai dengan takaran"
	// Every 2 Week on Friday at 07.00
	HARVEST = "Ada yang sudah bisa dipanen? Kalau iya nanti malam makan terakhir, besoknya panen ya"
)

type telegramMessage struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

var err error
var message string
var message_type *string
var chatIDs []string

func sendTelegramNotification(botToken string, chatID int64, message string) error {
	// Create the Telegram message payload
	telegramMsg := telegramMessage{
		ChatID: chatID,
		Text:   message,
	}

	// Convert the message payload to a JSON string
	msgBytes, err := json.Marshal(telegramMsg)
	if err != nil {
		return err
	}

	// Create a new HTTP request to send the Telegram message
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(telegramAPIURL, botToken), bytes.NewReader(msgBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the HTTP request to the Telegram API
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the response status code indicates success
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get bot token and chat ID from environment variables
	botToken := os.Getenv("BOT_TOKEN")
	chatIDsEnv := os.Getenv("CHAT_ID")

	// Split the string on the comma delimiter
	chatIDs = strings.Split(chatIDsEnv, ",")

	message_type = flag.String("type", "", "message type is routine or flok_check")

	flag.Parse()

	if len(*message_type) == 0 {
		log.Fatal("Type is empty!")
	}

	switch *message_type {
	case "routine":
		message = ROUTINE
	case "flok_check":
		message = FLOK_CHECK
	case "harvest":
		message = HARVEST
	}

	for _, chatID := range chatIDs {
		// Parse chat ID from string to int64
		chatIDint, err := strconv.ParseInt(chatID, 10, 64)
		if err != nil {
			log.Fatal("Error parsing chat ID")
		}

		// Send a notification message to the specified chat ID
		err = sendTelegramNotification(botToken, chatIDint, message)
		if err != nil {
			fmt.Println(err)
		}
	}
}
