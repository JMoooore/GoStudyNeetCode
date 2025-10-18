package main

import (
	"database/sql"
	"fmt"
	"time"
)

// ==================== Stats ====================

type OverallStats struct {
	// Overall completion counts
	TotalProblems     int
	CompletedProblems int
	RemainingProblems int

	// By difficulty
	EasyTotal       int
	EasyCompleted   int
	MediumTotal     int
	MediumCompleted int
	HardTotal       int
	HardCompleted   int

	// Review stats
	ProblemsNeedReview int
	OverdueReviews     int
	DueTodayReviews    int
	UpcomingReviews    int // Due within 3 days

	// Projection stats
	EstimatedDaysToComplete int // At 3 problems/day
	EstimatedCompletionDate string
}

func getOverallStats(db *sql.DB) (*OverallStats, error) {
	stats := &OverallStats{}

	// Get total problems by difficulty
	diffQuery := `
		SELECT
			LOWER(difficulty) as diff,
			COUNT(*) as total
		FROM problems
		GROUP BY LOWER(difficulty)
	`
	rows, err := db.Query(diffQuery)
	if err != nil {
		return nil, fmt.Errorf("query difficulty totals: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var diff string
		var total int
		if err := rows.Scan(&diff, &total); err != nil {
			return nil, err
		}
		stats.TotalProblems += total
		switch diff {
		case "easy":
			stats.EasyTotal = total
		case "medium":
			stats.MediumTotal = total
		case "hard":
			stats.HardTotal = total
		}
	}

	// Get completed problems by difficulty
	completedQuery := `
		SELECT
			LOWER(p.difficulty) as diff,
			COUNT(DISTINCT p.id) as completed
		FROM problems p
		INNER JOIN completions c ON p.id = c.problem_id
		GROUP BY LOWER(p.difficulty)
	`
	rows, err = db.Query(completedQuery)
	if err != nil {
		return nil, fmt.Errorf("query completed problems: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var diff string
		var completed int
		if err := rows.Scan(&diff, &completed); err != nil {
			return nil, err
		}
		stats.CompletedProblems += completed
		switch diff {
		case "easy":
			stats.EasyCompleted = completed
		case "medium":
			stats.MediumCompleted = completed
		case "hard":
			stats.HardCompleted = completed
		}
	}

	stats.RemainingProblems = stats.TotalProblems - stats.CompletedProblems

	// Get review status counts
	reviewQuery := `
		SELECT
			CASE
				WHEN julianday(c.next_review_date) - julianday(date('now')) < 0 THEN 'overdue'
				WHEN julianday(c.next_review_date) - julianday(date('now')) = 0 THEN 'today'
				WHEN julianday(c.next_review_date) - julianday(date('now')) <= 3 THEN 'upcoming'
				ELSE 'future'
			END as status,
			COUNT(*) as count
		FROM problems p
		INNER JOIN (
			SELECT problem_id, MAX(completed_at) as max_completed, next_review_date
			FROM completions
			GROUP BY problem_id
		) c ON p.id = c.problem_id
		WHERE c.next_review_date IS NOT NULL
		GROUP BY status
	`
	rows, err = db.Query(reviewQuery)
	if err != nil {
		return nil, fmt.Errorf("query review status: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		switch status {
		case "overdue":
			stats.OverdueReviews = count
		case "today":
			stats.DueTodayReviews = count
		case "upcoming":
			stats.UpcomingReviews = count
		}
	}

	stats.ProblemsNeedReview = stats.OverdueReviews + stats.DueTodayReviews + stats.UpcomingReviews

	// Calculate estimated days to complete
	// Assumption: 3 problems per day total (including both new problems and reviews)
	// Each new problem marked as "easy" generates 2 reviews (at day 4 and day 14)
	// Total work = remaining problems + (2 reviews per new problem) + current review backlog
	const reviewsPerProblem = 2.0

	totalWork := float64(stats.RemainingProblems)*(1.0+reviewsPerProblem) + float64(stats.ProblemsNeedReview)
	stats.EstimatedDaysToComplete = int((totalWork + 2.0) / 3.0) // Round up

	// Calculate estimated completion date
	estimatedDate := time.Now().AddDate(0, 0, stats.EstimatedDaysToComplete)
	stats.EstimatedCompletionDate = estimatedDate.Format("Jan 2, 2006")

	return stats, nil
}

func statCommandWithDB(db *sql.DB, args []string) error {
	stats, err := getOverallStats(db)
	if err != nil {
		return fmt.Errorf("get stats: %w", err)
	}

	fmt.Println()
	fmt.Println("📊 Your Study Statistics")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	// Overall Progress
	fmt.Println("Overall Progress:")
	fmt.Println("─────────────────────────────────────────────────────────")
	totalPercent := 0.0
	if stats.TotalProblems > 0 {
		totalPercent = float64(stats.CompletedProblems) / float64(stats.TotalProblems) * 100
	}
	fmt.Printf("  Total:      %d / %d problems completed (%.1f%%)\n",
		stats.CompletedProblems, stats.TotalProblems, totalPercent)
	fmt.Printf("  Remaining:  %d problems\n", stats.RemainingProblems)
	fmt.Println()

	// Progress by Difficulty
	fmt.Println("Progress by Difficulty:")
	fmt.Println("─────────────────────────────────────────────────────────")

	// Easy
	easyPercent := 0.0
	if stats.EasyTotal > 0 {
		easyPercent = float64(stats.EasyCompleted) / float64(stats.EasyTotal) * 100
	}
	fmt.Printf("  \033[32mEasy:\033[0m     %3d / %-3d (%5.1f%%)  %s\n",
		stats.EasyCompleted, stats.EasyTotal, easyPercent, progressBar(easyPercent, "green"))

	// Medium
	mediumPercent := 0.0
	if stats.MediumTotal > 0 {
		mediumPercent = float64(stats.MediumCompleted) / float64(stats.MediumTotal) * 100
	}
	fmt.Printf("  \033[33mMedium:\033[0m   %3d / %-3d (%5.1f%%)  %s\n",
		stats.MediumCompleted, stats.MediumTotal, mediumPercent, progressBar(mediumPercent, "yellow"))

	// Hard
	hardPercent := 0.0
	if stats.HardTotal > 0 {
		hardPercent = float64(stats.HardCompleted) / float64(stats.HardTotal) * 100
	}
	fmt.Printf("  \033[31mHard:\033[0m     %3d / %-3d (%5.1f%%)  %s\n",
		stats.HardCompleted, stats.HardTotal, hardPercent, progressBar(hardPercent, "red"))
	fmt.Println()

	// Review Status
	fmt.Println("Review Status:")
	fmt.Println("─────────────────────────────────────────────────────────")
	fmt.Printf("  🔴 Overdue:  %d problems\n", stats.OverdueReviews)
	fmt.Printf("  🟠 Today:    %d problems\n", stats.DueTodayReviews)
	fmt.Printf("  🟡 Soon:     %d problems (within 3 days)\n", stats.UpcomingReviews)
	fmt.Println()

	// Projections
	fmt.Println("Projections (at 3 problems/day assuming \033[32mEasy\033[0m completions):")
	fmt.Println("─────────────────────────────────────────────────────────")
	fmt.Printf("  Estimated days to complete:  %d days\n", stats.EstimatedDaysToComplete)
	fmt.Printf("  Estimated completion date:   %s\n", stats.EstimatedCompletionDate)
	fmt.Println()

	return nil
}

func progressBar(percent float64, color string) string {
	barLength := 20
	filled := min(int(percent/100.0*float64(barLength)), barLength)

	// ANSI color codes
	colorCodes := map[string]string{
		"green":  "\033[32m",
		"yellow": "\033[33m",
		"red":    "\033[31m",
		"cyan":   "\033[36m",
		"reset":  "\033[0m",
	}

	colorCode := colorCodes[color]
	if colorCode == "" {
		colorCode = colorCodes["cyan"]
	}
	resetCode := colorCodes["reset"]

	bar := "["
	for i := range barLength {
		if i < filled {
			bar += colorCode + "█" + resetCode
		} else {
			bar += "░"
		}
	}
	bar += "]"
	return bar
}
