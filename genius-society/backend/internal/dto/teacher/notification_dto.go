package teacherdto

import "time"

// NotificationDTO represents a notification sent to a user.
type NotificationDTO struct {
	ID        string     `json:"id"`
	UserID    string     `json:"userId"`
	Type      string     `json:"type"`
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	RefType   string     `json:"refType,omitempty"`
	RefID     string     `json:"refId,omitempty"`
	IsRead    bool       `json:"isRead"`
	ReadAt    *time.Time `json:"readAt,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
}

// UnreadCountDTO holds the total unread notifications count.
type UnreadCountDTO struct {
	Count int64 `json:"count"`
}

// NotificationMeta holds pagination metadata for notifications.
type NotificationMeta struct {
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"totalPages"`
}
