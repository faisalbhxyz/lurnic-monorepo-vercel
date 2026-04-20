-- +goose Up
-- Speed up cron queries: WHERE is_scheduled = true (course publish)
SET @__dl_idx := (
  SELECT COUNT(1)
  FROM information_schema.statistics
  WHERE table_schema = DATABASE()
    AND table_name = 'course_details'
    AND index_name = 'idx_course_details_is_scheduled'
);
SET @__dl_sql := IF(
  @__dl_idx = 0,
  'CREATE INDEX idx_course_details_is_scheduled ON course_details (is_scheduled)',
  'SELECT 1'
);
PREPARE __dl_stmt FROM @__dl_sql;
EXECUTE __dl_stmt;
DEALLOCATE PREPARE __dl_stmt;

-- Speed up cron: WHERE is_scheduled = true AND is_published = false (lesson publish)
SET @__dl_idx := (
  SELECT COUNT(1)
  FROM information_schema.statistics
  WHERE table_schema = DATABASE()
    AND table_name = 'course_lessons'
    AND index_name = 'idx_course_lessons_scheduled_published'
);
SET @__dl_sql := IF(
  @__dl_idx = 0,
  'CREATE INDEX idx_course_lessons_scheduled_published ON course_lessons (is_scheduled, is_published)',
  'SELECT 1'
);
PREPARE __dl_stmt FROM @__dl_sql;
EXECUTE __dl_stmt;
DEALLOCATE PREPARE __dl_stmt;

-- cleanup vars (avoid leaking between statements in same session)
SET @__dl_idx := NULL;
SET @__dl_sql := NULL;
-- +goose Down
SET @__dl_idx := (
  SELECT COUNT(1)
  FROM information_schema.statistics
  WHERE table_schema = DATABASE()
    AND table_name = 'course_lessons'
    AND index_name = 'idx_course_lessons_scheduled_published'
);
SET @__dl_sql := IF(
  @__dl_idx = 1,
  'DROP INDEX idx_course_lessons_scheduled_published ON course_lessons',
  'SELECT 1'
);
PREPARE __dl_stmt FROM @__dl_sql;
EXECUTE __dl_stmt;
DEALLOCATE PREPARE __dl_stmt;

SET @__dl_idx := (
  SELECT COUNT(1)
  FROM information_schema.statistics
  WHERE table_schema = DATABASE()
    AND table_name = 'course_details'
    AND index_name = 'idx_course_details_is_scheduled'
);
SET @__dl_sql := IF(
  @__dl_idx = 1,
  'DROP INDEX idx_course_details_is_scheduled ON course_details',
  'SELECT 1'
);
PREPARE __dl_stmt FROM @__dl_sql;
EXECUTE __dl_stmt;
DEALLOCATE PREPARE __dl_stmt;

SET @__dl_idx := NULL;
SET @__dl_sql := NULL;
