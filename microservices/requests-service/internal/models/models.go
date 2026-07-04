package models

import "time"

const (
	RequestStatusNew        = "new"
	RequestStatusInProgress = "in_progress"
	RequestStatusDone       = "done"
	RequestStatusCanceled   = "canceled"
)

type UserContext struct {
	UserID string
	Role   string
}

type Pagination struct {
	Page  int
	Limit int
}

type MaintenanceRequest struct {
	ID          string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID      string `gorm:"type:uuid;index;not null"`
	Category    string `gorm:"size:64;index;not null"`
	Description string `gorm:"not null"`
	Status      string `gorm:"size:32;index;not null"`
	AssignedTo  string `gorm:"type:uuid;index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type RequestComment struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	RequestID string `gorm:"type:uuid;index;not null"`
	UserID    string `gorm:"type:uuid;index;not null"`
	Text      string `gorm:"not null"`
	CreatedAt time.Time
}

type CreateRequestCommand struct {
	User        UserContext
	Category    string
	Description string
}

type ListRequestsFilter struct {
	Actor      UserContext
	Pagination Pagination
}

type RequestsPage struct {
	Items []MaintenanceRequest
	Page  int
	Limit int
	Total int64
}

type GetRequestCommand struct {
	Actor     UserContext
	RequestID string
}

type UpdateRequestStatusCommand struct {
	Actor      UserContext
	RequestID  string
	Status     string
	AssignedTo string
}

type AddRequestCommentCommand struct {
	Actor     UserContext
	RequestID string
	Text      string
}

type NotificationEvent struct {
	UserID   string
	Type     string
	Title    string
	Body     string
	EntityID string
}
