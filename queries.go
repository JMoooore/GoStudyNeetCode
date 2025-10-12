package main

import (
	"database/sql"
	"fmt"
)

// ==================== Problem Queries ====================

func selectStudyProblems(db *sql.DB, difficulty string, count int) ([]Problem, error) {
	query := buildStudyQuery(difficulty)

	var rows *sql.Rows
	var err error
	if difficulty == "any" {
		rows, err = db.Query(query, count)
	} else {
		rows, err = db.Query(query, difficulty, count)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var problems []Problem
	for rows.Next() {
		var p Problem
		if err := rows.Scan(&p.Title, &p.Difficulty, &p.Grouping, &p.LeetcodeNumber); err != nil {
			return nil, fmt.Errorf("scan problem: %w", err)
		}
		problems = append(problems, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return problems, nil
}

func buildStudyQuery(difficulty string) string {
	// Prioritizes:
	// 1. Reviews due today or past (next_review_date <= now) - oldest first
	// 2. Never attempted problems (new)
	// 3. Reviews upcoming (next_review_date > now) - nearest first
	query := `
		SELECT p.title, p.difficulty, p.grouping, p.leetcode_number
		FROM problems p
		LEFT JOIN (
			SELECT problem_id, MAX(completed_at) as last_completion, next_review_date
			FROM completions
			GROUP BY problem_id
		) c ON p.id = c.problem_id
	`

	if difficulty != "any" {
		query += " WHERE LOWER(p.difficulty) = LOWER(?)"
	}

	query += ` ORDER BY
		CASE
			WHEN c.next_review_date <= datetime('now') THEN 1
			WHEN c.next_review_date IS NULL THEN 2
			ELSE 3
		END,
		c.next_review_date ASC
		LIMIT ?`
	return query
}
