package main

type CliCommand struct {
	Name        string
	Description string
	Callback    func(args []string) error
}

type Problem struct {
	Title          string `json:"title"`
	Difficulty     string `json:"difficulty"`
	Grouping       string `json:"grouping"`
	LeetcodeNumber int    `json:"leetcode_number"`
}
