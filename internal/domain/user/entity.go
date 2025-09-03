package user

import (
	"errors"
	"time"

	"medika-backend/internal/domain/shared"
	"medika-backend/pkg/crypto"
)

// User aggregate root
type User struct {
	id             shared.UserID
	email          shared.Email
	name           shared.Name
	passwordHash   string
	role           Role
	organizationID *shared.OrganizationID
	phone          *shared.PhoneNumber
	avatarURL      *string
	isActive       bool
	profile        *Profile
	createdAt      time.Time
	updatedAt      time.Time
	version        int
}

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleDoctor   Role = "doctor"
	RoleNurse    Role = "nurse"
	RolePatient  Role = "patient"
	RoleCashier  Role = "cashier"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleDoctor, RoleNurse, RolePatient, RoleCashier:
		return true
	default:
		return false
	}
}

func (r Role) IsStaff() bool {
	return r == RoleAdmin || r == RoleDoctor || r == RoleNurse || r == RoleCashier
}

func (r Role) String() string {
	return string(r)
}

type Profile struct {
	userID          shared.UserID
	dateOfBirth     *time.Time
	gender          *shared.Gender
	address         *string
	emergencyContact *string
	medicalHistory  *string
	allergies       []string
	bloodType       *shared.BloodType
	// Doctor-specific fields
	specialty       *string
	licenseNumber   *string
	bio             *string
	experience      *int
	education       []string
	certifications  []string
	nextAvailable   *string
}

// Getter methods for Profile
func (p *Profile) UserID() shared.UserID           { return p.userID }
func (p *Profile) DateOfBirth() *time.Time         { return p.dateOfBirth }
func (p *Profile) Gender() *shared.Gender          { return p.gender }
func (p *Profile) Address() *string                { return p.address }
func (p *Profile) EmergencyContact() *string       { return p.emergencyContact }
func (p *Profile) MedicalHistory() *string         { return p.medicalHistory }
func (p *Profile) Allergies() []string             { return p.allergies }
func (p *Profile) BloodType() *shared.BloodType    { return p.bloodType }
// Doctor-specific getters
func (p *Profile) Specialty() *string              { return p.specialty }
func (p *Profile) LicenseNumber() *string          { return p.licenseNumber }
func (p *Profile) Bio() *string                    { return p.bio }
func (p *Profile) Experience() *int                { return p.experience }
func (p *Profile) Education() []string             { return p.education }
func (p *Profile) Certifications() []string        { return p.certifications }
func (p *Profile) NextAvailable() *string          { return p.nextAvailable }

// Constructor
func NewUser(
	email, name, password string,
	role Role,
	organizationID *string,
) (*User, error) {
	// Validate inputs
	userEmail, err := shared.NewEmail(email)
	if err != nil {
		return nil, err
	}

	userName, err := shared.NewName(name)
	if err != nil {
		return nil, err
	}

	if !role.IsValid() {
		return nil, errors.New("invalid user role")
	}

	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	// Hash password
	hashedPassword, err := crypto.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Handle organization ID
	var orgID *shared.OrganizationID
	if organizationID != nil {
		id, err := shared.NewOrganizationID(*organizationID)
		if err != nil {
			return nil, err
		}
		orgID = &id
	}

	// Business rule: Patients don't require organization ID initially
	if role != RolePatient && orgID == nil {
		return nil, errors.New("organization ID is required for staff members")
	}

	return &User{
		id:             shared.NewUserID(),
		email:          userEmail,
		name:           userName,
		passwordHash:   hashedPassword,
		role:           role,
		organizationID: orgID,
		isActive:       true,
		createdAt:      time.Now(),
		updatedAt:      time.Now(),
		version:        1,
	}, nil
}

// Domain methods
func (u *User) ID() shared.UserID                    { return u.id }
func (u *User) Email() shared.Email                  { return u.email }
func (u *User) Name() shared.Name                    { return u.name }
func (u *User) Role() Role                           { return u.role }
func (u *User) OrganizationID() *shared.OrganizationID { return u.organizationID }
func (u *User) Phone() *shared.PhoneNumber           { return u.phone }
func (u *User) AvatarURL() *string                   { return u.avatarURL }
func (u *User) IsActive() bool                       { return u.isActive }
func (u *User) Profile() *Profile                    { return u.profile }
func (u *User) CreatedAt() time.Time                 { return u.createdAt }
func (u *User) UpdatedAt() time.Time                 { return u.updatedAt }
func (u *User) Version() int                         { return u.version }
func (u *User) Status() string                       { 
	if u.isActive {
		return "active"
	}
	return "inactive"
}

func (u *User) PasswordHash() string {
	return u.passwordHash
}

func (u *User) VerifyPassword(password string) bool {
	return crypto.VerifyPassword(password, u.passwordHash)
}

