-- +goose Up
CREATE TABLE
    course_certificate_settings (
        id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        course_id INT UNSIGNED NOT NULL,
        is_enabled BOOLEAN NOT NULL DEFAULT FALSE,
        completion_percent TINYINT UNSIGNED NOT NULL DEFAULT 100,
        template_path VARCHAR(255) NOT NULL DEFAULT '/images/Certificat-14.jpg',
        title VARCHAR(255) NULL,
        subtitle_one VARCHAR(255) NULL,
        subtitle_two VARCHAR(255) NULL,
        owner_signature VARCHAR(500) NULL,
        instructor_signature VARCHAR(500) NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        FOREIGN KEY (course_id) REFERENCES course_details (id) ON DELETE CASCADE,
        UNIQUE KEY uq_course_certificate_settings_course_id (course_id)
    );

-- +goose StatementBegin
CREATE TABLE
    student_certificates (
        id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        tenant_id INT UNSIGNED NOT NULL,
        student_id INT UNSIGNED NOT NULL,
        course_id INT UNSIGNED NOT NULL,
        certificate_number VARCHAR(64) NOT NULL,
        student_name VARCHAR(255) NOT NULL,
        course_title VARCHAR(255) NOT NULL,
        progress_percent DECIMAL(5, 1) NOT NULL,
        template_path VARCHAR(255) NOT NULL,
        title VARCHAR(255) NULL,
        subtitle_one VARCHAR(255) NULL,
        subtitle_two VARCHAR(255) NULL,
        owner_signature VARCHAR(500) NULL,
        instructor_signature VARCHAR(500) NULL,
        issued_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE,
        FOREIGN KEY (student_id) REFERENCES students (id) ON DELETE CASCADE,
        FOREIGN KEY (course_id) REFERENCES course_details (id) ON DELETE CASCADE,
        UNIQUE KEY uq_student_certificates_number (certificate_number),
        UNIQUE KEY uq_student_certificates_student_course (tenant_id, student_id, course_id)
    );
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS student_certificates;

DROP TABLE IF EXISTS course_certificate_settings;
