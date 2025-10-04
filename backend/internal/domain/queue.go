package domain

// import "time"

type Game struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Max_slots        int    `json:"max_slots"`
	Current_people   int    `json:"current_people"`
	Duration_seconds int    `json:"duration_seconds"`
}

type ListGames []Game

type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginInfo struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RoleInfo struct {
	Role string `json:"role"`
}

type ChangeInfo struct {
	UserID int `json:"user_id"`
	GameID int `json:"game_id"`
}

type PosInfo struct {
	Pos int `json:"position"`
}

type IdInfo struct{
	Id int `json:"id"`
}