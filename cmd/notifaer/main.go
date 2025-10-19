package main

import (
	// "database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ArtemHvozdov/service-notification/internal/database"
	"github.com/ArtemHvozdov/service-notification/internal/models"
	"github.com/ArtemHvozdov/service-notification/internal/utils"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/telebot.v3"
)

var hurryUppMsg = []string {
  "@%s агов, красуне! Просто нагадую, що твоя відповідь у грі - це як сонце для квіточки дружби 🌸 Ти ще встигаєш відповісти на своє завдання! 😘",
  "@%s хей, зіронько! У тебе залишилось 12 годин, щоб відповісти! Я знаю, тобі є що сказати! 💛",
  "@%s на Землі зафіксовано критичний рівень нестачі твоїх відповідей у грі! 📢 Зайди в гру та врятуй ситуацію!",
  "@%s ти ж не хочеш, щоб я текстила тобі «Ну що, ти вже відповіла?» 😏 Ось і я так думаю… Я ж знаю, яка ти заклопотана, крихітко! Але навіть коротка відповідь - це велика приємність для твоїх подружок 💛",
  "@%s ого, ти щось задумала… 🤔 Уже 12 годин без відповіді у грі. Певне, ця відповідь буде просто неймовірною, як і ти сама! Гарного дня і стеж за часом! 🌟",
  "@%s гей, кицю, я знаю, що ти зайнята, але просто хочу нагадати про нашу гру 💛 Твоя відповідь важлива для твоєї бесті! 🐈",
  "@%s 🔮 я запитала у чарівної кулі, чи ти скоро відповіси у грі. Вона видала: \"Сумнівно\" 🤨 Давай змінимо її пророцтво, га? 😆",
}


func SendTelegramMsgToChat(bot *telebot.Bot, chatID int64, userName string) error {
	chat := &telebot.Chat{ID: chatID}
	txt := utils.GetRandomMsg(hurryUppMsg)
	message := fmt.Sprintf(txt, userName)
	_, err := bot.Send(chat, message)
	return err
}

func main() {
	log.Println("=== Notifier started ===")
	log.Printf("Execution time: %s", time.Now().Format("2006-01-02 15:04:05"))

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Failed to load .env file, using environment variables")
	}

	// Get environment variables
	dbDir := os.Getenv("DATABASE_DIR")
	if dbDir == "" {
		log.Fatal("DATABASE_DIR environment variable is required")
	}

	dbFile := os.Getenv("DATABASE_FILE")
	if dbFile == "" {
		log.Fatal("DATABASE_FILE environment variable is required")
	}

	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	if telegramToken == "" {
		log.Println("Warning: TELEGRAM_TOKEN not set")
	}

	// Initialization Telegram bot
	pref := telebot.Settings{
		Token:  telegramToken,
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatalf("Failed to create Telegram bot: %v", err)
	}

	// dbPath := dbDir + dbFile
	// log.Printf("Database path: %s", dbPath)

	// // Create database directory if it doesn't exist
	// if err := os.MkdirAll(dbDir, 0755); err != nil {
	// 	log.Fatalf("Error creating database directory: %v", err)
	// }

	dataDir := "/app/data"
	dataFile := "tg-game-bot.db"
	dbPath := fmt.Sprintf("%s/%s", dataDir, dataFile)

	db, err := database.NewDatabase(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB(db)

	activeGames, err := database.GetAllActiveGames()
	if err != nil {
		log.Fatalf("Failed to get active games: %v", err)
	}

	timeNow := time.Now().Unix()

	if len(activeGames) == 0 {
		log.Println("No active games found.")
	} else {
		log.Printf("Found %d active games:", len(activeGames))
	}

	remindersMap := make(map[int64][]*models.RemindedPlayer)

	for _, game := range activeGames {
		log.Printf("Active Game: ID=%d, Name=%s, IdChat=%d ,CurrentTaskID=%d",
			game.ID, game.Name, game.GameChatID, game.CurrentTaskID)

		log.Printf("Time now: %d, TimeUpdateTask: %d", timeNow, game.TimeUpdateTask)
		log.Printf("Delta time: %d", timeNow - game.TimeUpdateTask)

		// timeNow - game.TimeUpdateTask >= 600 && timeNow - game.TimeUpdateTask < 900
		if timeNow - game.TimeUpdateTask >= 600 && timeNow - game.TimeUpdateTask < 1200 {
			log.Print("Start checking players...")
            //continue
			players, err := database.GetAllPlayersByGameID(game.ID)
			if err != nil {
				log.Printf("Failed to get players for game ID %d: %v", game.ID, err)
			}

			for _, p := range players {
				hasResp, err := database.HasPlayerResponded(int(p.ID), game.ID, game.CurrentTaskID)
				if err != nil {
					log.Printf("Error checking response for player ID %d: %v", p.ID, err)
				}
				
				hasNotification, err := database.HasNotificationBeenSent(game.ID, int(game.GameChatID), game.CurrentTaskID, int(p.ID))
				if err != nil {
					log.Printf("Error checking notification for player ID %d: %v", p.ID, err)
				}
				
				if !hasResp && !hasNotification {
					log.Printf("Game chat id: %d", game.GameChatID)
					remindersMap[game.GameChatID] = append(remindersMap[game.GameChatID], &models.RemindedPlayer{
						ID:       int(p.ID),
						Username: p.UserName,
						GameID:   game.ID,
						TaskID:   game.CurrentTaskID,
					})
				}
			}
		}

    }
	
	// through the entire remindersMap
	for gameChatID, players := range remindersMap {
		log.Printf("Game Chat ID: %d, Players needing reminder: %d", gameChatID, len(players))
		
		for _, player := range players {
			log.Printf(" Chat: %d - Player ID: %d, Username: @%s", gameChatID, player.ID, player.Username)
			SendTelegramMsgToChat(bot, gameChatID, player.Username)
			database.MarkNotificationUser(player.GameID, int(gameChatID), player.TaskID, player.ID)
		}
	}

}
