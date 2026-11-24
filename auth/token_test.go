package auth

import (
	"testing"
	"time"
)

func TestTokenRoundTrip(t *testing.T) {
	m := NewTokenManager("test-secret-0123456789", time.Hour)
	tok, err := m.Issue(42)
	if err != nil {
		t.Fatal(err)
	}
	id, err := m.Parse(tok)
	if err != nil || id != 42 {
		t.Fatalf("got %d, %v", id, err)
	}

	if _, err := NewTokenManager("another-secret-0123456789", time.Hour).Parse(tok); err == nil {
		t.Fatal("token signed with a different secret must be rejected")
	}

	m.now = func() time.Time { return time.Now().Add(2 * time.Hour) }
	if _, err := m.Parse(tok); err == nil {
		t.Fatal("expired token must be rejected")
	}
}
