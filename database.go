package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func initDb() (*sql.DB, error) {
	// Get absolute path to the executable's directory
	exePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("get executable path: %w", err)
	}
	exeDir := filepath.Dir(exePath)

	// Use absolute paths for database and seed file
	dbPath := filepath.Join(exeDir, "app.db")
	seedPath := filepath.Join(exeDir, "neetcode_150.json")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// TODO: Remove this - temporary for testing
	// db.Exec("DROP TABLE IF EXISTS completions;")
	// db.Exec("DROP TABLE IF EXISTS problems;")

	if err := createTables(db); err != nil {
		db.Close()
		return nil, err
	}

	// Seed NeetCode 150 problems on first run (only if table is empty)
	if err := seedNeetCodeFromJSON(db, seedPath); err != nil {
		// Seeding is best-effort; if file missing, just continue with a note
		if !errors.Is(err, os.ErrNotExist) {
			db.Close()
			return nil, fmt.Errorf("failed to seed problems: %w", err)
		}
		fmt.Println("ℹ neetcode_150.json not found; skipping initial seed")
	}

	fmt.Println("✓ Database initialized")
	return db, nil
}

func createTables(db *sql.DB) error {
	createProblemsTable := `
		CREATE TABLE IF NOT EXISTS problems (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL UNIQUE,
			grouping TEXT,
			leetcode_number INTEGER,
			difficulty TEXT,
			notes TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`

	createCompletionsTable := `
		CREATE TABLE IF NOT EXISTS completions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			problem_id INTEGER NOT NULL,
			effort_rating INTEGER NOT NULL,
			interval_days INTEGER DEFAULT 1,
			easiness_factor REAL DEFAULT 2.5,
			repetitions INTEGER DEFAULT 0,
			next_review_date DATETIME,
			completed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (problem_id) REFERENCES problems(id)
		);`

	if _, err := db.Exec(createProblemsTable); err != nil {
		return fmt.Errorf("create problems table: %w", err)
	}

	if _, err := db.Exec(createCompletionsTable); err != nil {
		return fmt.Errorf("create completions table: %w", err)
	}

	return nil
}
