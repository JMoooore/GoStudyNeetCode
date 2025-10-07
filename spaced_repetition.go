package main

import (
	"database/sql"
	"fmt"
)

// ==================== Spaced Repetition (SM-2 Algorithm) ====================

func updateProblemCompletion(db *sql.DB, title string, effortRating int) error {
	problemID, err := getProblemID(db, title)
	if err != nil {
		return err
	}

	lastEF, lastInterval, lastReps := getLastCompletion(db, problemID)
	newInterval, newEF, newReps := calculateSM2(effortRating, lastEF, lastInterval, lastReps)

	return insertCompletion(db, problemID, effortRating, newInterval, newEF, newReps)
}

func getProblemID(db *sql.DB, title string) (int, error) {
	var problemID int
	err := db.QueryRow("SELECT id FROM problems WHERE title = ?", title).Scan(&problemID)
	if err != nil {
		return 0, fmt.Errorf("find problem: %w", err)
	}
	return problemID, nil
}

func getLastCompletion(db *sql.DB, problemID int) (ef float64, interval, reps int) {
	ef, interval, reps = 2.5, 1, 0
	db.QueryRow(`
		SELECT easiness_factor, interval_days, repetitions
		FROM completions
		WHERE problem_id = ?
		ORDER BY completed_at DESC
		LIMIT 1
	`, problemID).Scan(&ef, &interval, &reps)
	return
}

func calculateSM2(effortRating int, lastEF float64, lastInterval, lastReps int) (interval int, newEF float64, reps int) {
	// effortRating: 1=Easy, 2=Medium, 3=Hard
	// SM-2 quality scale: Easy=5, Medium=3, Hard=1
	quality := map[int]int{1: 5, 2: 3, 3: 1}[effortRating]

	// Calculate new easiness factor
	newEF = lastEF + (0.1 - float64(5-quality)*(0.08+float64(5-quality)*0.02))
	if newEF < 1.3 {
		newEF = 1.3
	}

	// Calculate interval and repetitions
	if quality < 3 {
		// Failed - reset to beginning
		interval, reps = 1, 0
	} else {
		reps = lastReps + 1
		switch reps {
		case 1:
			// First review: scale by quality
			if quality == 5 { // Easy
				interval = 4
			} else { // Medium (quality == 3)
				interval = 2
			}
		case 2:
			// Second review: scale by quality
			if quality == 5 { // Easy
				interval = 14
			} else { // Medium
				interval = 7
			}
		default:
			interval = int(float64(lastInterval) * newEF)
		}
	}

	return interval, newEF, reps
}

func insertCompletion(db *sql.DB, problemID, effortRating, interval int, ef float64, reps int) error {
	_, err := db.Exec(`
		INSERT INTO completions (problem_id, effort_rating, interval_days, easiness_factor, repetitions, next_review_date)
		VALUES (?, ?, ?, ?, ?, datetime('now', '+' || ? || ' days'))
	`, problemID, effortRating, interval, ef, reps, interval)
	return err
}
