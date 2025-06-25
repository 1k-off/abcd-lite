package oras

import (
	"strings"
	"testing"
)

func TestPullArgs_Valid(t *testing.T) {
	cmd := OrasCommand{
		Type:        Pull,
		Reference:   "myregistry.com/myimage:1.0.0",
		Username:    "user",
		Password:    "pass",
		Output:      "/tmp/out",
		Concurrency: 2,
	}
	args, err := cmd.Args()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{
		"pull", "--output", "/tmp/out", "--concurrency", "2",
		"--username", "user", "--password", "pass",
		"myregistry.com/myimage:1.0.0",
	}
	if strings.Join(args, " ") != strings.Join(expected, " ") {
		t.Errorf("expected %v, got %v", expected, args)
	}
}

func TestPullArgs_MissingReference(t *testing.T) {
	cmd := OrasCommand{
		Type: Pull,
	}
	_, err := cmd.Args()
	if err == nil || !strings.Contains(err.Error(), "reference must be specified") {
		t.Errorf("expected reference error, got %v", err)
	}
}

func TestPullArgs_Minimal(t *testing.T) {
	cmd := OrasCommand{
		Type:      Pull,
		Reference: "repo/image:tag",
	}
	args, err := cmd.Args()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"pull", "repo/image:tag"}
	if strings.Join(args, " ") != strings.Join(expected, " ") {
		t.Errorf("expected %v, got %v", expected, args)
	}
}
