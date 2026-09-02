package coursereview

import (
	"errors"
	"strings"

	mysqlDriver "github.com/go-sql-driver/mysql"
)

func isMissingTable(err error) bool {
	if err == nil {
		return false
	}

	var mysqlErr *mysqlDriver.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1146 {
		return true
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "doesn't exist") || strings.Contains(msg, "does not exist")
}
