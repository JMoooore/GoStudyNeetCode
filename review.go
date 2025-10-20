package main

import (
	"database/sql"
	"fmt"
	"time"
)

// ==================== Review History ====================

type ReviewInfo struct {
	Title           string
	Difficulty      string
	LastCompletedAt sql.NullString
	NextReviewDate  sql.NullString
	Repetitions     sql.NullInt64
	EasinessFactor  sql.NullFloat64
	DaysUntilReview sql.NullInt64
}

func getReviewHistory(db *sql.DB, difficulty string) ([]ReviewInfo, error) {
	query := `
		SELECT
			p.title,
			p.difficulty,
			c.completed_at,
			c.next_review_date,
			c.repetitions,
			c.easiness_factor,
			CAST((julianday(date(c.next_review_date)) - julianday(date('now'))) AS INTEGER) as days_until
		FROM problems p
		LEFT JOIN (
			SELECT
				problem_id,
				MAX(completed_at) as completed_at,
				next_review_date,
				repetitions,
				easiness_factor
			FROM completions
			GROUP BY problem_id
		) c ON p.id = c.problem_id
		WHERE c.completed_at IS NOT NULL
	`

	if difficulty != "any" {
		query += " AND LOWER(p.difficulty) = LOWER(?)"
	}

	query += " ORDER BY c.next_review_date ASC, p.title ASC"

	var rows *sql.Rows
	var err error
	if difficulty == "any" {
		rows, err = db.Query(query)
	} else {
		rows, err = db.Query(query, difficulty)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []ReviewInfo
	for rows.Next() {
		var r ReviewInfo
		if err := rows.Scan(&r.Title, &r.Difficulty, &r.LastCompletedAt, &r.NextReviewDate,
			&r.Repetitions, &r.EasinessFactor, &r.DaysUntilReview); err != nil {
			return nil, fmt.Errorf("scan review: %w", err)
		}
		reviews = append(reviews, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return reviews, nil
}

func formatReviewDate(dateStr sql.NullString) string {
	if !dateStr.Valid {
		return "Never"
	}

	// Try multiple date formats
	formats := []string{
		"2006-01-02 15:04:05",
		time.RFC3339,
		"2006-01-02T15:04:05Z",
	}

	var t time.Time
	var err error
	for _, format := range formats {
		t, err = time.Parse(format, dateStr.String)
		if err == nil {
			break
		}
	}

	if err != nil {
		return dateStr.String
	}

	return t.Format("Jan 2, 2006")
}

func getReviewStatus(days sql.NullInt64) (string, string) {
	if !days.Valid {
		return "⚪", "Never attempted"
	}

	d := days.Int64
	if d < 0 {
		return "🔴", fmt.Sprintf("Overdue by %d days", -d)
	} else if d == 0 {
		return "🟠", "Due today"
	} else if d <= 3 {
		return "🟡", fmt.Sprintf("Due in %d days", d)
	} else {
		return "🟢", fmt.Sprintf("Due in %d days", d)
	}
}
