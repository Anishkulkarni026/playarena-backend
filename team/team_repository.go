// team/team_repository.go
package team

import (
	"log"
	"github.com/JkD004/playarena-backend/db"
	"database/sql"
	"errors"
)

// CreateTeam inserts a new team into the DB
func CreateTeam(team *Team) error {
	query := "INSERT INTO teams (name, owner_id) VALUES (?, ?)"
	
	result, err := db.DB.Exec(query, team.Name, team.OwnerID)
	if err != nil {
		log.Println("Error inserting team:", err)
		return err
	}
	
	id, _ := result.LastInsertId()
	team.ID = id
	return nil
}

// AddMember adds a user to a team
func AddMember(teamID, userID int64, status string) error {
	query := "INSERT INTO team_members (team_id, user_id, status) VALUES (?, ?, ?)"

	_, err := db.DB.Exec(query, teamID, userID, status)
	if err != nil {
		log.Println("Error adding team member:", err)
		return err
	}
	return nil
}

// FindMembersByTeamID fetches all members of a team
func FindMembersByTeamID(teamID int64) ([]TeamMemberDetails, error) {
	query := `
		SELECT u.id, tm.team_id, u.first_name, u.last_name, u.email, tm.status
		FROM users u
		JOIN team_members tm ON u.id = tm.user_id
		WHERE tm.team_id = ?
		ORDER BY u.first_name
	`

	rows, err := db.DB.Query(query, teamID)
	if err != nil {
		log.Println("Error querying team members:", err)
		return nil, err
	}
	defer rows.Close()

	var members []TeamMemberDetails
	for rows.Next() {
		var member TeamMemberDetails
		if err := rows.Scan(
			&member.UserID,
			&member.TeamID,
			&member.FirstName,
			&member.LastName,
			&member.Email,
			&member.Status,
		); err != nil {
			log.Println("Error scanning team member:", err)
			continue
		}
		members = append(members, member)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if members == nil {
		members = make([]TeamMemberDetails, 0)
	}

	return members, nil
}
// IsUserMember checks if a user is a 'joined' member of a team
func IsUserMember(teamID int64, userID int64) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM team_members WHERE team_id = ? AND user_id = ? AND status = 'joined'"
	
	err := db.DB.QueryRow(query, teamID, userID).Scan(&count)
	if err != nil {
		log.Println("Error checking team membership:", err)
		return false, err
	}
	return count > 0, nil
}

