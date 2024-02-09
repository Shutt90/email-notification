package mail

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBuildMessage(t *testing.T) {
	actual := BuildMessage(
		[]string{"test@example.com", "123@google.com", "auth@gmail.com"},
		"shutt@example.com",
		"Hi there",
		"This is a text message",
	)

	expected := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\nFrom: shutt@example.com\nTo: test@example.com;123@google.com;auth@gmail.com\nSubject: Hi there!\nThis is a text message"

	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Fatalf("expected no diff but got: %s", diff)
	}
}
