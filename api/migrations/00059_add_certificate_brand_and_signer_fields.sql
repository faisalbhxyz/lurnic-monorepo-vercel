-- +goose Up
ALTER TABLE course_certificate_settings
ADD COLUMN brand_logo VARCHAR(500) NULL;

ALTER TABLE course_certificate_settings
ADD COLUMN organization_name VARCHAR(255) NULL;

ALTER TABLE course_certificate_settings
ADD COLUMN signer_name VARCHAR(255) NULL;

ALTER TABLE course_certificate_settings
ADD COLUMN signer_role VARCHAR(255) NULL;

ALTER TABLE course_certificate_settings
ADD COLUMN signer_org VARCHAR(255) NULL;

ALTER TABLE student_certificates
ADD COLUMN brand_logo VARCHAR(500) NULL;

ALTER TABLE student_certificates
ADD COLUMN organization_name VARCHAR(255) NULL;

ALTER TABLE student_certificates
ADD COLUMN signer_name VARCHAR(255) NULL;

ALTER TABLE student_certificates
ADD COLUMN signer_role VARCHAR(255) NULL;

ALTER TABLE student_certificates
ADD COLUMN signer_org VARCHAR(255) NULL;

ALTER TABLE student_certificates
ADD COLUMN pricing_model ENUM ('free', 'paid') NOT NULL DEFAULT 'free';

-- +goose Down
ALTER TABLE student_certificates DROP COLUMN pricing_model;

ALTER TABLE student_certificates DROP COLUMN signer_org;

ALTER TABLE student_certificates DROP COLUMN signer_role;

ALTER TABLE student_certificates DROP COLUMN signer_name;

ALTER TABLE student_certificates DROP COLUMN organization_name;

ALTER TABLE student_certificates DROP COLUMN brand_logo;

ALTER TABLE course_certificate_settings DROP COLUMN signer_org;

ALTER TABLE course_certificate_settings DROP COLUMN signer_role;

ALTER TABLE course_certificate_settings DROP COLUMN signer_name;

ALTER TABLE course_certificate_settings DROP COLUMN organization_name;

ALTER TABLE course_certificate_settings DROP COLUMN brand_logo;
