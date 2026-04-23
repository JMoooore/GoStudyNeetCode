package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

var shortToLong = map[string]string{
	"e": "easy",
	"m": "medium",
	"h": "hard",
	"a": "any",
}

func helpCommand(args []string) error {
	fmt.Println()
	fmt.Println("Available Commands:")
	fmt.Println("==================")

	cmds := getCommands(nil)
	seen := make(map[string]bool)
	names := make([]string, 0, len(cmds))
	for _, cmd := range cmds {
		if seen[cmd.Name] {
			continue
		}
		seen[cmd.Name] = true
		names = append(names, cmd.Name)
	}
	sort.Strings(names)
	for _, name := range names {
		cmd := cmds[name]
		aliasStr := ""
		if len(cmd.Aliases) > 0 {
			aliasStr = fmt.Sprintf("  (aliases: %s)", strings.Join(cmd.Aliases, ", "))
		}
		fmt.Printf("  %-10s - %s%s\n", cmd.Name, cmd.Description, aliasStr)
	}
	fmt.Println()
	fmt.Println("Example usage:")
	fmt.Println("  study --difficulty easy --count 5")
	fmt.Println("  study -d medium -c 3")
	fmt.Println("  study --category data_structures --count 3")
	fmt.Println("  study -t neetcode -d hard -c 2")
	fmt.Println()
	return nil
}

func exitCommand(args []string) error {
	fmt.Println("Thanks for using GoStudyNeetCode! Happy coding! 👋")
	os.Exit(0)
	return nil
}

func studyCommandWithDB(db *sql.DB, args []string) error {
	// Create a new FlagSet for this command
	fs := flag.NewFlagSet("study", flag.ContinueOnError)

	var difficulty string
	var category string
	var count int

	// Define flags
	fs.StringVar(&difficulty, "difficulty", "any", "Difficulty level (easy, medium, hard, any OR e, m, h, a)")
	fs.StringVar(&difficulty, "d", "any", "Short for difficulty")
	
	fs.StringVar(&category, "category", "any", "Category (neetcode, data_structures, any)")
	fs.StringVar(&category, "type", "any", "Short for category")
	fs.StringVar(&category, "t", "any", "Short for category")

	fs.IntVar(&count, "count", 1, "Number of questions")
	fs.IntVar(&count, "c", 1, "Short for count")

	// Parse the flags
	err := fs.Parse(args)
	if err != nil {
		return err
	}

	// Convert short form to long form if needed
	if val, ok := shortToLong[difficulty]; ok {
		difficulty = val
	}
	
	// Normalize category
	if category != "any" {
		category = strings.ToLower(category)
		if category != "neetcode" && category != "data_structures" {
			return fmt.Errorf("invalid category: %s (must be 'neetcode', 'data_structures', or 'any')", category)
		}
	}

	problems, err := selectStudyProblems(db, difficulty, category, count)
	if err != nil {
		return err
	}

	if len(problems) == 0 {
		fmt.Println("\n📚 No problems found matching your criteria.")
		fmt.Println()
		return nil
	}

	fmt.Println("\n📚 Your Study Problems:")
	fmt.Println("========================")
	for i, p := range problems {
		lcStr := ""
		if p.LeetcodeNumber > 0 {
			lcStr = fmt.Sprintf("[LC %d] ", p.LeetcodeNumber)
		}
		categoryLabel := ""
		if p.Category == "data_structures" {
			categoryLabel = " [Data Structure]"
		}
		fmt.Printf("%d. %s%s (%s) - %s%s\n", i+1, lcStr, p.Title, p.Difficulty, p.Grouping, categoryLabel)
	}
	fmt.Println()

	// Ask if user wants to mark any as completed
	reader := bufio.NewReader(os.Stdin)
	for len(problems) > 0 {
		fmt.Print("Mark any as completed? (y/n): ")
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response == "n" || response == "no" {
			break
		}

		if response == "y" || response == "yes" {
			fmt.Print("Enter problem number (e.g. 1 or 3): ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if input != "" {
				numStr := strings.TrimSpace(input)
				num, err := strconv.Atoi(numStr)
				if err != nil || num < 1 || num > len(problems) {
					fmt.Printf("Invalid problem number: %s\n", numStr)
					continue
				}

				problem := problems[num-1]

				// Ask for effort rating
				fmt.Printf("\nHow hard was '%s'? (1=Easy, 2=Medium, 3=Hard, 4=Remove from rotation): ", problem.Title)
				ratingStr, _ := reader.ReadString('\n')
				ratingStr = strings.TrimSpace(ratingStr)
				rating, err := strconv.Atoi(ratingStr)
				if err != nil || rating < 1 || rating > 4 {
					fmt.Println("Invalid rating, skipping...")
					continue
				}

				// Handle remove from rotation
				if rating == 4 {
					if err := archiveProblem(db, problem.Title); err != nil {
						fmt.Printf("Error archiving problem: %v\n", err)
					} else {
						fmt.Printf("\033[33m✓ Removed '%s' from rotation\033[0m\n", problem.Title)
					}
				} else {
					// Update the database
					if err := updateProblemCompletion(db, problem.Title, rating); err != nil {
						fmt.Printf("Error updating problem: %v\n", err)
					} else {
						fmt.Printf("\033[32m✓ Marked '%s' as completed with effort rating %d\033[0m\n", problem.Title, rating)
					}
				}

				// Remove from slice using 0-indexed position
				idx := num - 1
				problems = append(problems[:idx], problems[idx+1:]...)

				// Show updated list
				if len(problems) > 0 {
					fmt.Println("\nRemaining problems:")
					for i, p := range problems {
						lcStr := ""
						if p.LeetcodeNumber > 0 {
							lcStr = fmt.Sprintf("[LC %d] ", p.LeetcodeNumber)
						}
						fmt.Printf("%d. %s%s\n", i+1, lcStr, p.Title)
					}
					fmt.Println()
				}
			}
		}
	}

	return nil
}

