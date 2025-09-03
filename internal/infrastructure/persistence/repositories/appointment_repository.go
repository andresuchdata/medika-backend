package repositories

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"

	"medika-backend/internal/domain/appointment"
	"medika-backend/internal/infrastructure/persistence/models"
	"medika-backend/pkg/logger"
)

// AppointmentRepository implements appointment.Repository
type AppointmentRepository struct {
	db     bun.IDB
	logger logger.Logger
}

func NewAppointmentRepository(db *bun.DB) appointment.Repository {
	return &AppointmentRepository{
		db:     db,
		logger: logger.New(),
	}
}

func (r *AppointmentRepository) Create(ctx context.Context, apt *appointment.Appointment) error {
	model := r.toModel(apt)

	_, err := r.db.NewInsert().
		Model(model).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to create appointment: %w", err)
	}

	return nil
}

func (r *AppointmentRepository) GetByID(ctx context.Context, id string) (*appointment.Appointment, error) {
	model := &models.Appointment{}
	
	err := r.db.NewSelect().
		Model(model).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get appointment: %w", err)
	}

	return r.toDomain(model), nil
}

func (r *AppointmentRepository) GetByOrganization(ctx context.Context, organizationID string, limit, offset int) ([]*appointment.Appointment, error) {
	var models []models.Appointment
	
	err := r.db.NewSelect().
		Model(&models).
		Where("organization_id = ?", organizationID).
		Limit(limit).
		Offset(offset).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get appointments by organization: %w", err)
	}

	appointments := make([]*appointment.Appointment, len(models))
	for i, model := range models {
		appointments[i] = r.toDomain(&model)
	}

	return appointments, nil
}

func (r *AppointmentRepository) GetByPatient(ctx context.Context, patientID string, limit, offset int) ([]*appointment.Appointment, error) {
	var models []models.Appointment
	
	err := r.db.NewSelect().
		Model(&models).
		Where("patient_id = ?", patientID).
		Limit(limit).
		Offset(offset).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get appointments by patient: %w", err)
	}

	appointments := make([]*appointment.Appointment, len(models))
	for i, model := range models {
		appointments[i] = r.toDomain(&model)
	}

	return appointments, nil
}

func (r *AppointmentRepository) GetByDoctor(ctx context.Context, doctorID string, limit, offset int) ([]*appointment.Appointment, error) {
	var models []models.Appointment
	
	err := r.db.NewSelect().
		Model(&models).
		Where("doctor_id = ?", doctorID).
		Limit(limit).
		Offset(offset).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get appointments by doctor: %w", err)
	}

	appointments := make([]*appointment.Appointment, len(models))
	for i, model := range models {
		appointments[i] = r.toDomain(&model)
	}

	return appointments, nil
}

func (r *AppointmentRepository) Update(ctx context.Context, apt *appointment.Appointment) error {
	model := r.toModel(apt)

	_, err := r.db.NewUpdate().
		Model(model).
		Where("id = ?", apt.ID).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to update appointment: %w", err)
	}

	return nil
}

func (r *AppointmentRepository) UpdateStatus(ctx context.Context, id string, status appointment.AppointmentStatus) error {
	_, err := r.db.NewUpdate().
		Model((*models.Appointment)(nil)).
		Set("status = ?", string(status)).
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to update appointment status: %w", err)
	}

	return nil
}

func (r *AppointmentRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().
		Model((*models.Appointment)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete appointment: %w", err)
	}

	return nil
}

func (r *AppointmentRepository) CountByOrganization(ctx context.Context, organizationID string) (int, error) {
	count, err := r.db.NewSelect().
		Model((*models.Appointment)(nil)).
		Where("organization_id = ?", organizationID).
		Count(ctx)

	if err != nil {
		return 0, fmt.Errorf("failed to count appointments: %w", err)
	}

	return count, nil
}

func (r *AppointmentRepository) toModel(apt *appointment.Appointment) *models.Appointment {
	return &models.Appointment{
		ID:             apt.ID,
		PatientID:      apt.PatientID,
		DoctorID:       apt.DoctorID,
		OrganizationID: apt.OrganizationID,
		RoomID:         apt.RoomID,
		Date:           apt.Date,
		StartTime:      apt.StartTime,
		EndTime:        apt.EndTime,
		Duration:       apt.Duration,
		Status:         string(apt.Status),
		Type:           apt.Type,
		Notes:          apt.Notes,
		CreatedAt:      apt.CreatedAt,
		UpdatedAt:      apt.UpdatedAt,
	}
}

func (r *AppointmentRepository) toDomain(model *models.Appointment) *appointment.Appointment {
	return &appointment.Appointment{
		ID:             model.ID,
		PatientID:      model.PatientID,
		DoctorID:       model.DoctorID,
		OrganizationID: model.OrganizationID,
		RoomID:         model.RoomID,
		Date:           model.Date,
		StartTime:      model.StartTime,
		EndTime:        model.EndTime,
		Duration:       model.Duration,
		Status:         appointment.AppointmentStatus(model.Status),
		Type:           model.Type,
		Notes:          model.Notes,
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      model.UpdatedAt,
	}
}
