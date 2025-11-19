// team/team_service.go
package team

import (
	"errors"
	"log"
	"github.com/JkD004/playarena-backend/user"
	"github.com/JkD004/playarena-backend/notification"
)

// CreateNewTeam handles the logic for creating a new team
func CreateNewTeam(req *CreateTeamRequest, ownerID int64) (*Team, error) {
	if req.Name == "" {
		return nil, errors.New("team name cannot be empty")
	}

	newTeam := &Team{
		Name:    req.Name,
		OwnerID: ownerID,
	}

	// 1. Save the new team to the 'teams' table
	err := CreateTeam(newTeam)
	if err != nil {
		log.Println("Service error creating team:", err)
		return nil, errors.New("could not create team")
	}

	// 2. Automatically add the owner as a 'joined' member
	err = AddMember(newTeam.ID, ownerID, "joined")
	if err != nil {
		log.Printf("CRITICAL: Failed to add owner %d as member to new team %d: %v\n", ownerID, newTeam.ID, err)
	}

	return newTeam, nil
}

// InviteMember handles the logic for inviting a new member
func InviteMember(inviterID, teamID int64, inviteeEmail string) error {
	// 1. Check if the person sending the invite is a 'joined' member
	// We assume IsUserMember checks if they are "joined" (not pending)
	isMember, err := IsUserMember(teamID, inviterID)
	if err != nil {
		return errors.New("database error")
	}
	if !isMember {
		return errors.New("you are not a member of this team")
	}

	// 2. Find the user being invited
	invitee, err := user.GetUserByEmail(inviteeEmail)
	if err != nil {
		return errors.New("user with that email not found")
	}

	// 3. Check if the invitee is the same as the inviter
	if invitee.ID == inviterID {
		return errors.New("you cannot invite yourself")
	}

	// 4. Check if the user is already in the team
	status, err := FindMemberStatus(teamID, invitee.ID)
	if err != nil {
		// If error is not nil, it might just mean they aren't found, which is good.
		// But if FindMemberStatus returns "not found" as an error, we handle it.
		// Assuming FindMemberStatus returns "" string if not found.
	}
	
	if status == "joined" {
		return errors.New("user is already a member of this team")
	}
	if status == "pending" {
		return errors.New("user has already been invited")
	}

	// 5. Add the new member with 'pending' status
	// NOTE: This uses AddMember from repository, not AddTeamMember
	err = AddMember(teamID, invitee.ID, "pending")
	if err != nil {
		return errors.New("failed to invite user")
	}

	// 6. Send Notification
	_ = notification.CreateNotification(invitee.ID, "You have been invited to join a new team!", "info")
	
	return nil
}

// GetTeamsForUser fetches all teams for a given user
func GetTeamsForUser(userID int64) ([]UserTeamDetails, error) {
	return FindTeamsByUserID(userID)
}

// UpdateUserTeamStatus handles the logic for a user accepting/rejecting an invite
func UpdateUserTeamStatus(userID, teamID int64, newStatus string) error {
	if newStatus == "joined" {
		return UpdateMemberStatus(teamID, userID, "joined")
	} 
	
	if newStatus == "rejected" {
		return RemoveMember(teamID, userID)
	}

	return errors.New("invalid status: must be 'joined' or 'rejected'")
}

// GetTeamMembers fetches all members for a team
func GetTeamMembers(userID, teamID int64) ([]TeamMemberDetails, error) {
	// Check if the user requesting the list is a member
	isMember, err := IsUserMember(teamID, userID)
	if err != nil {
		return nil, errors.New("database error checking membership")
	}
	if !isMember {
		return nil, errors.New("you are not a member of this team")
	}

	return FindMembersByTeamID(teamID)
}

// PostChatMessage handles logic for sending a new message
func PostChatMessage(userID, teamID int64, message string) (*ChatMessage, error) {
	isMember, err := IsUserMember(teamID, userID)
	if err != nil {
		return nil, errors.New("database error checking membership")
	}
	if !isMember {
		return nil, errors.New("you are not a member of this team")
	}

	if message == "" {
		return nil, errors.New("message cannot be empty")
	}

	return SaveChatMessage(teamID, userID, message)
}

// GetTeamMessages fetches all chat messages for a team
func GetTeamMessages(userID, teamID int64) ([]ChatMessageDetails, error) {
	isMember, err := IsUserMember(teamID, userID)
	if err != nil {
		return nil, errors.New("database error checking membership")
	}
	if !isMember {
		return nil, errors.New("you are not a member of this team")
	}

	return FindChatMessagesByTeamID(teamID)
}

// --- Legacy/Duplicate Wrappers (Kept to prevent router errors) ---

func InviteUserToTeam(email string, teamID int64) error {
    // This is a wrapper for the old router call, but it lacks inviterID.
    // Ideally, update your handler to call InviteMember with the inviter's ID.
    // For now, we just call the user lookup logic.
    userToInvite, err := user.GetUserByEmail(email)
    if (err != nil) { return err }
    return AddMember(teamID, userToInvite.ID, "pending")
}

func GetMembersForTeam(teamID int64) ([]TeamMemberDetails, error) {
    return FindMembersByTeamID(teamID)
}

func SendTeamMessage(teamID, userID int64, content string) error {
	return CreateMessage(teamID, userID, content)
}

func GetTeamChat(teamID int64) ([]ChatMessage, error) {
	return GetMessagesByTeamID(teamID)
}