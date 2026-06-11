package database

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		t.Skip("DATABASE_URL not set")
	}

	db, err := New(databaseURL)
	if err != nil {
		t.Fatalf("failed to create database service: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.DB.Ping(ctx); err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}
}

func TestNew_InvalidURL(t *testing.T) {
	_, err := New("not-a-valid-postgres-url")

	if err == nil {
		t.Fatal("expected error for invalid database url")
	}
}

func TestDatabaseQuery(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		t.Skip("DATABASE_URL not set")
	}

	db, err := New(databaseURL)
	if err != nil {
		t.Fatalf("failed to create database service: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result int

	err = db.DB.QueryRow(
		ctx,
		"SELECT 1",
	).Scan(&result)

	if err != nil {
		t.Fatalf("query failed: %v", err)
	}

	if result != 1 {
		t.Fatalf("expected 1, got %d", result)
	}
}

func TestDatabaseVersion(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		t.Skip("DATABASE_URL not set")
	}

	db, err := New(databaseURL)
	if err != nil {
		t.Fatalf("failed to create database service: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var version string

	err = db.DB.QueryRow(
		ctx,
		"SELECT version()",
	).Scan(&version)

	if err != nil {
		t.Fatalf("failed to query postgres version: %v", err)
	}

	if version == "" {
		t.Fatal("postgres version is empty")
	}

	t.Logf("PostgreSQL Version: %s", version)
}