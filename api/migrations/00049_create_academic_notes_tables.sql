-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    academic_note_classes (
        id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        tenant_id INT UNSIGNED NOT NULL,
        title VARCHAR(150) NOT NULL,
        slug VARCHAR(180) NOT NULL,
        icon_label VARCHAR(10) NULL,
        icon_color VARCHAR(20) NULL,
        position INT DEFAULT 0,
        is_published TINYINT(1) DEFAULT 1,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE
    );

-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE
    academic_note_subjects (
        id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        class_id BIGINT UNSIGNED NOT NULL,
        title VARCHAR(150) NOT NULL,
        slug VARCHAR(180) NOT NULL,
        position INT DEFAULT 0,
        is_published TINYINT(1) DEFAULT 1,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        FOREIGN KEY (class_id) REFERENCES academic_note_classes (id) ON DELETE CASCADE
    );

-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE
    academic_note_papers (
        id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        subject_id BIGINT UNSIGNED NOT NULL,
        title VARCHAR(150) NOT NULL,
        slug VARCHAR(180) NOT NULL,
        icon_label VARCHAR(10) NULL,
        icon_color VARCHAR(20) NULL,
        position INT DEFAULT 0,
        is_published TINYINT(1) DEFAULT 1,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        FOREIGN KEY (subject_id) REFERENCES academic_note_subjects (id) ON DELETE CASCADE
    );

-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE
    academic_notes (
        id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        paper_id BIGINT UNSIGNED NOT NULL,
        title VARCHAR(255) NOT NULL,
        subtitle VARCHAR(255) NULL,
        thumbnail TEXT NULL,
        pdf_url TEXT NOT NULL,
        pdf_file_name VARCHAR(255) NULL,
        position INT DEFAULT 0,
        is_published TINYINT(1) DEFAULT 1,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        FOREIGN KEY (paper_id) REFERENCES academic_note_papers (id) ON DELETE CASCADE
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS academic_notes;

-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS academic_note_papers;

-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS academic_note_subjects;

-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS academic_note_classes;

-- +goose StatementEnd
