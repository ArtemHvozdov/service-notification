package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/ArtemHvozdov/service-notification/internal/models"
)

var db *sql.DB

func NewDatabase(dbPath string) (*sql.DB, error) {
	var err error
	db, err = sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_foreign_keys=on&_busy_timeout=60000")
	if err != nil {
		log.Fatalf("Error connection database: %v", err)
		return nil, err
	}

	db.SetMaxOpenConns(5) 
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 5)

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	return db, nil	
}

func CloseDB(Db *sql.DB) {
	if Db != nil {
		if err := Db.Close(); err != nil {
			er1 := fmt.Errorf("error closing database connection: %v", err)
			log.Println(er1)
		} else {
			log.Println("The database connection was closed successfully.")
		}
	}
}

func GetAllActiveGames() ([]models.Game, error) {
	query := `SELECT id, name, current_task_id, total_players, status FROM games WHERE status = ?`
	rows, err := db.Query(query, models.StatusGamePlaying)
	if err != nil {
		log.Printf("Failed to get all active games: %v", err)
		
		return nil, err
	}
	defer rows.Close()

	var games []models.Game
	for rows.Next() {
		var game models.Game
		err := rows.Scan(&game.ID, &game.Name, &game.CurrentTaskID, &game.TotalPlayers, &game.Status)
		if err != nil {
			log.Printf("Error scanning game: %v", err)
			return nil, err
		}
		games = append(games, game)
	}

	return games, nil
}
