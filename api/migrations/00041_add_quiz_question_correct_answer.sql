-- +goose Up
SET @__dl_col := (
  SELECT COUNT(1)
  FROM information_schema.COLUMNS
  WHERE table_schema = DATABASE()
    AND table_name = 'quiz_questions'
    AND column_name = 'correct_answer'
);
SET @__dl_sql := IF(
  @__dl_col = 0,
  'ALTER TABLE quiz_questions ADD COLUMN correct_answer JSON NULL DEFAULT NULL',
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
    AND column_name = 'correct_answer'
);
SET @__dl_sql := IF(
  @__dl_col = 1,
  'ALTER TABLE quiz_questions DROP COLUMN correct_answer',
  'SELECT 1'
);
PREPARE __dl_stmt FROM @__dl_sql;
EXECUTE __dl_stmt;
DEALLOCATE PREPARE __dl_stmt;
SET @__dl_col := NULL;
SET @__dl_sql := NULL;
