// notification/notification_repository.go
package notification

import (
	"log"
	"github.com/JkD004/playarena-backend/db"
)

// CreateNotification inserts a new notification
func CreateNotification(userID int64, message, notifType string) error {
	query := `INSERT INTO notifications (user_id, message, type) VALUES (?, ?, ?)`
	_, err := db.DB.Exec(query, userID, message, notifType)
	if err != nil {
		log.Println("Error creating notification:", err)
		return err
	}
	return nil
}

// GetNotificationsByUserID fetches unread notifications for a user
func GetNotificationsByUserID(userID int64) ([]Notification, error) {
	query := `
		SELECT id, user_id, message, type, is_read, created_at 
		FROM notifications 
		WHERE user_id = ? 
		ORDER BY created_at DESC 
		LIMIT 10
	`
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifs []Notification
	for rows.Next() {
		var n Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Message, &n.Type, &n.IsRead, &n.CreatedAt); err != nil {
			continue
		}
		notifs = append(notifs, n)
	}
	
	if notifs == nil {
		notifs = make([]Notification, 0)
	}
	return notifs, nil
}

// MarkAsRead marks a specific notification as read
func MarkAsRead(notificationID int64) error {
	query := `UPDATE notifications SET is_read = TRUE WHERE id = ?`
	_, err := db.DB.Exec(query, notificationID)
	return err
}