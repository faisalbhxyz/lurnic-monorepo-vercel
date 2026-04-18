package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dumpPath := "lurnic-Database backup.sql"
	if len(os.Args) > 1 {
		dumpPath = os.Args[1]
	}

	raw := strings.TrimSpace(os.Getenv("GOOSE_DBSTRING"))
	if raw == "" {
		log.Fatal("GOOSE_DBSTRING is required")
	}

	user, pass, host, port, dbName, q := parseToParts(raw)
	if dbName == "" {
		log.Fatal("database name missing in GOOSE_DBSTRING (path after host)")
	}

	tlsQ := ensureTLS(q)
	adminDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/?%s", user, pass, host, port, tlsQ)
	dumpDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", user, pass, host, port, dbName, mergeQuery(tlsQ, "multiStatements=true"))

	sqlBytes, err := os.ReadFile(dumpPath)
	if err != nil {
		log.Fatal(err)
	}

	adminDB, err := sql.Open("mysql", adminDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer adminDB.Close()
	if err := adminDB.Ping(); err != nil {
		log.Fatal("admin ping:", err)
	}

	log.Printf("recreating database %q", dbName)
	ident := strings.ReplaceAll(dbName, "`", "")
	if _, err := adminDB.Exec("DROP DATABASE IF EXISTS `" + ident + "`"); err != nil {
		log.Fatal("drop database:", err)
	}
	if _, err := adminDB.Exec("CREATE DATABASE `" + ident + "` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"); err != nil {
		log.Fatal("create database:", err)
	}
	adminDB.Close()

	dumpDB, err := sql.Open("mysql", dumpDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer dumpDB.Close()
	if err := dumpDB.Ping(); err != nil {
		log.Fatal("dump ping:", err)
	}

	sqlStr := string(sqlBytes)
	// phpMyAdmin emits ALTER TABLE ... MODIFY ... AUTO_INCREMENT which TiDB rejects (Error 8200).
	if i := strings.Index(sqlStr, "-- AUTO_INCREMENT for dumped tables"); i >= 0 {
		if j := strings.Index(sqlStr, "-- Constraints for dumped tables"); j > i {
			sqlStr = sqlStr[:i] + "-- (skipped: AUTO_INCREMENT MODIFY blocks — TiDB incompatible)\n\n" + sqlStr[j:]
		}
	}

	log.Println("importing dump...")
	if _, err := dumpDB.Exec(sqlStr); err != nil {
		log.Fatal("import:", err)
	}
	log.Println("import finished OK")
}

func mergeQuery(a, b string) string {
	a = strings.TrimPrefix(a, "?")
	if a == "" {
		return b
	}
	if strings.Contains(a, b) {
		return a
	}
	return a + "&" + b
}

func ensureTLS(q string) string {
	if strings.Contains(q, "tls=") {
		return q
	}
	if q == "" {
		return "tls=true"
	}
	return q + "&tls=true"
}

// parseToParts accepts either mysql:// URL or Go driver DSN user:pass@tcp(host:port)/db?params
func parseToParts(raw string) (user, pass, host, port, dbName, query string) {
	s := strings.TrimSpace(raw)
	if strings.HasPrefix(s, "mysql://") {
		u, err := url.Parse(s)
		if err != nil {
			log.Fatal("parse URL:", err)
		}
		pass, _ = u.User.Password()
		user = u.User.Username()
		host = u.Hostname()
		port = u.Port()
		if port == "" {
			port = "3306"
		}
		dbName = strings.TrimPrefix(u.Path, "/")
		query = u.RawQuery
		return user, pass, host, port, dbName, query
	}

	// user:pass@tcp(host:port)/dbname?query
	at := strings.LastIndex(s, "@tcp(")
	if at < 0 {
		log.Fatal("invalid GOOSE_DBSTRING: expected mysql:// or @tcp(")
	}
	userinfo := s[:at]
	c := strings.Index(userinfo, ":")
	if c < 0 {
		log.Fatal("invalid GOOSE_DBSTRING: missing user:password")
	}
	user = userinfo[:c]
	pass = userinfo[c+1:]

	rest := s[at+len("@tcp("):]
	end := strings.Index(rest, ")")
	if end < 0 {
		log.Fatal("invalid GOOSE_DBSTRING: missing ) after host:port")
	}
	hp := strings.Split(rest[:end], ":")
	if len(hp) != 2 {
		log.Fatal("invalid GOOSE_DBSTRING: host:port")
	}
	host, port = hp[0], hp[1]

	after := rest[end+1:]
	if !strings.HasPrefix(after, ")/") {
		log.Fatal("invalid GOOSE_DBSTRING: expected )/ after host:port")
	}
	after = after[len(")/"):]
	if i := strings.Index(after, "?"); i >= 0 {
		dbName = after[:i]
		query = after[i+1:]
	} else {
		dbName = after
	}
	return user, pass, host, port, dbName, query
}
