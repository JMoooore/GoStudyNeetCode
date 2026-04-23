package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
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
	err := db.QueryRow(`
		SELECT easiness_factor, interval_days, repetitions
		FROM completions
		WHERE problem_id = ?
		ORDER BY completed_at DESC
		LIMIT 1
	`, problemID).Scan(&ef, &interval, &reps)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		fmt.Fprintf(os.Stderr, "warning: fetching last completion for problem %d: %v\n", problemID, err)
	}
	return
}

func calculateSM2(effortRating int, lastEF float64, lastInterval, lastReps int) (interval int, newEF float64, reps int) {
	quality := map[int]int{1: 5, 2: 3, 3: 1}[effortRating]

	newEF = lastEF + (0.1 - float64(5-quality)*(0.08+float64(5-quality)*0.02))
	if newEF < 1.3 {
		newEF = 1.3
	}

	if quality < 3 {
		// Hard - don't fully reset, just short interval
		interval, reps = 3, lastReps
	} else {
		reps = lastReps + 1
		switch reps {
		case 1:
			if quality == 5 {
				interval = 7
			} else {
				interval = 4
			}
		case 2:
			if quality == 5 {
				interval = 21
			} else {
				interval = 14
			}
		default:
			// Use the new interval based on lastInterval and newEF
			interval = int(float64(lastInterval) * newEF)
			// Floor to prevent reviews piling up - but intervals will still grow when above minimum
			if interval < 14 {
				interval = 14
			}
		}
	}

	return interval, newEF, reps
}

func insertCompletion(db *sql.DB, problemID, effortRating, interval int, ef float64, reps int) error {
	_, err := db.Exec(`
		INSERT INTO completions (problem_id, effort_rating, interval_days, easiness_factor, repetitions, next_review_date, completed_at)
		VALUES (?, ?, ?, ?, ?, datetime('now', 'localtime', '+' || ? || ' days'), datetime('now', 'localtime'))
	`, problemID, effortRating, interval, ef, reps, interval)
	return err
}

func archiveProblem(db *sql.DB, title string) error {
	problemID, err := getProblemID(db, title)
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE problems SET archived = 1 WHERE id = ?", problemID)
	return err
}
