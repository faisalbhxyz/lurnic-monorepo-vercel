-- +goose Up
-- +goose StatementBegin
CREATE TABLE student_daily_watch (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id INT UNSIGNED NOT NULL,
    student_id INT UNSIGNED NOT NULL,
    watch_date DATE NOT NULL,
    timezone VARCHAR(64) NOT NULL,
    video_seconds INT NOT NULL DEFAULT 0,
    live_class_seconds INT NOT NULL DEFAULT 0,
    quiz_seconds INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uq_daily_watch (tenant_id, student_id, watch_date),
    KEY idx_daily_watch_student (tenant_id, student_id, watch_date),
    KEY idx_daily_watch_updated (tenant_id, student_id, updated_at),
    CONSTRAINT fk_daily_watch_tenant FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE,
    CONSTRAINT fk_daily_watch_student FOREIGN KEY (student_id) REFERENCES students (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE student_watch_events (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id INT UNSIGNED NOT NULL,
    student_id INT UNSIGNED NOT NULL,
    course_id INT UNSIGNED NULL,
    lesson_id INT UNSIGNED NULL,
    source VARCHAR(32) NOT NULL,
    watched_seconds INT NOT NULL,
    client_event_id VARCHAR(64) NOT NULL,
    watched_at TIMESTAMP NOT NULL,
    watch_date DATE NOT NULL,
    timezone VARCHAR(64) NOT NULL,
    device_platform VARCHAR(16) NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_watch_event (tenant_id, student_id, client_event_id),
    KEY idx_watch_events_day (tenant_id, student_id, watch_date),
    KEY idx_watch_events_lesson (tenant_id, lesson_id),
    KEY idx_watch_events_client (client_event_id),
    CONSTRAINT fk_watch_events_tenant FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE,
    CONSTRAINT fk_watch_events_student FOREIGN KEY (student_id) REFERENCES students (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS student_watch_events;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS student_daily_watch;
-- +goose StatementEnd
