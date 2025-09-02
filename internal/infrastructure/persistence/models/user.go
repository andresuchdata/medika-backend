package models

import (
	"time"

	"github.com/uptrace/bun"
)

// User model for Bun ORM
type User struct {
	bun.BaseModel `bun:"table:users"`

	ID             string    `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	Email          string    `bun:"email,unique,notnull"`
	Name           string    `bun:"name,notnull"`
	PasswordHash   string    `bun:"password_hash,notnull"`
	Role           string    `bun:"role,notnull"`
	OrganizationID *string   `bun:"organization_id,type:uuid"`
	Phone          *string   `bun:"phone"`
	AvatarURL      *string   `bun:"avatar_url"`
	IsActive       bool      `bun:"is_active,default:true"`
	CreatedAt      time.Time `bun:"created_at,default:current_timestamp"`
	UpdatedAt      time.Time `bun:"updated_at,default:current_timestamp"`
	Version        int       `bun:"version,default:1"`

	// Relations
	Organization *Organization `bun:"rel:belongs-to,join:organization_id=id"`
	Profile      *UserProfile  `bun:"rel:has-one,join:id=user_id"`
	Appointments []Appointment `bun:"rel:has-many,join:id=patient_id"`
}

// UserProfile model
type UserProfile struct {
	bun.BaseModel `bun:"table:user_profiles"`

	ID               string     `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	UserID           string     `bun:"user_id,type:uuid,unique,notnull"`
	DateOfBirth      *time.Time `bun:"date_of_birth"`
	Gender           *string    `bun:"gender"`
	Address          *string    `bun:"address"`
	EmergencyContact *string    `bun:"emergency_contact"`
	MedicalHistory   *string    `bun:"medical_history"`
	Allergies        []string   `bun:"allergies,type:text[]"`
	BloodType        *string    `bun:"blood_type"`
	CreatedAt        time.Time  `bun:"created_at,default:current_timestamp"`
	UpdatedAt        time.Time  `bun:"updated_at,default:current_timestamp"`

	// Relations
	User *User `bun:"rel:belongs-to,join:user_id=id"`
}

// Organization model
type Organization struct {
	bun.BaseModel `bun:"table:organizations"`

	ID            string          `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	Name          string          `bun:"name,notnull"`
	Type          string          `bun:"type,notnull"`
	Address       string          `bun:"address,notnull"`
	Phone         string          `bun:"phone,notnull"`
	Email         string          `bun:"email,notnull"`
	Website       *string         `bun:"website"`
	BusinessHours []BusinessHours `bun:"business_hours,type:jsonb"`
	IsActive      bool            `bun:"is_active,default:true"`
	CreatedAt     time.Time       `bun:"created_at,default:current_timestamp"`
	UpdatedAt     time.Time       `bun:"updated_at,default:current_timestamp"`

	// Relations
	Users []User `bun:"rel:has-many,join:id=organization_id"`
}

// BusinessHours embedded struct
type BusinessHours struct {
	Day    int    `json:"day"`    // 0-6 (Sunday-Saturday)
	Open   string `json:"open"`   // "09:00"
	Close  string `json:"close"`  // "17:00"
	IsOpen bool   `json:"is_open"`
}

// Appointment model
type Appointment struct {
	bun.BaseModel `bun:"table:appointments"`

	ID             string    `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	PatientID      string    `bun:"patient_id,type:uuid,notnull"`
	DoctorID       string    `bun:"doctor_id,type:uuid,notnull"`
	OrganizationID string    `bun:"organization_id,type:uuid,notnull"`
	RoomID         *string   `bun:"room_id,type:uuid"`
	Date           time.Time `bun:"date,notnull"`
	StartTime      string    `bun:"start_time,notnull"`
	EndTime        string    `bun:"end_time,notnull"`
	Duration       int       `bun:"duration,notnull"` // minutes
	Status         string    `bun:"status,notnull"`
	Type           string    `bun:"type,notnull"`
	Notes          *string   `bun:"notes"`
	CreatedAt      time.Time `bun:"created_at,default:current_timestamp"`
	UpdatedAt      time.Time `bun:"updated_at,default:current_timestamp"`

	// Relations
	Patient      *User         `bun:"rel:belongs-to,join:patient_id=id"`
	Doctor       *User         `bun:"rel:belongs-to,join:doctor_id=id"`
	Organization *Organization `bun:"rel:belongs-to,join:organization_id=id"`
	Room         *Room         `bun:"rel:belongs-to,join:room_id=id"`
}