func reviewCommandWithDB(db *sql.DB, args []string) error {
	fs := flag.NewFlagSet("review", flag.ContinueOnError)

	var difficulty string
	fs.StringVar(&difficulty, "difficulty", "any", "Filter by difficulty (easy, medium, hard, any)")
	fs.StringVar(&difficulty, "d", "any", "Short for difficulty")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	if val, ok := shortToLong[difficulty]; ok {
		difficulty = val
	}

	reviews, err := getReviewHistory(db, difficulty)
	if err != nil {
		return err
	}

	if len(reviews) == 0 {
		fmt.Println("\nNo completed problems yet. Complete some problems first!")
		return nil
	}

	fmt.Println("\n📊 Review History:")
	fmt.Println("====================================================================================")
	fmt.Printf("%-4s %-40s %-10s %-15s %-15s\n", "Stat", "Problem", "Difficulty", "Last Done", "Next Review")
	fmt.Println("------------------------------------------------------------------------------------")

	for _, r := range reviews {
		status, _ := getReviewStatus(r.DaysUntilReview)
		lastDone := formatReviewDate(r.LastCompletedAt)
		nextReview := formatReviewDate(r.NextReviewDate)

		fmt.Printf("%s   %-40s %-10s %-15s %-15s\n",
			status,
			truncate(r.Title, 40),
			r.Difficulty,
			lastDone,
			nextReview,
		)
	}
	fmt.Println()

	return nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func factorsCommandWithDB(db *sql.DB, args []string) error {
	fs := flag.NewFlagSet("factors", flag.ContinueOnError)

	var difficulty string
	fs.StringVar(&difficulty, "difficulty", "any", "Filter by difficulty (easy, medium, hard, any)")
	fs.StringVar(&difficulty, "d", "any", "Short for difficulty")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	if val, ok := shortToLong[difficulty]; ok {
		difficulty = val
	}

	reviews, err := getReviewHistory(db, difficulty)
	if err != nil {
		return err
	}

	if len(reviews) == 0 {
		fmt.Println("\nNo completed problems yet. Complete some problems first!")
		return nil
	}

	// Calculate statistics
	var totalEF float64
	var countEF int
	var minEF, maxEF float64
	minEF = 999.0
	maxEF = 0.0

	// Count by EF ranges
	var lowEF, mediumEF, highEF, veryHighEF int
	// low: < 2.0, medium: 2.0-2.5, high: 2.5-3.0, veryHigh: > 3.0

	fmt.Println("\n📈 Your Learning Factors (Easiness Factors):")
	fmt.Println("═══════════════════════════════════════════════════════════════════════════════════════")
	fmt.Printf("%-40s %-10s %-12s %-12s %-15s\n", "Problem", "Difficulty", "EF Factor", "Repetitions", "Next Review")
	fmt.Println("───────────────────────────────────────────────────────────────────────────────────────")

	for _, r := range reviews {
		if !r.EasinessFactor.Valid {
			continue
		}

		ef := r.EasinessFactor.Float64
		totalEF += ef
		countEF++

		if ef < minEF {
			minEF = ef
		}
		if ef > maxEF {
			maxEF = ef
		}

		// Categorize
		if ef < 2.0 {
			lowEF++
		} else if ef < 2.5 {
			mediumEF++
		} else if ef <= 3.0 {
			highEF++
		} else {
			veryHighEF++
		}

		// Format EF with color coding
		efColor := ""
		efReset := "\033[0m"
		if ef < 2.0 {
			efColor = "\033[31m" // Red - too low
		} else if ef < 2.5 {
			efColor = "\033[33m" // Yellow - moderate
		} else if ef <= 3.0 {
			efColor = "\033[32m" // Green - good
		} else {
			efColor = "\033[35m" // Magenta - very high (might be too fast)
		}

		reps := "N/A"
		if r.Repetitions.Valid {
			reps = fmt.Sprintf("%d", r.Repetitions.Int64)
		}

		nextReview := "N/A"
		if r.NextReviewDate.Valid {
			nextReview = formatReviewDate(r.NextReviewDate)
		}

		fmt.Printf("%-40s %-10s %s%5.2f%s      %-12s %-15s\n",
			truncate(r.Title, 40),
			r.Difficulty,
			efColor,
			ef,
			efReset,
			reps,
			nextReview,
		)
	}

	fmt.Println()

	// Show statistics
	if countEF > 0 {
		avgEF := totalEF / float64(countEF)
		fmt.Println("📊 Statistics:")
		fmt.Println("───────────────────────────────────────────────────────────────────────────────────────")
		fmt.Printf("  Average EF:  %.2f\n", avgEF)
		fmt.Printf("  Minimum EF:  %.2f\n", minEF)
		fmt.Printf("  Maximum EF:  %.2f\n", maxEF)
		fmt.Println()
		fmt.Println("📊 Distribution:")
		fmt.Println("───────────────────────────────────────────────────────────────────────────────────────")
		fmt.Printf("  \033[31mLow (< 2.0):\033[0m      %d problems (intervals grow slowly)\n", lowEF)
		fmt.Printf("  \033[33mMedium (2.0-2.5):\033[0m  %d problems (moderate growth)\n", mediumEF)
		fmt.Printf("  \033[32mHigh (2.5-3.0):\033[0m    %d problems (good growth rate)\n", highEF)
		fmt.Printf("  \033[35mVery High (> 3.0):\033[0m %d problems (intervals grow very fast)\n", veryHighEF)
		fmt.Println()
		fmt.Println("💡 Understanding Easiness Factors:")
		fmt.Println("───────────────────────────────────────────────────────────────────────────────────────")
		fmt.Println("  • EF determines how quickly review intervals increase")
		fmt.Println("  • After 2 reviews, interval = previous_interval × EF")
		fmt.Println("  • Higher EF = faster interval growth = fewer reviews")
		fmt.Println("  • If intervals are growing too fast, your EF values may be too high")
		fmt.Println("  • Marking problems as 'Hard' (3) will lower EF")
		fmt.Println("  • Marking problems as 'Easy' (1) will raise EF")
		fmt.Println("  • Default starting EF is 2.5")
		fmt.Println()
	}

	return nil
}

func getCommands(db *sql.DB) map[string]CliCommand {
	commands := []CliCommand{
		{
			Name:        "help",
			Aliases:     []string{"h", "?"},
			Description: "Display command for utilizing CLI tool",
			Callback:    helpCommand,
		},
		{
			Name:        "exit",
			Aliases:     []string{"q", "quit"},
			Description: "Exit the application",
			Callback:    exitCommand,
		},
		{
			Name:        "study",
			Aliases:     []string{"s"},
			Description: "Get your daily questions to study",
			Callback: func(args []string) error {
				return studyCommandWithDB(db, args)
			},
		},
		{
			Name:        "review",
			Aliases:     []string{"r"},
			Description: "View your review history and upcoming reviews",
			Callback: func(args []string) error {
				return reviewCommandWithDB(db, args)
			},
		},
		{
			Name:        "stat",
			Aliases:     []string{"st"},
			Description: "View your overall study statistics",
			Callback: func(args []string) error {
				return statCommandWithDB(db, args)
			},
		},
		{
			Name:        "factors",
			Aliases:     []string{"f"},
			Description: "View your learning factors (easiness factors) for all problems",
			Callback: func(args []string) error {
				return factorsCommandWithDB(db, args)
			},
		},
	}

	m := make(map[string]CliCommand, len(commands))
	for _, cmd := range commands {
		m[cmd.Name] = cmd
		for _, alias := range cmd.Aliases {
			m[alias] = cmd
		}
	}
	return m
}
