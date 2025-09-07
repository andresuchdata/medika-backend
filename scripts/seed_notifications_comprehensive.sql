-- Comprehensive notification seed data with consistent linked IDs
-- This script creates notifications for ALL user types in the existing database
-- Ensures every user (admin, doctor, nurse, cashier, patient) has relevant notifications

-- Clear existing notifications (if any)
DELETE FROM notifications;

-- Insert comprehensive notification data for ALL user types
INSERT INTO notifications (
    id,
    user_id,
    title,
    message,
    type,
    priority,
    is_read,
    channels,
    data,
    scheduled_for,
    sent_at,
    created_at
) VALUES 
-- Appointment-related notifications for Dr. Sarah Johnson (admin)
(
    '11190001-1111-1111-1111-111111111111',
    '11120001-1111-1111-1111-111111111111', -- Dr. Sarah Johnson
    'New Appointment Scheduled',
    'A new consultation appointment has been scheduled for tomorrow at 10:00 AM with patient John Doe.',
    'appointment',
    'medium',
    false,
    ARRAY['email', 'in_app'],
    '{"appointment_id": "11170001-1111-1111-1111-111111111111", "patient_name": "John Doe", "appointment_time": "2025-01-08T10:00:00Z", "action_required": false}',
    NULL,
    NOW(),
    NOW() - INTERVAL '2 hours'
),

-- High priority appointment notification
(
    '11190002-1111-1111-1111-111111111111',
    '11120001-1111-1111-1111-111111111111', -- Dr. Sarah Johnson
    'Urgent Appointment Request',
    'Patient Jane Smith has requested an urgent consultation. Please review and confirm availability.',
    'appointment',
    'high',
    false,
    ARRAY['email', 'in_app', 'sms'],
    '{"appointment_id": "11170002-1111-1111-1111-111111111111", "patient_name": "Jane Smith", "urgency": "high", "action_required": true}',
    NULL,
    NOW(),
    NOW() - INTERVAL '1 hour'
),

-- Queue-related notification
(
    '11190003-1111-1111-1111-111111111111',
    '11120001-1111-1111-1111-111111111111', -- Dr. Sarah Johnson
    'Patient Ready for Consultation',
    'Patient in queue position 1 is ready for consultation. Please call them to your office.',
    'alert',
    'medium',
    false,
    ARRAY['in_app'],
    '{"queue_id": "11180001-1111-1111-1111-111111111111", "appointment_id": "11170001-1111-1111-1111-111111111111", "position": 1, "patient_name": "John Doe", "action_required": true}',
    NULL,
    NOW(),
    NOW() - INTERVAL '30 minutes'
),

-- System notification
(
    '11190004-1111-1111-1111-111111111111',
    '11120001-1111-1111-1111-111111111111', -- Dr. Sarah Johnson
    'System Maintenance Scheduled',
    'Scheduled system maintenance will occur tonight from 11:00 PM to 1:00 AM. Please save your work.',
    'system',
    'low',
    true,
    ARRAY['email', 'in_app'],
    '{"maintenance_start": "2025-01-07T23:00:00Z", "maintenance_end": "2025-01-08T01:00:00Z", "affected_services": ["database", "file_storage"]}',
    NULL,
    NOW() - INTERVAL '4 hours',
    NOW() - INTERVAL '4 hours'
),

-- Lab results notification
(
    '11190005-1111-1111-1111-111111111111',
    '11120001-1111-1111-1111-111111111111', -- Dr. Sarah Johnson
    'Lab Results Available',
    'Blood test results for patient John Doe are now available for review.',
    'lab',
    'medium',
    false,
    ARRAY['email', 'in_app'],
    '{"patient_id": "11150001-1111-1111-1111-111111111111", "patient_name": "John Doe", "test_type": "blood_test", "results_status": "normal", "action_required": true}',
    NULL,
    NOW(),
    NOW() - INTERVAL '15 minutes'
),

