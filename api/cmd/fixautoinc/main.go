package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type colInfo struct {
	Name     string
	Ordinal  int
	Nullable string
}

func main() {
	dsn := strings.TrimSpace(os.Getenv("GOOSE_DBSTRING"))
	if dsn == "" {
		log.Fatal("GOOSE_DBSTRING is missing")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	tables := []string{
		"course_details",
		"course_chapters",
		"course_lessons",
		"lesson_resources",
		"course_general_settings",
		"course_instructors",
		"course_quizzes",
		"quiz_questions",
		"course_assignments",
	}

	for _, t := range tables {
		if err := fixTable(db, t); err != nil {
			log.Fatalf("fix %s: %v", t, err)
		}
		fmt.Printf("fixed: %s\n", t)
	}
}

func fixTable(db *sql.DB, table string) error {
	ai, err := tableHasAutoIncrementID(db, table)
	if err != nil {
		return err
	}
	if ai {
		return nil
	}

	cols, err := getColumns(db, table)
	if err != nil {
		return err
	}
	if len(cols) == 0 {
		return fmt.Errorf("no columns found for %s", table)
	}

	createSQL, err := showCreateTable(db, table)
	if err != nil {
		return err
	}

	newTable := table + "__ai_fix"
	bakTable := table + "__ai_bak"

	if err := execIgnoreAlreadyExists(db, fmt.Sprintf("DROP TABLE IF EXISTS `%s`", newTable)); err != nil {
		return err
	}

	newCreate := rewriteCreateTableAddAutoIncrement(createSQL, table, newTable)
	if _, err := db.Exec(newCreate); err != nil {
		return fmt.Errorf("create %s: %w", newTable, err)
	}

	colList := make([]string, 0, len(cols))
	for _, c := range cols {
		colList = append(colList, fmt.Sprintf("`%s`", c.Name))
	}

	ins := fmt.Sprintf(
		"INSERT INTO `%s` (%s) SELECT %s FROM `%s`",
		newTable,
		strings.Join(colList, ", "),
		strings.Join(colList, ", "),
		table,
	)
	if _, err := db.Exec(ins); err != nil {
		return fmt.Errorf("copy data: %w", err)
	}

	_ = execIgnoreAlreadyExists(db, fmt.Sprintf("DROP TABLE IF EXISTS `%s`", bakTable))

	rename := fmt.Sprintf(
		"RENAME TABLE `%s` TO `%s`, `%s` TO `%s`",
		table, bakTable,
		newTable, table,
	)
	if _, err := db.Exec(rename); err != nil {
		return fmt.Errorf("rename swap: %w", err)
	}

	return nil
}

func tableHasAutoIncrementID(db *sql.DB, table string) (bool, error) {
	const q = `
SELECT EXTRA
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = ?
  AND COLUMN_NAME = 'id'
LIMIT 1;
`
	var extra sql.NullString
	if err := db.QueryRow(q, table).Scan(&extra); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s has no id column", table)
		}
		return false, err
	}
	return strings.Contains(strings.ToLower(extra.String), "auto_increment"), nil
}

func getColumns(db *sql.DB, table string) ([]colInfo, error) {
	const q = `
SELECT COLUMN_NAME, ORDINAL_POSITION, IS_NULLABLE
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = ?
ORDER BY ORDINAL_POSITION ASC;
`
	rows, err := db.Query(q, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cols []colInfo
	for rows.Next() {
		var c colInfo
		if err := rows.Scan(&c.Name, &c.Ordinal, &c.Nullable); err != nil {
			return nil, err
		}
		cols = append(cols, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Defensive: ensure stable ordering
	sort.Slice(cols, func(i, j int) bool { return cols[i].Ordinal < cols[j].Ordinal })
	return cols, nil
}

func showCreateTable(db *sql.DB, table string) (string, error) {
	row := db.QueryRow(fmt.Sprintf("SHOW CREATE TABLE `%s`", table))
	var tbl, createSQL string
	if err := row.Scan(&tbl, &createSQL); err != nil {
		return "", err
	}
	return createSQL, nil
}

func rewriteCreateTableAddAutoIncrement(createSQL, oldName, newName string) string {
	// Replace table name first.
	out := strings.Replace(createSQL, fmt.Sprintf("CREATE TABLE `%s`", oldName), fmt.Sprintf("CREATE TABLE `%s`", newName), 1)

	// Ensure the `id` column includes AUTO_INCREMENT (case-insensitive).
	// Example: `id` int unsigned NOT NULL,
	idLine := regexp.MustCompile("(?im)^\\s*`id`\\s+([^,]+),\\s*$")
	m := idLine.FindStringSubmatch(out)
	if len(m) == 2 {
		def := m[1]
		if !regexp.MustCompile("(?i)auto_increment").MatchString(def) {
			def = strings.TrimSpace(def) + " AUTO_INCREMENT"
			out = idLine.ReplaceAllString(out, fmt.Sprintf("  `id` %s,", def))
		}
	}

	return out
}

func execIgnoreAlreadyExists(db *sql.DB, stmt string) error {
	_, err := db.Exec(stmt)
	return err
}
