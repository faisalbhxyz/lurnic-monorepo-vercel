-- +goose Up
SET @__dl_col := (
  SELECT COUNT(1)
  FROM information_schema.COLUMNS
  WHERE table_schema = DATABASE()
    AND table_name = 'quiz_questions'
    AND column_name = 'options'
);
SET @__dl_sql := IF(
  @__dl_col = 0,
  'ALTER TABLE quiz_questions ADD COLUMN options JSON NULL DEFAULT NULL',
  'SELECT 1'
);
PREPARE __dl_stmt FROM @__dl_sql;
EXECUTE __dl_stmt;
DEALLOCATE PREPARE __dl_stmt;
SET @__dl_col := NULL;
SET @__dl_sql := NULL;

-- +goose Down
SET @__dl_col := (
  SELECT COUNT(1)
  FROM information_schema.COLUMNS
  WHERE table_schema = DATABASE()
    AND table_name = 'quiz_questions'
    AND column_name = 'options'
);
SET @__dl_sql := IF(
  @__dl_col = 1,
  'ALTER TABLE quiz_questions DROP COLUMN options',
  'SELECT 1'
);
PREPARE __dl_stmt FROM @__dl_sql;
EXECUTE __dl_stmt;
DEALLOCATE PREPARE __dl_stmt;
SET @__dl_col := NULL;
SET @__dl_sql := NULL;
