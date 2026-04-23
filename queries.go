package main

import (
	"database/sql"
	"fmt"
)

// ==================== Problem Queries ====================

func selectStudyProblems(db *sql.DB, difficulty string, category string, count int) ([]Problem, error) {
	query := buildStudyQuery(difficulty, category)

	var rows *sql.Rows
	var err error
	
	// Build args based on what filters are used
	args := []interface{}{}
	if difficulty != "any" {
		args = append(args, difficulty)
	}
	if category != "any" {
		args = append(args, category)
	}
	args = append(args, count)
	
	rows, err = db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var problems []Problem
	for rows.Next() {
		var p Problem
		if err := rows.Scan(&p.Title, &p.Difficulty, &p.Grouping, &p.LeetcodeNumber, &p.Category); err != nil {
			return nil, fmt.Errorf("scan problem: %w", err)
		}
		problems = append(problems, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return problems, nil
}

func buildStudyQuery(difficulty string, category string) string {
	// Prioritizes:
	// 1. Reviews due today or past (next_review_date <= now) - oldest first
	// 2. Never attempted problems (new)
	// 3. Reviews upcoming (next_review_date > now) - nearest first
	query := `
		SELECT p.title, p.difficulty, p.grouping, p.leetcode_number, p.category
		FROM problems p
		LEFT JOIN (
			SELECT problem_id, completed_at as last_completion, next_review_date
			FROM (
				SELECT problem_id, completed_at, next_review_date,
					ROW_NUMBER() OVER (PARTITION BY problem_id ORDER BY completed_at DESC) as rn
				FROM completions
			)
			WHERE rn = 1
		) c ON p.id = c.problem_id
		WHERE p.archived = 0
	`

	// Add difficulty filter
	if difficulty != "any" {
		query += " AND LOWER(p.difficulty) = LOWER(?)"
	}
	
	// Add category filter
	if category != "any" {
		query += " AND LOWER(p.category) = LOWER(?)"
	}

	query += ` ORDER BY
		CASE
			WHEN date(c.next_review_date) <= date('now', 'localtime') THEN 1
			WHEN c.next_review_date IS NULL THEN 2
			ELSE 3
		END,
		CASE
			WHEN date(c.next_review_date) <= date('now', 'localtime') THEN c.next_review_date
			ELSE NULL
		END ASC,
		RANDOM()
		LIMIT ?`
	return query
}
