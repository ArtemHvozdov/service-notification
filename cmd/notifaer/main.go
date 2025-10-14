package main

import (
	// "database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ArtemHvozdov/service-notification/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	//"gopkg.in/telebot.v3"
)

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
	// pref := telebot.Settings{
	// 	Token:  telegramToken,
	// }

	// bot, err := telebot.NewBot(pref)
	// if err != nil {
	// 	log.Fatalf("Failed to create Telegram bot: %v", err)
	// }

	// dbPath := dbDir + dbFile
	// log.Printf("Database path: %s", dbPath)

	// // Create database directory if it doesn't exist
	// if err := os.MkdirAll(dbDir, 0755); err != nil {
	// 	log.Fatalf("Error creating database directory: %v", err)
	// }

	dataDir := "/app/data"
dataFile := "tg-game-bot.db"
dbPath := fmt.Sprintf("%s/%s", dataDir, dataFile)

	// if err := os.MkdirAll("dataDir", 0755); err != nil {
	// 	log.Println("Error creating folder:")
	// 	//log.Fatalf("Error creating folder %s: %v", dataDir, err)
	// }

	db, err := database.NewDatabase(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB(db)

	activeGames, err := database.GetAllActiveGames()
	if err != nil {
		log.Fatalf("Failed to get active games: %v", err)
	}

	if len(activeGames) == 0 {
		log.Println("No active games found.")
	} else {
		log.Printf("Found %d active games:", len(activeGames))
	}

	for _, game := range activeGames {
		log.Printf("Active Game: ID=%d, Name=%s, CurrentTaskID=%d, TotalPlayers=%d, Status=%s",
			game.ID, game.Name, game.CurrentTaskID, game.TotalPlayers, game.Status)
	}

	// Open database connection	
	// db, err := sql.Open("sqlite3", dbPath+
	// 	"?_journal_mode=WAL"+
	// 	"&_foreign_keys=on"+
	// 	"&_busy_timeout=5000") 
	// if err != nil {
	// 	log.Fatalf("Failed to open database: %v", err)
	// }

	// Setting pool of connection
	// db.SetMaxOpenConns(5) 
	// db.SetMaxIdleConns(5)
	// db.SetConnMaxLifetime(time.Minute * 5)

	// defer db.Close()

	// Test database connection
	// if err := db.Ping(); err != nil {
	// 	log.Fatalf("Failed to ping database: %v", err)
	// }
}



// 1. Получить все активные игры из БД (таблица games, есть айди чата в игре) +
// 2. Если прошло больше N минут с обрновления таски в игре, проверить кто еще не ответил
//    2.1 