// Room model
type Room struct {
	bun.BaseModel `bun:"table:rooms"`

	ID             string    `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	OrganizationID string    `bun:"organization_id,type:uuid,notnull"`
	Name           string    `bun:"name,notnull"`
	Type           string    `bun:"type,notnull"`
	Capacity       int       `bun:"capacity,notnull"`
	IsAvailable    bool      `bun:"is_available,default:true"`
	Equipment      []string  `bun:"equipment,type:text[]"`
	CreatedAt      time.Time `bun:"created_at,default:current_timestamp"`
	UpdatedAt      time.Time `bun:"updated_at,default:current_timestamp"`

	// Relations
	Organization *Organization `bun:"rel:belongs-to,join:organization_id=id"`
	Appointments []Appointment `bun:"rel:has-many,join:id=room_id"`
}

// PatientQueue model
type PatientQueue struct {
	bun.BaseModel `bun:"table:patient_queues"`

	ID                string    `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	AppointmentID     string    `bun:"appointment_id,type:uuid,notnull"`
	OrganizationID    string    `bun:"organization_id,type:uuid,notnull"`
	Position          int       `bun:"position,notnull"`
	EstimatedWaitTime int       `bun:"estimated_wait_time,notnull"` // minutes
	Status            string    `bun:"status,notnull"`
	CreatedAt         time.Time `bun:"created_at,default:current_timestamp"`
	UpdatedAt         time.Time `bun:"updated_at,default:current_timestamp"`

	// Relations
	Appointment  *Appointment  `bun:"rel:belongs-to,join:appointment_id=id"`
	Organization *Organization `bun:"rel:belongs-to,join:organization_id=id"`
}

// Notification model
type Notification struct {
	bun.BaseModel `bun:"table:notifications"`

	ID           string                 `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	UserID       string                 `bun:"user_id,type:uuid,notnull"`
	Type         string                 `bun:"type,notnull"`
	Title        string                 `bun:"title,notnull"`
	Message      string                 `bun:"message,notnull"`
	Data         map[string]interface{} `bun:"data,type:jsonb"`
	IsRead       bool                   `bun:"is_read,default:false"`
	Channels     []string               `bun:"channels,type:text[]"`
	Priority     string                 `bun:"priority,notnull"`
	ScheduledFor *time.Time             `bun:"scheduled_for"`
	SentAt       *time.Time             `bun:"sent_at"`
	CreatedAt    time.Time              `bun:"created_at,default:current_timestamp"`

	// Relations
	User *User `bun:"rel:belongs-to,join:user_id=id"`
}

// Media model
type Media struct {
	bun.BaseModel `bun:"table:media"`

	ID               string                 `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	Filename         string                 `bun:"filename,notnull"`
	OriginalName     string                 `bun:"original_name,notnull"`
	MimeType         string                 `bun:"mime_type,notnull"`
	Size             int64                  `bun:"size,notnull"`
	URL              string                 `bun:"url,notnull"`
	ThumbnailURL     *string                `bun:"thumbnail_url"`
	UploadedBy       string                 `bun:"uploaded_by,type:uuid,notnull"`
	OrganizationID   string                 `bun:"organization_id,type:uuid,notnull"`
	Metadata         map[string]interface{} `bun:"metadata,type:jsonb"`
	ProcessingStatus string                 `bun:"processing_status,notnull"`
	CreatedAt        time.Time              `bun:"created_at,default:current_timestamp"`
	UpdatedAt        time.Time              `bun:"updated_at,default:current_timestamp"`

	// Relations
	UploadedByUser *User         `bun:"rel:belongs-to,join:uploaded_by=id"`
	Organization   *Organization `bun:"rel:belongs-to,join:organization_id=id"`
}
