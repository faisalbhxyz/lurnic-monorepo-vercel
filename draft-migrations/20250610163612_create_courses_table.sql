-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    courses (
        id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        description TEXT NULL,
        visibility ENUM ('public', 'private', 'protected') NOT NULL DEFAULT 'public',
        is_scheduled BOOLEAN NOT NULL DEFAULT FALSE,
        schedule_date DATE NULL,
        schedule_time TIME NULL,
        show_comming_soon BOOLEAN NOT NULL DEFAULT FALSE,
        featured_image VARCHAR(255) NULL,
        intro_video VARCHAR(255) NULL,
        pricing_model ENUM ('free', 'paid') NOT NULL DEFAULT 'free',
        tags JSON NULL,
        author_id INT UNSIGNED NOT NULL,
        difficulty_level ENUM ('all', 'beginner', 'intermediate', 'expert') NOT NULL DEFAULT 'all',
        is_public_course BOOLEAN NOT NULL DEFAULT FALSE,
        maximum_student INT NOT NULL DEFAULT 0,
        tenant_id INT UNSIGNED NOT NULL,
        FOREIGN KEY (author_id) REFERENCES users (id),
        FOREIGN KEY (tenant_id) REFERENCES tenants (id)
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS courses;

-- +goose StatementEnd