package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
)

// ==================== Seeding ====================

func seedNeetCodeFromJSON(db *sql.DB, jsonPath string) error {
	var count int
	if err := db.QueryRow("SELECT COUNT(1) FROM problems").Scan(&count); err != nil {
		return fmt.Errorf("count problems: %w", err)
	}
	if count > 0 {
		return nil // Already seeded
	}

	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return err
	}

	var problems []Problem
	if err := json.Unmarshal(data, &problems); err != nil {
		return fmt.Errorf("parse %s: %w", jsonPath, err)
	}

	if err := insertProblems(db, problems); err != nil {
		return err
	}

	fmt.Printf("✓ Seeded %d NeetCode problems\n", len(problems))
	return nil
}

func insertProblems(db *sql.DB, problems []Problem) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT OR IGNORE INTO problems (title, difficulty, grouping, leetcode_number) VALUES (?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, p := range problems {
		if p.Title == "" {
			continue
		}
		if _, err := stmt.Exec(p.Title, p.Difficulty, p.Grouping, p.LeetcodeNumber); err != nil {
			return fmt.Errorf("insert problem %q: %w", p.Title, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
