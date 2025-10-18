package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"os"
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

	for _, cmd := range getCommands(nil) {
		fmt.Printf("  %-10s - %s\n", cmd.Name, cmd.Description)
	}
	fmt.Println()
	fmt.Println("Example usage:")
	fmt.Println("  study --difficulty easy --count 5")
	fmt.Println("  study -d medium -c 3")
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
	var count int

	// Define flags
	fs.StringVar(&difficulty, "difficulty", "any", "Difficulty level (easy, medium, hard, any OR e, m, h, a)")
	fs.StringVar(&difficulty, "d", "any", "Short for difficulty")

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

	problems, err := selectStudyProblems(db, difficulty, count)
	if err != nil {
		return err
	}

	fmt.Println("\n📚 Your Study Problems:")
	fmt.Println("========================")
	for i, p := range problems {
		fmt.Printf("%d. [LC %d] %s (%s) - %s\n", i+1, p.LeetcodeNumber, p.Title, p.Difficulty, p.Grouping)
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
				fmt.Printf("\nHow hard was '%s'? (1=Easy, 2=Medium, 3=Hard): ", problem.Title)
				ratingStr, _ := reader.ReadString('\n')
				ratingStr = strings.TrimSpace(ratingStr)
				rating, err := strconv.Atoi(ratingStr)
				if err != nil || rating < 1 || rating > 3 {
					fmt.Println("Invalid rating, skipping...")
					continue
				}

				// Update the database
				if err := updateProblemCompletion(db, problem.Title, rating); err != nil {
					fmt.Printf("Error updating problem: %v\n", err)
				} else {
					fmt.Printf("\033[32m✓ Marked '%s' as completed with effort rating %d\033[0m\n", problem.Title, rating)

					// Remove from slice using 0-indexed position
					idx := num - 1
					problems = append(problems[:idx], problems[idx+1:]...)

					// Show updated list
					if len(problems) > 0 {
						fmt.Println("\nRemaining problems:")
						for i, p := range problems {
							fmt.Printf("%d. [LC %d] %s\n", i+1, p.LeetcodeNumber, p.Title)
						}
						fmt.Println()
					}
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

func getCommands(db *sql.DB) map[string]CliCommand {
	return map[string]CliCommand{
		"help": {
			Name:        "help",
			Description: "Display command for utilizing CLI tool",
			Callback:    helpCommand,
		},
		"exit": {
			Name:        "exit",
			Description: "Exit the application",
			Callback:    exitCommand,
		},
		"study": {
			Name:        "study",
			Description: "Get your daily questions to study",
			Callback: func(args []string) error {
				return studyCommandWithDB(db, args)
			},
		},
		"review": {
			Name:        "review",
			Description: "View your review history and upcoming reviews",
			Callback: func(args []string) error {
				return reviewCommandWithDB(db, args)
			},
		},
		"stat": {
			Name:        "stat",
			Description: "View your overall study statistics",
			Callback: func(args []string) error {
				return statCommandWithDB(db, args)
			},
		},
	}
}