// FindMemberStatus checks if a user is already in the team (in any status)
func FindMemberStatus(teamID int64, userID int64) (string, error) {
	var status string
	query := "SELECT status FROM team_members WHERE team_id = ? AND user_id = ?"

	err := db.DB.QueryRow(query, teamID, userID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // Not a member, no error
		}
		return "", err // A real database error
	}
	return status, nil // Returns "pending" or "joined"
}
// FindTeamsByUserID fetches all teams (and status) for a user
func FindTeamsByUserID(userID int64) ([]UserTeamDetails, error) {
	query := `
		SELECT t.id, t.name, t.owner_id, tm.status, tm.joined_at
		FROM teams t
		JOIN team_members tm ON t.id = tm.team_id
		WHERE tm.user_id = ?
		ORDER BY t.name
	`

	rows, err := db.DB.Query(query, userID)
	if err != nil {
		log.Println("Error querying user teams:", err)
		return nil, err
	}
	defer rows.Close()

	var teams []UserTeamDetails
	for rows.Next() {
		var team UserTeamDetails
		if err := rows.Scan(
			&team.TeamID,
			&team.TeamName,
			&team.OwnerID,
			&team.UserStatus,
			&team.JoinedAt,
		); err != nil {
			log.Println("Error scanning user team:", err)
			continue
		}
		teams = append(teams, team)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if teams == nil {
		teams = make([]UserTeamDetails, 0) // Return empty slice, not nil
	}

	return teams, nil
}

// UpdateMemberStatus updates a user's status in a specific team
// It ensures a user can only update their status if it's 'pending'
func UpdateMemberStatus(teamID int64, userID int64, newStatus string) error {
	// Only allow updating if the current status is 'pending'
	query := `
		UPDATE team_members 
		SET status = ?, joined_at = CASE WHEN ? = 'joined' THEN NOW() ELSE NULL END
		WHERE team_id = ? AND user_id = ? AND status = 'pending'
	`
	
	result, err := db.DB.Exec(query, newStatus, newStatus, teamID, userID)
	if err != nil {
		log.Println("Error updating member status:", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		// This means they weren't pending or didn't have an invite
		return errors.New("no pending invitation found for this team")
	}

	return nil
}

// We also need a function to handle 'leaving' or 'rejecting' (deleting the row)
func RemoveMember(teamID int64, userID int64) error {
	query := "DELETE FROM team_members WHERE team_id = ? AND user_id = ?"
	
	_, err := db.DB.Exec(query, teamID, userID)
	if err != nil {
		log.Println("Error removing team member:", err)
		return err
	}
	return nil
}

// SaveChatMessage saves a new chat message to the DB
func SaveChatMessage(teamID, userID int64, message string) (*ChatMessage, error) {
	query := "INSERT INTO team_messages (team_id, user_id, message_content) VALUES (?, ?, ?)"
	
	result, err := db.DB.Exec(query, teamID, userID, message)
	if err != nil {
		log.Println("Error saving chat message:", err)
		return nil, err
	}

	id, _ := result.LastInsertId()
	
	// Create a partial ChatMessage object to return (created_at will be set by DB)
	chatMessage := &ChatMessage{
		ID:             id,
		TeamID:         teamID,
		UserID:         userID,
		MessageContent: message,
	}
	return chatMessage, nil
}

// FindChatMessagesByTeamID fetches all messages for a team
func FindChatMessagesByTeamID(teamID int64) ([]ChatMessageDetails, error) {
	query := `
		SELECT m.id, m.team_id, m.user_id, u.first_name, u.last_name, m.message_content, m.created_at
		FROM team_messages m
		JOIN users u ON m.user_id = u.id
		WHERE m.team_id = ?
		ORDER BY m.created_at ASC
	`
	rows, err := db.DB.Query(query, teamID)
	if err != nil {
		log.Println("Error querying chat messages:", err)
		return nil, err
	}
	defer rows.Close()

	var messages []ChatMessageDetails
	for rows.Next() {
		var msg ChatMessageDetails
		if err := rows.Scan(
			&msg.ID,
			&msg.TeamID,
			&msg.UserID,
			&msg.FirstName,
			&msg.LastName,
			&msg.MessageContent,
			&msg.CreatedAt,
		); err != nil {
			log.Println("Error scanning chat message:", err)
			continue
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if messages == nil {
		messages = make([]ChatMessageDetails, 0)
	}

	return messages, nil
}


// CreateMessage saves a new chat message
func CreateMessage(teamID, userID int64, content string) error {
	query := `INSERT INTO team_messages (team_id, user_id, message_content) VALUES (?, ?, ?)`
	_, err := db.DB.Exec(query, teamID, userID, content)
	if err != nil {
		log.Println("Error saving message:", err)
		return err
	}
	return nil
}

// GetMessagesByTeamID fetches chat history for a team
func GetMessagesByTeamID(teamID int64) ([]ChatMessage, error) {
	query := `
		SELECT 
			m.id, m.team_id, m.user_id, u.first_name, u.last_name, m.message_content, m.created_at
		FROM team_messages m
		JOIN users u ON m.user_id = u.id
		WHERE m.team_id = ?
		ORDER BY m.created_at ASC
	`
	rows, err := db.DB.Query(query, teamID)
	if err != nil {
		log.Println("Error fetching messages:", err)
		return nil, err
	}
	defer rows.Close()

	var messages []ChatMessage
	for rows.Next() {
		var msg ChatMessage
		if err := rows.Scan(
			&msg.ID, &msg.TeamID, &msg.UserID, 
			&msg.FirstName, &msg.LastName, 
			&msg.MessageContent, &msg.CreatedAt,
		); err != nil {
			continue
		}
		messages = append(messages, msg)
	}
	
	if messages == nil {
		messages = make([]ChatMessage, 0)
	}
	return messages, nil
}