func (u *User) UpdateProfile(
	dateOfBirth *time.Time,
	gender *string,
	address *string,
	phone *string,
) error {
	// Validate and convert inputs
	var phoneNumber *shared.PhoneNumber
	if phone != nil {
		pn, err := shared.NewPhoneNumber(*phone)
		if err != nil {
			return err
		}
		phoneNumber = &pn
	}

	var genderValue *shared.Gender
	if gender != nil {
		g, err := shared.NewGender(*gender)
		if err != nil {
			return err
		}
		genderValue = &g
	}

	// Update phone on user
	u.phone = phoneNumber

	// Create or update profile
	if u.profile == nil {
		u.profile = &Profile{
			userID: u.id,
		}
	}

	u.profile.dateOfBirth = dateOfBirth
	u.profile.gender = genderValue
	u.profile.address = address
	u.updateTimestamp()

	return nil
}

func (u *User) UpdateMedicalInfo(
	emergencyContact *string,
	medicalHistory *string,
	allergies []string,
	bloodType *string,
) error {
	// Only patients and medical staff can have medical info
	if u.role != RolePatient && !u.role.IsStaff() {
		return errors.New("medical information not applicable for this user role")
	}

	var bloodTypeValue *shared.BloodType
	if bloodType != nil {
		bt, err := shared.NewBloodType(*bloodType)
		if err != nil {
			return err
		}
		bloodTypeValue = &bt
	}

	if u.profile == nil {
		u.profile = &Profile{
			userID: u.id,
		}
	}

	u.profile.emergencyContact = emergencyContact
	u.profile.medicalHistory = medicalHistory
	u.profile.allergies = allergies
	u.profile.bloodType = bloodTypeValue
	u.updateTimestamp()

	return nil
}

func (u *User) UpdateDoctorProfile(
	specialty *string,
	licenseNumber *string,
	bio *string,
	experience *int,
	education []string,
	certifications []string,
	nextAvailable *string,
) error {
	// Only doctors can have doctor profile information
	if u.role != RoleDoctor {
		return errors.New("doctor profile information not applicable for this user role")
	}

	if u.profile == nil {
		u.profile = &Profile{
			userID: u.id,
		}
	}

	u.profile.specialty = specialty
	u.profile.licenseNumber = licenseNumber
	u.profile.bio = bio
	u.profile.experience = experience
	u.profile.education = education
	u.profile.certifications = certifications
	u.profile.nextAvailable = nextAvailable
	u.updateTimestamp()

	return nil
}

func (u *User) UpdateAvatar(avatarURL string) error {
	// Validate avatar URL (basic validation)
	if avatarURL == "" {
		return errors.New("avatar URL cannot be empty")
	}
	
	// Basic URL format validation
	if len(avatarURL) > 500 {
		return errors.New("avatar URL too long")
	}

	u.avatarURL = &avatarURL
	u.updateTimestamp()
	
	return nil
}

func (u *User) Deactivate() {
	u.isActive = false
	u.updateTimestamp()
}

func (u *User) Activate() {
	u.isActive = true
	u.updateTimestamp()
}

func (u *User) ChangePassword(newPassword string) error {
	if len(newPassword) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	hashedPassword, err := shared.HashPassword(newPassword)
	if err != nil {
		return err
	}

	u.passwordHash = hashedPassword
	u.updateTimestamp()

	return nil
}

func (u *User) updateTimestamp() {
	u.updatedAt = time.Now()
	u.version++
}

// Business rules
func (u *User) CanAccessPatient(patientID shared.UserID, patientOrgID shared.OrganizationID) bool {
	switch u.role {
	case RoleAdmin:
		// Admin can access all patients in their organization
		return u.organizationID != nil && *u.organizationID == patientOrgID
	case RoleDoctor, RoleNurse:
		// Medical staff can access patients in their organization
		return u.organizationID != nil && *u.organizationID == patientOrgID
	case RoleCashier:
		// Cashier can access patients for billing purposes
		return u.organizationID != nil && *u.organizationID == patientOrgID
	case RolePatient:
		// Patients can only access their own data
		return u.id == patientID
	default:
		return false
	}
}

func (u *User) CanManageAppointments() bool {
	return u.role.IsStaff()
}

func (u *User) CanAccessMedicalRecords(patientID shared.UserID, patientOrgID shared.OrganizationID) bool {
	switch u.role {
	case RoleDoctor, RoleNurse:
		// Medical staff can access medical records
		return u.organizationID != nil && *u.organizationID == patientOrgID
	case RolePatient:
		// Patients can access their own medical records
		return u.id == patientID
	default:
		return false
	}
}

// Reconstruction for repository
func ReconstructUser(
	id shared.UserID,
	email shared.Email,
	name shared.Name,
	passwordHash string,
	role Role,
	organizationID *shared.OrganizationID,
	phone *shared.PhoneNumber,
	avatarURL *string,
	isActive bool,
	profile *Profile,
	createdAt, updatedAt time.Time,
	version int,
) *User {
	return &User{
		id:             id,
		email:          email,
		name:           name,
		passwordHash:   passwordHash,
		role:           role,
		organizationID: organizationID,
		phone:          phone,
		avatarURL:      avatarURL,
		isActive:       isActive,
		profile:        profile,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
		version:        version,
	}
}
