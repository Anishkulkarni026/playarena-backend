// team/team_handler.go
package team

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"
)


// Struct for the status update request
type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"` // "joined" or "rejected"
}

// UpdateMemberStatusHandler handles PATCH /api/v1/teams/:id/status
func UpdateMemberStatusHandler(c *gin.Context) {
	// Get team ID from URL
	teamID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	// Get user's ID from token (this is the user updating their *own* status)
	userID := c.MustGet("userID").(int64)

	// Get the new status from the JSON body
	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, 'status' is required"})
		return
	}

	// Call the service
	err = UpdateUserTeamStatus(userID, teamID, req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Team status updated successfully"})
}

// Struct for the invite request body
type InviteRequest struct {
	Email string `json:"email" binding:"required"`
}

// InviteMemberHandler handles POST /api/v1/teams/:id/invite
func InviteMemberHandler(c *gin.Context) {
	// Get team ID from URL
	teamID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	// Get inviter's ID from token
	inviterID := c.MustGet("userID").(int64)

	// Get invitee's email from JSON body
	var req InviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, 'email' is required"})
		return
	}

	// Call the service
	err = InviteMember(inviterID, teamID, req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation sent successfully"})
}


// GetTeamMembersHandler handles GET /api/v1/teams/:id/members
func GetTeamMembersHandler(c *gin.Context) {
	// Get team ID from URL
	teamID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	// Get user's ID from token
	userID := c.MustGet("userID").(int64)

	// Call the service
	members, err := GetTeamMembers(userID, teamID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, members)
}

// CreateTeamHandler handles POST /api/v1/teams
func CreateTeamHandler(c *gin.Context) {
	var req CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Get userID from the middleware context
	userID := c.MustGet("userID").(int64)

	team, err := CreateNewTeam(&req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, team)
}

// GetMyTeamsHandler handles GET /api/v1/teams/mine
func GetMyTeamsHandler(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	teams, err := GetTeamsForUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch your teams"})
		return
	}

	c.JSON(http.StatusOK, teams)
}

// PostMessageHandler handles POST /api/v1/teams/:id/chat

// GetTeamMessagesHandler handles GET /api/v1/teams/:id/chat
// PostMessageHandler handles sending a message
func PostMessageHandler(c *gin.Context) {
	teamID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}
	
	userID := c.MustGet("userID").(int64)

	var req struct {
		Message string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message content required"})
		return
	}

	err = SendTeamMessage(teamID, userID, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Message sent"})
}

// GetTeamMessagesHandler handles fetching chat history
func GetTeamMessagesHandler(c *gin.Context) {
	teamID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	messages, err := GetTeamChat(teamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chat"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