-- Emergency notification (critical priority)
(
    '11190006-1111-1111-1111-111111111111',
    '11120001-1111-1111-1111-111111111111', -- Dr. Sarah Johnson
    'Emergency Alert',
    'Emergency patient has arrived and requires immediate attention. Please respond immediately.',
    'emergency',
    'critical',
    false,
    ARRAY['email', 'in_app', 'sms', 'phone'],
    '{"emergency_type": "cardiac", "patient_name": "Emergency Patient", "severity": "critical", "location": "Emergency Room", "action_required": true}',
    NULL,
    NOW(),
    NOW() - INTERVAL '5 minutes'
),

-- Notifications for Dr. Michael Smith
(
    '11190007-1111-1111-1111-111111111111',
    '11130001-1111-1111-1111-111111111111', -- Dr. Michael Smith
    'Appointment Reminder',
    'You have a follow-up appointment with patient Jane Smith in 30 minutes.',
    'appointment',
    'medium',
    false,
    ARRAY['in_app'],
    '{"appointment_id": "11170002-1111-1111-1111-111111111111", "patient_name": "Jane Smith", "appointment_time": "2025-01-07T14:30:00Z", "action_required": false}',
    NULL,
    NOW(),
    NOW() - INTERVAL '10 minutes'
),

-- Patient check-in notification
(
    '11190008-1111-1111-1111-111111111111',
    '11130001-1111-1111-1111-111111111111', -- Dr. Michael Smith
    'Patient Checked In',
    'Patient John Doe has checked in for their appointment and is waiting in the lobby.',
    'patient',
    'low',
    false,
    ARRAY['in_app'],
    '{"patient_id": "11150001-1111-1111-1111-111111111111", "patient_name": "John Doe", "appointment_id": "11170001-1111-1111-1111-111111111111", "check_in_time": "2025-01-07T14:00:00Z"}',
    NULL,
    NOW(),
    NOW() - INTERVAL '5 minutes'
),

-- Notifications for Dr. Jennifer Jones
(
    '11190009-1111-1111-1111-111111111111',
    '11130002-1111-1111-1111-111111111111', -- Dr. Jennifer Jones
    'Routine Checkup Scheduled',
    'A routine checkup has been scheduled for next week with patient Robert Wilson.',
    'appointment',
    'low',
    true,
    ARRAY['email', 'in_app'],
    '{"appointment_id": "11170005-1111-1111-1111-111111111111", "patient_name": "Robert Wilson", "appointment_time": "2025-01-14T09:00:00Z", "action_required": false}',
    NULL,
    NOW() - INTERVAL '1 day',
    NOW() - INTERVAL '1 day'
),

-- Queue status update
(
    '11190010-1111-1111-1111-111111111111',
    '11130002-1111-1111-1111-111111111111', -- Dr. Jennifer Jones
    'Queue Status Update',
    'Your current queue has 3 patients waiting. Average wait time is 15 minutes.',
    'alert',
    'low',
    false,
    ARRAY['in_app'],
    '{"queue_length": 3, "average_wait_time": 15, "next_patient": "Alice Johnson", "action_required": false}',
    NULL,
    NOW(),
    NOW() - INTERVAL '20 minutes'
),

-- System notification for all users
(
    '11190011-1111-1111-1111-111111111111',
    '11120001-1111-1111-1111-111111111111', -- Dr. Sarah Johnson
    'New Feature Available',
    'The new patient portal feature is now available. Patients can now book appointments online.',
    'system',
    'low',
    true,
    ARRAY['email', 'in_app'],
    '{"feature_name": "Patient Portal", "release_date": "2025-01-07", "documentation_url": "/docs/patient-portal"}',
    NULL,
    NOW() - INTERVAL '2 days',
    NOW() - INTERVAL '2 days'
),

-- Lab results for Dr. Michael Smith
(
    '11190012-1111-1111-1111-111111111111',
    '11130001-1111-1111-1111-111111111111', -- Dr. Michael Smith
    'Lab Results Ready',
    'X-ray results for patient Jane Smith are ready for review.',
    'lab',
    'medium',
    false,
    ARRAY['in_app'],
    '{"patient_id": "11150002-1111-1111-1111-111111111111", "patient_name": "Jane Smith", "test_type": "x_ray", "results_status": "abnormal", "action_required": true}',
    NULL,
    NOW(),
    NOW() - INTERVAL '45 minutes'
),

