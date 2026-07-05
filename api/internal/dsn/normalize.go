package dsn

import (
	"fmt"
	"net/url"
	"strings"
)

// Normalize converts mysql:// URLs or Go driver DSNs into a normalized Go MySQL DSN
// with TLS/timeouts suitable for TiDB Cloud and managed MySQL from Docker.
func Normalize(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", fmt.Errorf("GOOSE_DBSTRING is empty")
	}

	var goDSN string
	switch {
	case strings.HasPrefix(s, "mysql://"):
		var err error
		goDSN, err = fromMySQLURL(s)
		if err != nil {
			return "", err
		}
	case strings.Contains(s, "@tcp("):
		goDSN = s
	default:
		return "", fmt.Errorf("GOOSE_DBSTRING must be mysql:// URL or user:pass@tcp(host:port)/db form")
	}

	return applyDefaults(goDSN), nil
}

func fromMySQLURL(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", fmt.Errorf("parse mysql URL: %w", err)
	}
	pass, _ := u.User.Password()
	user := u.User.Username()
	if user == "" {
		return "", fmt.Errorf("mysql URL missing username")
	}
	host := u.Hostname()
	if host == "" {
		return "", fmt.Errorf("mysql URL missing host")
	}
	port := u.Port()
	if port == "" {
		port = "3306"
	}
	dbName := strings.TrimPrefix(u.Path, "/")
	if dbName == "" {
		return "", fmt.Errorf("mysql URL missing database name")
	}

	base := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, dbName)
	if u.RawQuery == "" {
		return base, nil
	}
	return base + "?" + u.RawQuery, nil
}

func applyDefaults(goDSN string) string {
	base, q := splitQuery(goDSN)
	q = ensureQueryKV(q, "tls", "skip-verify")
	q = ensureQueryKV(q, "timeout", "5s")
	q = ensureQueryKV(q, "readTimeout", "15s")
	q = ensureQueryKV(q, "writeTimeout", "15s")
	if q == "" {
		return base
	}
	return base + "?" + q
}

func splitQuery(s string) (base, query string) {
	if i := strings.Index(s, "?"); i >= 0 {
		return s[:i], s[i+1:]
	}
	return s, ""
}

func ensureQueryKV(query, key, value string) string {
	if query == "" {
		return key + "=" + value
	}
	parts := strings.Split(query, "&")
	for _, p := range parts {
		if p == key || strings.HasPrefix(p, key+"=") {
			return query
		}
	}
	return query + "&" + key + "=" + value
}
