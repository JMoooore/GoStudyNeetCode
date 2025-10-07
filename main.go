package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func startRepl(db *sql.DB) {
	reader := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("GoStudy > ")
		reader.Scan()

		input := reader.Text()
		parts := strings.Fields(input)

		if len(parts) == 0 {
			continue
		}

		commandName := parts[0]
		args := parts[1:]

		command, exists := getCommands(db)[commandName]
		if exists {
			err := command.Callback(args)
			if err != nil {
				fmt.Println(err)
			}
			continue
		}

		fmt.Printf("Unknown command: %s. Type 'help' for available commands.\n", commandName)
	}
}

func main() {
	fmt.Println("=================================")
	fmt.Println("  Welcome to GoStudyNeetCode!    ")
	fmt.Println("=================================")
	fmt.Println()

	db, err := initDb()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	fmt.Println("Type 'help' to see available commands")
	startRepl(db)
}
