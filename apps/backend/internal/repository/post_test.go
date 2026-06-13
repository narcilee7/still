package repository

import (
	"testing"

	stillv1 "github.com/still-mvp/still/apps/backend/gen/still/v1"
)

func TestStatusToString(t *testing.T) {
	cases := []struct {
		input    stillv1.PostStatus
		expected string
	}{
		{stillv1.PostStatus_POST_STATUS_PENDING, "PENDING"},
		{stillv1.PostStatus_POST_STATUS_APPROVED, "APPROVED"},
		{stillv1.PostStatus_POST_STATUS_REJECTED, "REJECTED"},
		{stillv1.PostStatus_POST_STATUS_UNSPECIFIED, "PENDING"},
	}
	for _, c := range cases {
		got := StatusToString(c.input)
		if got != c.expected {
			t.Errorf("StatusToString(%v) = %q, want %q", c.input, got, c.expected)
		}
	}
}

func TestParseStatus(t *testing.T) {
	cases := []struct {
		input    string
		expected stillv1.PostStatus
	}{
		{"PENDING", stillv1.PostStatus_POST_STATUS_PENDING},
		{"APPROVED", stillv1.PostStatus_POST_STATUS_APPROVED},
		{"REJECTED", stillv1.PostStatus_POST_STATUS_REJECTED},
		{"unknown", stillv1.PostStatus_POST_STATUS_UNSPECIFIED},
		{"", stillv1.PostStatus_POST_STATUS_UNSPECIFIED},
	}
	for _, c := range cases {
		got := parseStatus(c.input)
		if got != c.expected {
			t.Errorf("parseStatus(%q) = %v, want %v", c.input, got, c.expected)
		}
	}
}
