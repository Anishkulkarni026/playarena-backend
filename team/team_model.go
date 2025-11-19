// team/team_model.go
package team

import (
	"time"
	"database/sql"
)


type Team struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	OwnerID   int64     `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
}

type TeamMember struct {
	ID       int64     `json:"id"`
	TeamID   int64     `json:"team_id"`
	UserID   int64     `json:"user_id"`
	Status   string    `json:"status"`
	JoinedAt time.Time `json:"joined_at"`
}

// Struct for the "Create Team" request
type CreateTeamRequest struct {
	Name string `json:"name" binding:"required"`
}

// UserTeamDetails holds info about a team a user is in
type UserTeamDetails struct {
	TeamID     int64     `json:"team_id"`
	TeamName   string    `json:"team_name"`
	OwnerID    int64     `json:"owner_id"`
	UserStatus string    `json:"user_status"` // "pending" or "joined"
	JoinedAt   sql.NullTime `json:"joined_at"`
}
// TeamMemberDetails holds public info about a team member
type TeamMemberDetails struct {
	UserID    int64  `json:"user_id"`
	TeamID    int64  `json:"team_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Status    string `json:"status"`
}



type ChatMessage struct {
	ID             int64     `json:"id"`
	TeamID         int64     `json:"team_id"`
	UserID         int64     `json:"user_id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	MessageContent string    `json:"message_content"`
	CreatedAt      time.Time `json:"created_at"`
}

// PostMessageRequest is the struct for the API request body
type PostMessageRequest struct {
	Message string `json:"message" binding:"required"`
}

// ChatMessageDetails includes sender information
type ChatMessageDetails struct {
	ID             int64     `json:"id"`
	TeamID         int64     `json:"team_id"`
	UserID         int64     `json:"user_id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	MessageContent string    `json:"message_content"`
	CreatedAt      time.Time `json:"created_at"`
}