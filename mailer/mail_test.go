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

	expected := "From: shutt@example.com\nTo: test@example.com; 123@google.com; auth@gmail.com; \nSubject: Hi there\n\nThis is a text message"

	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Fatalf("expected no diff but got: %s", diff)
	}
}
