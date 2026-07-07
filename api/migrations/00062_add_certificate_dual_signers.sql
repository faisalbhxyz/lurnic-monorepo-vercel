-- +goose Up
ALTER TABLE course_certificate_settings
ADD COLUMN dual_signers_enabled TINYINT(1) NOT NULL DEFAULT 0;

ALTER TABLE course_certificate_settings
ADD COLUMN signer2_name VARCHAR(255) NULL;

ALTER TABLE course_certificate_settings
ADD COLUMN signer2_role VARCHAR(255) NULL;

ALTER TABLE course_certificate_settings
ADD COLUMN signer2_org VARCHAR(255) NULL;

ALTER TABLE student_certificates
ADD COLUMN dual_signers_enabled TINYINT(1) NOT NULL DEFAULT 0;

ALTER TABLE student_certificates
ADD COLUMN signer2_name VARCHAR(255) NULL;

ALTER TABLE student_certificates
ADD COLUMN signer2_role VARCHAR(255) NULL;

ALTER TABLE student_certificates
ADD COLUMN signer2_org VARCHAR(255) NULL;

-- +goose Down
ALTER TABLE student_certificates DROP COLUMN signer2_org;
ALTER TABLE student_certificates DROP COLUMN signer2_role;
ALTER TABLE student_certificates DROP COLUMN signer2_name;
ALTER TABLE student_certificates DROP COLUMN dual_signers_enabled;

ALTER TABLE course_certificate_settings DROP COLUMN signer2_org;
ALTER TABLE course_certificate_settings DROP COLUMN signer2_role;
ALTER TABLE course_certificate_settings DROP COLUMN signer2_name;
ALTER TABLE course_certificate_settings DROP COLUMN dual_signers_enabled;
