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
	query := `SELECT id, name, game_chat_id,  current_task_id, time_update_task, total_players, status FROM games WHERE status = ?`
	rows, err := db.Query(query, models.StatusGamePlaying)
	if err != nil {
		log.Printf("Failed to get all active games: %v", err)
		
		return nil, err
	}
	defer rows.Close()

	var games []models.Game
	for rows.Next() {
		var game models.Game
		err := rows.Scan(&game.ID, &game.Name, &game.GameChatID, &game.CurrentTaskID, &game.TimeUpdateTask, &game.TotalPlayers, &game.Status)
		if err != nil {
			log.Printf("Error scanning game: %v", err)
			return nil, err
		}
		games = append(games, game)
	}

	return games, nil
}


func GetAllPlayersByGameID(gameId int) ([]models.Player, error) {
	query := `SELECT id, username, name, game_id, status, skipped, role FROM players WHERE game_id = ?`
	rows, err := db.Query(query, gameId)
	if err != nil {
		log.Println("Failed to get all players by game ID")
		return nil, err
	}
	defer rows.Close()

	var players []models.Player
	for rows.Next() {
		var player models.Player
		err := rows.Scan(&player.ID, &player.UserName, &player.Name, &player.GameID, &player.Status, &player.Skipped, &player.Role)
		if err != nil {
			log.Printf("Error scanning player: %v", err)
			return nil, err
		}
		players = append(players, player)
	}

	return players, nil
}

func HasPlayerResponded(playerID, gameID, taskID int) (bool, error) {
    query := `
        SELECT EXISTS(
            SELECT 1 
            FROM player_responses 
            WHERE player_id = ? 
                AND game_id = ? 
                AND task_id = ?
        )
    `
    
    var exists bool
    err := db.QueryRow(query, playerID, gameID, taskID).Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("failed to check player response: %w", err)
    }
    
    return exists, nil
}

func MarkNotificationUser(gameID, gameChatID, taskID, userID int) error {
    query := `
        INSERT INTO notifications (game_id, game_chat_id, task_id, user_id, notification_sent)
        VALUES (?, ?, ?, ?, 1)
    `
    
    _, err := db.Exec(query, gameID, gameChatID, taskID, userID)
    if err != nil {
        return fmt.Errorf("failed to create notification record: %w", err)
    }
	
	log.Printf("Marked notification sent for user ID %d, game ID %d, task ID %d", userID, gameID, taskID)
    
    return nil
}

func HasNotificationBeenSent(gameID, gameChatID, taskID, userID int) (bool, error) {
    query := `
        SELECT EXISTS(
            SELECT 1 
            FROM notifications 
            WHERE game_id = ? 
                AND game_chat_id = ?
                AND task_id = ?
                AND user_id = ?
                AND notification_sent = 1
        )
    `
    
    var exists bool
    err := db.QueryRow(query, gameID, gameChatID, taskID, userID).Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("failed to check notification: %w", err)
    }
    
    return exists, nil
}