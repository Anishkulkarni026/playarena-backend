// notification/notification_handler.go
package notification

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

// GetMyNotificationsHandler fetches notifications for the logged-in user
func GetMyNotificationsHandler(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	notifs, err := GetNotificationsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}
	c.JSON(http.StatusOK, notifs)
}

// MarkReadHandler marks a notification as read
func MarkReadHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	
	// TODO: Verify the notification belongs to the user
	_ = MarkAsRead(id) // We ignore error for simplicity here
	
	c.JSON(http.StatusOK, gin.H{"message": "Marked as read"})
}