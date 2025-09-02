package shared

import (
	"crypto/rand"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// ID Value Objects
type UserID struct{ value string }
type PatientID struct{ value string }
type DoctorID struct{ value string }
type AppointmentID struct{ value string }
type OrganizationID struct{ value string }
type NotificationID struct{ value string }

// User ID
func NewUserID() UserID {
	return UserID{value: uuid.New().String()}
}

func NewUserIDFromString(s string) (UserID, error) {
	if _, err := uuid.Parse(s); err != nil {
		return UserID{}, fmt.Errorf("invalid user ID: %w", err)
	}
	return UserID{value: s}, nil
}

func (id UserID) String() string { return id.value }
func (id UserID) IsEmpty() bool  { return id.value == "" }

// Patient ID
func NewPatientID() PatientID {
	return PatientID{value: uuid.New().String()}
}

func NewPatientIDFromString(s string) (PatientID, error) {
	if _, err := uuid.Parse(s); err != nil {
		return PatientID{}, fmt.Errorf("invalid patient ID: %w", err)
	}
	return PatientID{value: s}, nil
}

func (id PatientID) String() string { return id.value }
func (id PatientID) IsEmpty() bool  { return id.value == "" }

// Organization ID
func NewOrganizationID(s string) (OrganizationID, error) {
	if _, err := uuid.Parse(s); err != nil {
		return OrganizationID{}, fmt.Errorf("invalid organization ID: %w", err)
	}
	return OrganizationID{value: s}, nil
}

func (id OrganizationID) String() string { return id.value }
func (id OrganizationID) IsEmpty() bool  { return id.value == "" }

// Appointment ID
func NewAppointmentID() AppointmentID {
	return AppointmentID{value: uuid.New().String()}
}

func NewAppointmentIDFromString(s string) (AppointmentID, error) {
	if _, err := uuid.Parse(s); err != nil {
		return AppointmentID{}, fmt.Errorf("invalid appointment ID: %w", err)
	}
	return AppointmentID{value: s}, nil
}

func (id AppointmentID) String() string { return id.value }
func (id AppointmentID) IsEmpty() bool  { return id.value == "" }

// Email Value Object
type Email struct{ value string }

func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return Email{}, fmt.Errorf("invalid email format: %s", email)
	}
	
	return Email{value: email}, nil
}

func (e Email) String() string { return e.value }

// Name Value Object
type Name struct{ value string }

func NewName(name string) (Name, error) {
	name = strings.TrimSpace(name)
	
	if len(name) < 2 {
		return Name{}, fmt.Errorf("name must be at least 2 characters")
	}
	
	if len(name) > 100 {
		return Name{}, fmt.Errorf("name must be less than 100 characters")
	}
	
	// Only letters, spaces, and common name characters
	nameRegex := regexp.MustCompile(`^[a-zA-Z\s\-\'\.]+$`)
	if !nameRegex.MatchString(name) {
		return Name{}, fmt.Errorf("name contains invalid characters")
	}
	
	return Name{value: name}, nil
}

func (n Name) String() string { return n.value }

// Phone Number Value Object
type PhoneNumber struct{ value string }

func NewPhoneNumber(phone string) (PhoneNumber, error) {
	// Remove all non-digits
	cleaned := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")
	
	// Validate length (10-15 digits as per E.164)
	if len(cleaned) < 10 || len(cleaned) > 15 {
		return PhoneNumber{}, fmt.Errorf("invalid phone number length: %s", phone)
	}
	
	// Add + prefix for international format
	formatted := "+" + cleaned
	
	return PhoneNumber{value: formatted}, nil
}

func (p PhoneNumber) String() string { return p.value }

// Blood Type Value Object
type BloodType struct{ value string }

var validBloodTypes = map[string]bool{
	"A+": true, "A-": true, "B+": true, "B-": true,
	"AB+": true, "AB-": true, "O+": true, "O-": true,
}

func NewBloodType(bloodType string) (BloodType, error) {
	bloodType = strings.ToUpper(strings.TrimSpace(bloodType))
	
	if !validBloodTypes[bloodType] {
		return BloodType{}, fmt.Errorf("invalid blood type: %s", bloodType)
	}
	
	return BloodType{value: bloodType}, nil
}

func (bt BloodType) String() string { return bt.value }

func (bt BloodType) IsCompatibleWith(donor BloodType) bool {
	compatibility := map[string][]string{
		"A+":  {"A+", "A-", "O+", "O-"},
		"A-":  {"A-", "O-"},
		"B+":  {"B+", "B-", "O+", "O-"},
		"B-":  {"B-", "O-"},
		"AB+": {"A+", "A-", "B+", "B-", "AB+", "AB-", "O+", "O-"},
		"AB-": {"A-", "B-", "AB-", "O-"},
		"O+":  {"O+", "O-"},
		"O-":  {"O-"},
	}
	
	compatible := compatibility[bt.value]
	for _, c := range compatible {
		if c == donor.value {
			return true
		}
	}
	return false
}

// Gender Value Object
type Gender struct{ value string }

const (
	GenderMale   = "male"
	GenderFemale = "female"
	GenderOther  = "other"
)

func NewGender(gender string) (Gender, error) {
	gender = strings.ToLower(strings.TrimSpace(gender))
	
	switch gender {
	case GenderMale, GenderFemale, GenderOther:
		return Gender{value: gender}, nil
	default:
		return Gender{}, fmt.Errorf("invalid gender: %s", gender)
	}
}

func (g Gender) String() string { return g.value }

// Appointment Time Value Object
type AppointmentTime struct{ value time.Time }

func NewAppointmentTime(t time.Time) (AppointmentTime, error) {
	now := time.Now()
	
	// Must be in the future
	if t.Before(now) {
		return AppointmentTime{}, fmt.Errorf("appointment time must be in the future")
	}
	
	// Must be within business hours (8 AM - 6 PM)
	hour := t.Hour()
	if hour < 8 || hour >= 18 {
		return AppointmentTime{}, fmt.Errorf("appointment must be between 8 AM and 6 PM")
	}
	
	// Must be on weekdays
	if t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
		return AppointmentTime{}, fmt.Errorf("appointments not available on weekends")
	}
	
	// Must be at 15-minute intervals
	if t.Minute()%15 != 0 {
		return AppointmentTime{}, fmt.Errorf("appointments must be scheduled at 15-minute intervals")
	}
	
	return AppointmentTime{value: t}, nil
}

func (at AppointmentTime) Value() time.Time { return at.value }
func (at AppointmentTime) IsToday() bool {
	return at.value.Format("2006-01-02") == time.Now().Format("2006-01-02")
}

// Password handling utilities
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Generate secure random string
func GenerateSecureToken(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}
	return string(bytes), nil
}
