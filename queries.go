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
	// 1. Problems due for review (next_review_date <= now)
	// 2. Never attempted problems
	// 3. Ordered by next_review_date (oldest first)
	query := `
		SELECT p.title, p.difficulty, p.grouping, p.leetcode_number
		FROM problems p
		LEFT JOIN (
			SELECT problem_id, MAX(completed_at) as last_completion, next_review_date
			FROM completions
			GROUP BY problem_id
		) c ON p.id = c.problem_id
		WHERE (c.next_review_date IS NULL OR c.next_review_date <= datetime('now'))
	`

	if difficulty != "any" {
		query += " AND LOWER(p.difficulty) = LOWER(?)"
	}

	query += " ORDER BY c.next_review_date IS NULL DESC, c.next_review_date ASC LIMIT ?"
	return query
}
