-- Seed notifications for existing users
-- This script creates sample notifications for testing

-- Clear existing notifications first
DELETE FROM notifications;

-- Get some existing user IDs and create notifications
INSERT INTO notifications (user_id, title, message, type, priority, is_read, channels, data, created_at) VALUES
-- Notifications for first user (assuming user exists)
((SELECT id FROM users LIMIT 1 OFFSET 0), 'New Appointment Request', 'John Doe requested an appointment for tomorrow at 2:00 PM', 'appointment', 'high', false, '{"in_app"}', '{"appointment_id": "apt-001", "patient_name": "John Doe", "requested_time": "2024-01-16T14:00:00Z"}', NOW() - INTERVAL '2 minutes'),

((SELECT id FROM users LIMIT 1 OFFSET 0), 'Patient Check-in', 'Jane Smith has checked in for her 10:30 AM appointment', 'patient', 'medium', false, '{"in_app"}', '{"patient_id": "pat-002", "patient_name": "Jane Smith", "appointment_time": "2024-01-15T10:30:00Z"}', NOW() - INTERVAL '15 minutes'),

((SELECT id FROM users LIMIT 1 OFFSET 0), 'Lab Results Ready', 'Blood test results for Mike Johnson are now available', 'lab', 'low', true, '{"in_app"}', '{"patient_id": "pat-003", "patient_name": "Mike Johnson", "lab_type": "blood_test", "results_url": "/lab-results/123"}', NOW() - INTERVAL '1 hour'),

((SELECT id FROM users LIMIT 1 OFFSET 0), 'Emergency Alert', 'Emergency patient arrived - Dr. Sarah Johnson needed in ER', 'emergency', 'critical', true, '{"in_app", "email"}', '{"patient_id": "pat-004", "patient_name": "Emergency Patient", "room": "ER-1", "severity": "critical"}', NOW() - INTERVAL '2 hours'),

((SELECT id FROM users LIMIT 1 OFFSET 0), 'Schedule Update', 'Your 3:00 PM appointment has been rescheduled to 4:00 PM', 'schedule', 'medium', true, '{"in_app"}', '{"appointment_id": "apt-005", "old_time": "2024-01-15T15:00:00Z", "new_time": "2024-01-15T16:00:00Z"}', NOW() - INTERVAL '3 hours'),

-- Notifications for second user (if exists)
((SELECT id FROM users LIMIT 1 OFFSET 1), 'System Maintenance', 'Scheduled maintenance will occur tonight from 11 PM to 1 AM', 'system', 'low', false, '{"in_app"}', '{"maintenance_type": "database", "start_time": "2024-01-15T23:00:00Z", "end_time": "2024-01-16T01:00:00Z"}', NOW() - INTERVAL '30 minutes'),

((SELECT id FROM users LIMIT 1 OFFSET 1), 'New Message', 'Dr. Sarah Johnson sent you a message regarding patient care protocol', 'message', 'medium', true, '{"in_app"}', '{"sender_id": "user-001", "sender_name": "Dr. Sarah Johnson", "message_type": "protocol_update"}', NOW() - INTERVAL '1 hour'),

-- Notifications for third user (if exists)
((SELECT id FROM users LIMIT 1 OFFSET 2), 'Monthly Report Available', 'December patient statistics and department performance report is ready', 'system', 'low', true, '{"in_app"}', '{"report_type": "monthly_stats", "month": "December", "year": "2024", "report_url": "/reports/monthly-2024-12"}', NOW() - INTERVAL '4 hours'),

((SELECT id FROM users LIMIT 1 OFFSET 2), 'Appointment Cancelled', 'Sarah Wilson has cancelled her appointment scheduled for tomorrow at 2 PM', 'appointment', 'low', true, '{"in_app"}', '{"appointment_id": "apt-006", "patient_name": "Sarah Wilson", "cancelled_time": "2024-01-16T14:00:00Z", "reason": "patient_request"}', NOW() - INTERVAL '2 hours'),

-- Additional notifications for variety
((SELECT id FROM users LIMIT 1 OFFSET 0), 'Medication Reminder', 'Patient needs to take medication at 8:00 AM', 'patient', 'medium', false, '{"in_app"}', '{"patient_id": "pat-007", "medication": "Metformin", "dosage": "500mg", "time": "08:00"}', NOW() - INTERVAL '5 minutes'),

((SELECT id FROM users LIMIT 1 OFFSET 1), 'Insurance Verification', 'Insurance verification required for patient John Smith', 'patient', 'high', false, '{"in_app"}', '{"patient_id": "pat-008", "patient_name": "John Smith", "insurance_provider": "Blue Cross", "verification_deadline": "2024-01-16T17:00:00Z"}', NOW() - INTERVAL '10 minutes'),

((SELECT id FROM users LIMIT 1 OFFSET 2), 'Equipment Maintenance', 'X-ray machine in Room 3 requires scheduled maintenance', 'system', 'medium', true, '{"in_app"}', '{"equipment_id": "xray-003", "room": "Room 3", "maintenance_type": "scheduled", "scheduled_date": "2024-01-20T09:00:00Z"}', NOW() - INTERVAL '1 day'),

-- More recent notifications
((SELECT id FROM users LIMIT 1 OFFSET 0), 'New Patient Registration', 'Emily Davis has completed her registration and requires approval', 'patient', 'medium', false, '{"in_app"}', '{"patient_id": "pat-009", "patient_name": "Emily Davis", "registration_date": "2024-01-15T08:30:00Z", "approval_required": true}', NOW() - INTERVAL '30 seconds'),

((SELECT id FROM users LIMIT 1 OFFSET 1), 'Critical Patient Alert', 'Patient Michael Johnson in Room 304 requires immediate attention', 'alert', 'critical', false, '{"in_app", "email", "sms"}', '{"patient_id": "pat-010", "patient_name": "Michael Johnson", "room": "304", "alert_type": "medical_emergency", "severity": "critical"}', NOW() - INTERVAL '1 minute'),

((SELECT id FROM users LIMIT 1 OFFSET 2), 'Appointment Reminder', 'You have an appointment with John Smith at 10:00 AM today', 'appointment', 'high', false, '{"in_app"}', '{"appointment_id": "apt-007", "patient_name": "John Smith", "appointment_time": "2024-01-15T10:00:00Z", "room": "Room 2"}', NOW() - INTERVAL '2 minutes');
