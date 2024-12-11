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

	// Flok Application
	FLOK_APPLICATION = "\nMolase 250ml per meter kubik (didihkan terlebih dahulu)\nBakteri 25 gram per meter kubik\nTepung terigu 250 gram per meter kubik"

	// Flok is not enough
	FLOK_NOT_ENOUGH = "Jika flok kurang maka lakukan penambahan flok:" + FLOK_APPLICATION
	// Flok Overflow
	FLOK_OVERFLOW = "Jika flok berlebih maka lakukan:\n• Puasakan ikan sampai jendela makan berikutnya\n• Pada jendela makan berikutnya berikan pakan 80% dari biasanya\n• *JIKA NAFSU MAKAN IKAN BERKURANG* buang air dasar kolam kurang sebanyak 10cm dan tambahkan air baru sebanyak yang dibuang"

	// Flok Death
	FLOK_DEATH = "*JIKA BANYAK FLOK MATI SEGERA CEK FLOK!*"

	// Every 07.00 and 18.00
	ROUTINE = "Cek *PH*, *TDS*, *Suhu Kolam* dan *beri makan Ikan ya* 🐟🐟🐟\n\n*JIKA PH KURANG DARI 6 JANGAN DIBERI MAKAN!*\n\nYang harus dilakukan saat PH kurang dari 6 adalah:\n• Puasakan ikan sampai jendela makan berikutnya\n• Berikan dolomit 200gr per meter kubik yang dilarutkan terlebih dahulu dalam air sebelum ditebar ke kolam bioflok\n• Apabila hujan tidak berhenti berhari-hari, ikan bisa dipuasakan *MAKSIMAL* 3 hari. Lebih dari itu hubungi pemilik kolam"
	// Every Week on Friday at 18.01
	FLOK_CHECK = "Cek flok di setiap kolam. Jangan lupa matikan airasi minimal 5 menit sebelum mengambil sample.\n\n*JIKA FLOK MASIH KURANG BERIKAN CAMPURAN UNTUK MOLASE, BAKTERI DAN TEPUNG SESUAI DENGAN TAKARAN!*\n" + FLOK_APPLICATION

	// Every Monday and Thursday
	CLEAN_PRE_FILTER = "Jadwal hari ini bersihkan *Pre-Filter*"

	// Every 2 weeks on Wednesday
	CLEAN_MECHANIC_FILTER = "Jadwal hari ini bersihkan *Filter Mekanik*"

	// 2 Times a day every 08.00 and 19.00
	FEEDING = "Ayo berikan pakan ikan!"

	// Once a month at last day of month at 19.00
	HARVEST = "Ada kolam yang sudah bisa dipanen? Kalau iya nanti malam makan terakhir, besoknya panen ya"
)

type telegramMessage struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

var message string
var message_type *string
var chatIDs []string

func sendTelegramNotification(botToken string, chatID int64, message string) error {
	// Create the Telegram message payload
	telegramMsg := telegramMessage{
		ChatID:    chatID,
		Text:      message,
		ParseMode: "Markdown",
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

	message_type = flag.String("type", "", "message type is routine, harvest, flok_check or flok_death")

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
	case "flok_death":
		message = FLOK_DEATH
	case "flok_overflow":
		message = FLOK_OVERFLOW
	case "flok_not_enough":
		message = FLOK_NOT_ENOUGH
	case "clean_pre_filter":
		message = CLEAN_PRE_FILTER
	case "clean_mechanic_filter":
		message = CLEAN_MECHANIC_FILTER
	case "feeding":
		message = FEEDING
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