-- Appointment cancellation
(
    '11190013-1111-1111-1111-111111111111',
    '11130003-1111-1111-1111-111111111111', -- Dr. Robert Brown
    'Appointment Cancelled',
    'Patient John Doe has cancelled their appointment scheduled for tomorrow.',
    'appointment',
    'low',
    false,
    ARRAY['in_app'],
    '{"appointment_id": "11170003-1111-1111-1111-111111111111", "patient_name": "John Doe", "cancellation_reason": "Patient request", "original_time": "2025-01-08T11:00:00Z"}',
    NULL,
    NOW(),
    NOW() - INTERVAL '1 hour'
),

-- Queue position update
(
    '11190014-1111-1111-1111-111111111111',
    '11120001-1111-1111-1111-111111111111', -- Dr. Sarah Johnson
    'Queue Position Updated',
    'Patient Alice Johnson has moved to position 2 in the queue.',
    'alert',
    'low',
    false,
    ARRAY['in_app'],
    '{"queue_id": "11180002-1111-1111-1111-111111111111", "appointment_id": "11170002-1111-1111-1111-111111111111", "patient_name": "Alice Johnson", "new_position": 2, "previous_position": 3}',
    NULL,
    NOW(),
    NOW() - INTERVAL '10 minutes'
),

-- System backup notification
(
    '11190015-1111-1111-1111-111111111111',
    '11120001-1111-1111-1111-111111111111', -- Dr. Sarah Johnson
    'Backup Completed',
    'Daily system backup has been completed successfully. All data is secure.',
    'system',
    'low',
    true,
    ARRAY['in_app'],
    '{"backup_type": "daily", "backup_size": "2.5GB", "backup_duration": "15 minutes", "status": "success"}',
    NULL,
    NOW() - INTERVAL '6 hours',
    NOW() - INTERVAL '6 hours'
),

-- ===== NOTIFICATIONS FOR SYSTEM ADMINISTRATOR =====
(
    '11190016-1111-1111-1111-111111111111',
    '11120002-1111-1111-1111-111111111111', -- System Administrator
    'System Performance Alert',
    'Database response time has increased to 2.5 seconds. Consid
    er optimization.',
    'system',
    'high',
    false,
    ARRAY['email', 'in_app'],
    '{"metric": "database_response_time", "value": "2.5s", "threshold": "2.0s", "action_required": true}',
    NULL,
    NOW(),
    NOW() - INTERVAL '1 hour'
),

(
    '11190017-1111-1111-1111-111111111111',
    '11120002-1111-1111-1111-111111111111', -- System Administrator
    'Security Audit Required',
    'Monthly security audit is due. Please schedule and complete the audit.',
    'system',
    'medium',
    false,
    ARRAY['email', 'in_app'],
    '{"audit_type": "security", "due_date": "2025-01-15", "action_required": true}',
    NULL,
    NOW(),
    NOW() - INTERVAL '2 hours'
),

-- ===== NOTIFICATIONS FOR NURSES =====
(
    '11190018-1111-1111-1111-111111111111',
    '11140001-1111-1111-1111-111111111111', -- Nurse Emily Wilson
    'Patient Vitals Check Required',
    'Patient John Doe needs vital signs checked. Please visit room 101.',
    'patient',
    'medium',
    false,
    ARRAY['in_app'],
    '{"patient_id": "11150001-1111-1111-1111-111111111111", "patient_name": "John Doe", "room": "101", "vitals_type": "blood_pressure", "action_required": true}',
    NULL,
    NOW(),
    NOW() - INTERVAL '20 minutes'
),

(
    '11190019-1111-1111-1111-111111111111',
    '11140001-1111-1111-1111-111111111111', -- Nurse Emily Wilson
    'Medication Administration',
    'Patient Jane Smith is due for medication at 2:00 PM. Please prepare and administer.',
    'patient',
    'high',
    false,
    ARRAY['in_app'],
    '{"patient_id": "11150002-1111-1111-1111-111111111111", "patient_name": "Jane Smith", "medication": "Insulin", "dosage": "10 units", "due_time": "14:00", "action_required": true}',
    NULL,
    NOW(),
    NOW() - INTERVAL '10 minutes'
),

