# GoStudyNeetCode

A CLI-based spaced repetition study tool for mastering the NeetCode 150 problems. Built with Go and SQLite.

## Features

- **Spaced Repetition Algorithm**: Intelligently schedules problem reviews based on your performance
- **Interactive CLI**: Clean REPL interface for studying and tracking progress
- **SQLite Database**: Persistent storage of your study progress
- **NeetCode 150 Integration**: Pre-configured with the popular NeetCode 150 problem set

## Installation

```bash
git clone https://github.com/yourusername/GoStudyNeetCode.git
cd GoStudyNeetCode
go build
```

## Usage

Run the application:
```bash
./GoStudyNeetCode
```

Available commands within the REPL:
- Type `help` to see all available commands
- Study new problems and review existing ones
- Track your progress through the NeetCode 150

## Setup

### Seeding NeetCode 150

On first run, the app will attempt to seed the `problems` table with the NeetCode 150 list from `neetcode_150.json` in the project root.

Expected JSON format:
```json
[
  { "title": "Two Sum", "difficulty": "Easy" },
  { "title": "Valid Anagram", "difficulty": "Easy" },
  { "title": "Group Anagrams", "difficulty": "Medium" }
]
```

If the file is missing, seeding is skipped and you can add it later and rerun with a fresh database.

## How It Works

This tool uses a spaced repetition algorithm to help you retain problem-solving patterns. When you review problems:
- Correct answers increase the review interval
- Incorrect answers reset the problem for more frequent review
- The algorithm adapts to your learning pace

## Technologies

- **Go**: Core application logic
- **SQLite**: Lightweight database for progress tracking
- **go-sqlite3**: SQLite driver for Go

## License

MIT

## Contributing

Feel free to open issues or submit pull requests!
