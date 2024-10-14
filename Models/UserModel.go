package Models

import "time"

type Users struct {
	UserId   string `json:"userid" db:"userid"`
	UserName string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
	// use camelcase and snake case carefully
	CreatedAt  time.Time `json:"created_at" db:"createdAt"`
	ArchivedAt time.Time `json:"archived_at" db:"archived_at"`
	Todos      Todos
}

type UserLogin struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type UserCtx struct {
	UserID    string `json:"userId" db:"user_id"`
	Name      string `json:"name" db:"name"`
	Email     string `json:"email" db:"email"`
	SessionID string `json:"sessionId" db:"session_id"`
}