(
    '11190020-1111-1111-1111-111111111111',
    '11140002-1111-1111-1111-111111111111', -- Nurse Lisa Davis
    'Patient Discharge Ready',
    'Patient Bob Johnson is ready for discharge. Please complete discharge procedures.',
    'patient',
    'medium',
    false,
    ARRAY['in_app'],
    '{"patient_id": "11150003-1111-1111-1111-111111111111", "patient_name": "Bob Johnson", "discharge_time": "15:00", "action_required": true}',
    NULL,
    NOW(),
    NOW() - INTERVAL '30 minutes'
),

(
    '11190021-1111-1111-1111-111111111111',
    '11140002-1111-1111-1111-111111111111', -- Nurse Lisa Davis
    'Lab Sample Collection',
    'Blood sample collection required for patient John Doe. Please collect and send to lab.',
    'lab',
    'medium',
    false,
    ARRAY['in_app'],
    '{"patient_id": "11150001-1111-1111-1111-111111111111", "patient_name": "John Doe", "test_type": "blood_work", "collection_time": "immediate", "action_required": true}',
    NULL,
    NOW(),
    NOW() - INTERVAL '15 minutes'
),

-- ===== NOTIFICATIONS FOR CASHIERS =====
(
    '11190022-1111-1111-1111-111111111111',
    '11160001-1111-1111-1111-111111111111', -- Maria Rodriguez (Cashier)
    'Payment Pending',
    'Payment of $150.00 is pending for patient Jane Smith. Please process payment.',
    'patient',
    'medium',
    false,
    ARRAY['in_app'],
    '{"patient_id": "11150002-1111-1111-1111-111111111111", "patient_name": "Jane Smith", "amount": "$150.00", "payment_type": "consultation_fee", "action_required": true}',
    NULL,
    NOW(),
    NOW() - INTERVAL '45 minutes'
),

(
    '11190023-1111-1111-1111-111111111111',
    '11160001-1111-1111-1111-111111111111', -- Maria Rodriguez (Cashier)
    'Insurance Verification',
    'Insurance verification required for patient Bob Johnson before appointment.',
    'patient',
    'low',
    false,
    ARRAY['in_app'],
    '{"patient_id": "11150003-1111-1111-1111-111111111111", "patient_name": "Bob Johnson", "insurance_provider": "Blue Cross", "verification_status": "pending"}',
    NULL,
    NOW(),
    NOW() - INTERVAL '1 hour'
),

(
    '11190024-1111-1111-1111-111111111111',
    '11160002-1111-1111-1111-111111111111', -- Carlos Martinez (Cashier)
    'Daily Report Due',
    'Daily financial report is due at 5:00 PM. Please complete and submit.',
    'system',
    'medium',
    false,
    ARRAY['in_app'],
    '{"report_type": "daily_financial", "due_time": "17:00", "action_required": true}',
    NULL,
    NOW(),
    NOW() - INTERVAL '2 hours'
),

(
    '11190025-1111-1111-1111-111111111111',
    '11160002-1111-1111-1111-111111111111', -- Carlos Martinez (Cashier)
    'Refund Processed',
    'Refund of $75.00 has been processed for patient John Doe. Receipt generated.',
    'patient',
    'low',
    true,
    ARRAY['in_app'],
    '{"patient_id": "11150001-1111-1111-1111-111111111111", "patient_name": "John Doe", "refund_amount": "$75.00", "reason": "cancelled_appointment"}',
    NULL,
    NOW() - INTERVAL '3 hours',
    NOW() - INTERVAL '3 hours'
),

-- ===== NOTIFICATIONS FOR PATIENTS =====
(
    '11190026-1111-1111-1111-111111111111',
    '11150001-1111-1111-1111-111111111111', -- John Doe (Patient)
    'Appointment Reminder',
    'Your appointment with Dr. Michael Smith is tomorrow at 10:00 AM. Please arrive 15 minutes early.',
    'appointment',
    'medium',
    false,
    ARRAY['email', 'in_app'],
    '{"appointment_id": "11170001-1111-1111-1111-111111111111", "doctor_name": "Dr. Michael Smith", "appointment_time": "2025-01-08T10:00:00Z", "action_required": false}',
    NULL,
    NOW(),
    NOW() - INTERVAL '1 hour'
),

