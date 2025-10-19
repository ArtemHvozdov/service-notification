package models

type Game struct {
	ID         int
	Name       string
	GameChatID int64
	CurrentTaskID int
	TimeUpdateTask int64 // Unix timestamp of the last task update
	TotalPlayers  int    // max 5
	Status        string // "waiting", "playing", "finished"
}

type Player struct {
	ID       int64
	UserName string
	Name     string
	Status   string
	Skipped  int
	GameID   int
	Role     string // "admin", "player"
	GameChatID int64
}

type RemindedPlayer struct {
	ID int
	Username string
	GameID   int
	TaskID   int
}

const (
	StatusGameWaiting  = "waiting"
	StatusGamePlaying  = "playing"
	StatusGameFinished = "finished"
)