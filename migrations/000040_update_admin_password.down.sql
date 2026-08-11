-- Restore original admin password 'admin123'
UPDATE users
SET password_hash = '$2a$10$bE87GqCBB5/QMO.grdBik.lS4jna0wVGid3B2qdiNc1o6V4Ljl4Ay'
WHERE username = 'admin';
