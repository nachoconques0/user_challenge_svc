package event

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

const (
	// UserCreated event type means that the user was created
	UserCreated string = "USER_CREATED"
	// UserUpdated event type means that the user was updated
	UserUpdated string = "USER_UPDATED"
	// UserSoftDeleted event type means that the user was soft deleted
	UserSoftDeleted string = "USER_SOFT_DELETED"
)

// User event type represnt the entity that will be stored when user has any changes
type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID      `gorm:"not null"`
	EventType string         `gorm:"not null"`
	Payload   datatypes.JSON `gorm:"type:jsonb;not null"`
	Published bool           `gorm:"not null;default:false"`
	CreatedAt time.Time
}

// TableName returns the user event table
func (User) TableName() string {
	return "challenge.user_event"
}

type CreatedPayload struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	TraceID  string `json:"trace_id"`
}

type UpdatedPayload struct {
	UserID   string `json:"user_id"`
	Nickname string `json:"nickname"`
	TraceID  string `json:"trace_id"`
}

type DeletedPayload struct {
	UserID  string `json:"user_id"`
	TraceID string `json:"trace_id"`
}
