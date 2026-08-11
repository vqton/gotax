-- Update admin password from 'admin123' to 'Admin@123456!'
-- Consistent with cmd/seed/main.go seedAdminPassword constant
UPDATE users
SET password_hash = '$2a$10$8APuAy0MlmqaacNduu1B3eDpIxuhD2lKHRWZC43T.StfzFoA5/ljm'
WHERE username = 'admin';