(
    '11190027-1111-1111-1111-111111111111',
    '11150001-1111-1111-1111-111111111111', -- John Doe (Patient)
    'Lab Results Available',
    'Your blood test results are now available in your patient portal.',
    'lab',
    'medium',
    false,
    ARRAY['email', 'in_app'],
    '{"test_type": "blood_test", "results_status": "normal", "view_portal": true}',
    NULL,
    NOW(),
    NOW() - INTERVAL '30 minutes'
),

(
    '11190028-1111-1111-1111-111111111111',
    '11150002-1111-1111-1111-111111111111', -- Jane Smith (Patient)
    'Prescription Ready',
    'Your prescription for medication is ready for pickup at the pharmacy.',
    'patient',
    'medium',
    false,
    ARRAY['email', 'in_app'],
    '{"medication": "Insulin", "pharmacy": "Main Street Pharmacy", "pickup_time": "anytime", "action_required": true}',
    NULL,
    NOW(),
    NOW() - INTERVAL '2 hours'
),

(
    '11190029-1111-1111-1111-111111111111',
    '11150002-1111-1111-1111-111111111111', -- Jane Smith (Patient)
    'Follow-up Appointment',
    'Your follow-up appointment with Dr. Jennifer Jones is scheduled for next week.',
    'appointment',
    'low',
    false,
    ARRAY['email', 'in_app'],
    '{"appointment_id": "11170002-1111-1111-1111-111111111111", "doctor_name": "Dr. Jennifer Jones", "appointment_time": "2025-01-14T14:30:00Z"}',
    NULL,
    NOW(),
    NOW() - INTERVAL '4 hours'
),

(
    '11190030-1111-1111-1111-111111111111',
    '11150003-1111-1111-1111-111111111111', -- Bob Johnson (Patient)
    'Insurance Update Required',
    'Please update your insurance information in your patient profile.',
    'patient',
    'low',
    false,
    ARRAY['email', 'in_app'],
    '{"insurance_provider": "Blue Cross", "update_required": "policy_number", "action_required": true}',
    NULL,
    NOW(),
    NOW() - INTERVAL '6 hours'
),

(
    '11190031-1111-1111-1111-111111111111',
    '11150003-1111-1111-1111-111111111111', -- Bob Johnson (Patient)
    'Payment Confirmation',
    'Your payment of $200.00 has been processed successfully. Thank you!',
    'patient',
    'low',
    true,
    ARRAY['email', 'in_app'],
    '{"amount": "$200.00", "payment_method": "credit_card", "transaction_id": "TXN123456"}',
    NULL,
    NOW() - INTERVAL '1 day',
    NOW() - INTERVAL '1 day'
),

-- ===== ADDITIONAL DOCTOR NOTIFICATIONS =====
(
    '11190032-1111-1111-1111-111111111111',
    '11130003-1111-1111-1111-111111111111', -- Dr. Robert Brown
    'Patient Consultation Request',
    'Patient Bob Johnson has requested a consultation. Please review and respond.',
    'appointment',
    'medium',
    false,
    ARRAY['in_app'],
    '{"patient_id": "11150003-1111-1111-1111-111111111111", "patient_name": "Bob Johnson", "request_type": "consultation", "action_required": true}',
    NULL,
    NOW(),
    NOW() - INTERVAL '25 minutes'
),

(
    '11190033-1111-1111-1111-111111111111',
    '11130003-1111-1111-1111-111111111111', -- Dr. Robert Brown
    'Medical Records Update',
    'Patient medical records have been updated. Please review the changes.',
    'patient',
    'low',
    false,
    ARRAY['in_app'],
    '{"patient_id": "11150001-1111-1111-1111-111111111111", "patient_name": "John Doe", "update_type": "vital_signs", "updated_by": "Nurse Emily Wilson"}',
    NULL,
    NOW(),
    NOW() - INTERVAL '40 minutes'
);

-- No sequence update needed as we're using gen_random_uuid() for ID generation
