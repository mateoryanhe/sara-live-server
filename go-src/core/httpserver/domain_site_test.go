package httpserver

import (
	"context"
	"path/filepath"
	"testing"
)

func TestRefreshDomainSiteRegistryAtomicallyReplacesSnapshot(t *testing.T) {
	original := domainSiteRegistry.Load()
	t.Cleanup(func() {
		domainSiteRegistry.Store(original)
	})

	rootOne := t.TempDir()
	rootTwo := t.TempDir()
	active, err := RefreshDomainSiteRegistry(context.Background(), []DomainSiteRegistration{
		{Domain: "One.Example.com, two.example.com.", Root: rootOne},
		{Domain: "three.example.com", Root: rootTwo},
	})
	if err != nil {
		t.Fatalf("RefreshDomainSiteRegistry() error = %v", err)
	}
	if len(active) != 3 {
		t.Fatalf("active registrations = %d, want 3", len(active))
	}
	snapshot := domainSiteRegistry.Load()
	if snapshot == nil {
		t.Fatal("domain registry snapshot is nil")
	}
	if got := snapshot.roots["one.example.com"]; got != filepath.Clean(rootOne) {
		t.Fatalf("one.example.com root = %q, want %q", got, filepath.Clean(rootOne))
	}
	if got := snapshot.roots["two.example.com"]; got != filepath.Clean(rootOne) {
		t.Fatalf("two.example.com root = %q, want %q", got, filepath.Clean(rootOne))
	}
	if got := snapshot.roots["three.example.com"]; got != filepath.Clean(rootTwo) {
		t.Fatalf("three.example.com root = %q, want %q", got, filepath.Clean(rootTwo))
	}

	beforeInvalidRefresh := snapshot
	_, err = RefreshDomainSiteRegistry(context.Background(), []DomainSiteRegistration{
		{Domain: "duplicate.example.com", Root: rootOne},
		{Domain: "DUPLICATE.EXAMPLE.COM.", Root: rootTwo},
	})
	if err == nil {
		t.Fatal("duplicate domain refresh unexpectedly succeeded")
	}
	if got := domainSiteRegistry.Load(); got != beforeInvalidRefresh {
		t.Fatal("invalid refresh replaced the active snapshot")
	}
}

func TestNormalizeRegisteredDomain(t *testing.T) {
	valid := map[string]string{
		"Example.COM.": "example.com",
		"localhost":    "localhost",
		"127.0.0.1":    "127.0.0.1",
	}
	for input, want := range valid {
		got, err := normalizeRegisteredDomain(input)
		if err != nil {
			t.Fatalf("normalizeRegisteredDomain(%q) error = %v", input, err)
		}
		if got != want {
			t.Fatalf("normalizeRegisteredDomain(%q) = %q, want %q", input, got, want)
		}
	}

	invalid := []string{
		"https://example.com",
		"example.com/path",
		"example.com:443",
		"*.example.com",
		"bad..example.com",
		"-bad.example.com",
	}
	for _, input := range invalid {
		if _, err := normalizeRegisteredDomain(input); err == nil {
			t.Fatalf("normalizeRegisteredDomain(%q) unexpectedly succeeded", input)
		}
	}
}

func TestReplaceDomainSiteRegistrationKeepsOtherSites(t *testing.T) {
	original := domainSiteRegistry.Load()
	t.Cleanup(func() {
		domainSiteRegistry.Store(original)
	})

	oldRoot := t.TempDir()
	otherRoot := t.TempDir()
	newRoot := t.TempDir()
	_, err := RefreshDomainSiteRegistry(context.Background(), []DomainSiteRegistration{
		{Domain: "old.example.com", Root: oldRoot},
		{Domain: "other.example.com", Root: otherRoot},
	})
	if err != nil {
		t.Fatalf("initial refresh error = %v", err)
	}
	active, err := ReplaceDomainSiteRegistration(context.Background(), []string{"old.example.com"}, DomainSiteRegistration{
		Domain: "new.example.com",
		Root:   newRoot,
	})
	if err != nil {
		t.Fatalf("ReplaceDomainSiteRegistration() error = %v", err)
	}
	if len(active) != 2 {
		t.Fatalf("active registrations = %d, want 2", len(active))
	}
	snapshot := domainSiteRegistry.Load()
	if _, exists := snapshot.roots["old.example.com"]; exists {
		t.Fatal("old domain still exists after replacement")
	}
	if got := snapshot.roots["new.example.com"]; got != filepath.Clean(newRoot) {
		t.Fatalf("new root = %q, want %q", got, filepath.Clean(newRoot))
	}
	if got := snapshot.roots["other.example.com"]; got != filepath.Clean(otherRoot) {
		t.Fatalf("other root = %q, want %q", got, filepath.Clean(otherRoot))
	}
}
