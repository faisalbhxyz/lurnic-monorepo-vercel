package dsn

import "testing"

func TestNormalizeMySQLURL(t *testing.T) {
	got, err := Normalize("mysql://user:secret@tidb.example.com:4000/lurnic?charset=utf8mb4&parseTime=True")
	if err != nil {
		t.Fatal(err)
	}
	wantPrefix := "user:secret@tcp(tidb.example.com:4000)/lurnic?"
	if got[:len(wantPrefix)] != wantPrefix {
		t.Fatalf("prefix = %q, want %q", got, wantPrefix)
	}
	if !containsAll(got, "tls=skip-verify", "timeout=5s") {
		t.Fatalf("missing defaults in %q", got)
	}
}

func TestNormalizeGoDSN(t *testing.T) {
	got, err := Normalize("user:pass@tcp(db.example.com:3306)/lurnic?charset=utf8mb4")
	if err != nil {
		t.Fatal(err)
	}
	if !containsAll(got, "user:pass@tcp(db.example.com:3306)/lurnic", "tls=skip-verify") {
		t.Fatalf("unexpected DSN: %q", got)
	}
}

func containsAll(s string, parts ...string) bool {
	for _, p := range parts {
		if !contains(s, p) {
			return false
		}
	}
	return true
}

func contains(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
