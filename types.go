package main

import "database/sql"

type CliCommand struct {
	Name        string
	Description string
	Callback    func(args []string) error
}

type Config struct {
	db          sql.DB
	nextProblem *string
	prevProblem *string
}

type Problem struct {
	Title          string `json:"title"`
	Difficulty     string `json:"difficulty"`
	Grouping       string `json:"grouping"`
	LeetcodeNumber int    `json:"leetcode_number"`
}
