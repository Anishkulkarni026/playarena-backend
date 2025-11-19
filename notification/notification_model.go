// notification/notification_model.go
package notification

import "time"

type Notification struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Message   string    `json:"message"`
	Type      string    `json:"type"` // 'info', 'success', etc.
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}