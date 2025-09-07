package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/uptrace/bun"
)

// Notification represents a notification in the database
type Notification struct {
	bun.BaseModel `bun:"table:notifications"`

	ID           string    `bun:"id,pk" json:"id"`
	UserID       string    `bun:"user_id,notnull" json:"user_id"`
	Type         string    `bun:"type,notnull" json:"type"`
	Title        string    `bun:"title,notnull" json:"title"`
	Message      string    `bun:"message,notnull" json:"message"`
	Data         JSONB     `bun:"data,type:jsonb" json:"data"`
	IsRead       bool      `bun:"is_read,default:false" json:"is_read"`
	Channels     []string  `bun:"channels,array" json:"channels"`
	Priority     string    `bun:"priority,notnull" json:"priority"`
	ScheduledFor *time.Time `bun:"scheduled_for" json:"scheduled_for"`
	SentAt       *time.Time `bun:"sent_at" json:"sent_at"`
	CreatedAt    time.Time `bun:"created_at,notnull,default:now()" json:"created_at"`

	// Relations
	User *User `bun:"rel:belongs-to,join:user_id=id" json:"user,omitempty"`
}

// JSONB is a custom type for handling JSONB columns
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, j)
}
