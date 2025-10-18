# GoStudyNeetCode

[![Go Version](https://img.shields.io/badge/Go-1.16+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](http://makeapullrequest.com)

> **Master LeetCode & NeetCode 150 problems with spaced repetition** - A powerful CLI tool that helps you remember coding patterns and ace technical interviews using proven learning science.

A terminal-based spaced repetition system (SRS) for mastering coding interview problems. Built with Go and SQLite, this tool intelligently schedules LeetCode/NeetCode problem reviews to maximize retention and prepare you for FAANG interviews.

## ✨ Features

- **🧠 Spaced Repetition Algorithm**: Scientifically-proven learning technique that schedules coding problem reviews based on your performance - study smarter, not harder
- **⚡ Interactive CLI/REPL**: Clean terminal interface for focused studying without browser distractions
- **💾 SQLite Database**: Lightweight, persistent storage tracks your progress locally - no account required
- **📚 NeetCode 150 Integration**: Pre-configured with the complete NeetCode 150 problem set covering all essential interview patterns
- **🎯 Progress Tracking**: Monitor your mastery of data structures, algorithms, and problem-solving patterns
- **🚀 Offline-First**: Study anywhere without internet dependency

## 🎬 Demo

```
$ ./GoStudyNeetCode
Welcome to GoStudyNeetCode - Your Interview Prep Companion!
Type 'help' to see available commands.

> study
Loading your next problem...
[Add a screenshot or ASCII demo here]
```

<!-- TODO: Add demo GIF or screenshot -->

## 🚀 Quick Start

### Prerequisites
- Go 1.16 or higher
- Git

### Installation

```bash
# Clone the repository
git clone https://github.com/JMoooore/GoStudyNeetCode.git

cd GoStudyNeetCode

# Build the application
go build

# Run it!
./GoStudyNeetCode
```

### Alternative: Install with `go install`

```bash
go install github.com/JMoooore/GoStudyNeetCode@latest
```

### Optional: Create an Alias

For quicker access, add an alias to your shell configuration file (`~/.bashrc`, `~/.zshrc`, etc.):

```bash
# If you built locally
alias neetcode='/path/to/GoStudyNeetCode/GoStudyNeetCode'

# Or if you used go install (ensure $GOPATH/bin or $HOME/go/bin is in your PATH)
alias neetcode='GoStudyNeetCode'
```

After adding the alias, reload your shell configuration:
```bash
source ~/.zshrc  # or ~/.bashrc
```

Then you can simply run:
```bash
neetcode
```

## 📖 Usage

Launch the interactive study session:
```bash
./GoStudyNeetCode
```

### Available Commands
Once inside the REPL, you can:
- **`study`** - Start reviewing problems due for practice
- **`help`** - Display all available commands
- **`review`** - View your progress on individual problems
- **`stat`** - View your overall progress and statistics including completion estimate
- **`exit`** - Save and exit the application

The spaced repetition algorithm automatically determines which problems you should review based on your past performance.

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

## 🧪 How It Works

GoStudyNeetCode implements a **spaced repetition system (SRS)** - the same learning technique used by Anki, SuperMemo, and other proven study tools.

### The Algorithm
- ✅ **Correct answers** → Review interval increases exponentially (1 day → 3 days → 7 days → 14 days...)
- ❌ **Incorrect answers** → Problem resets for frequent review until mastered
- 🎯 **Adaptive scheduling** → The algorithm learns your pace and adjusts difficulty accordingly

This ensures you focus on weak areas while maintaining knowledge of mastered patterns - perfect for **interview preparation** and **long-term retention** of coding concepts.

## 🛠️ Built With

- **[Go](https://go.dev/)** - Fast, compiled performance with simple concurrency
- **[SQLite](https://www.sqlite.org/)** - Zero-configuration embedded database
- **[go-sqlite3](https://github.com/mattn/go-sqlite3)** - CGo-free SQLite driver

## 🗺️ Roadmap

- [ ] Web dashboard for statistics visualization
- [ ] Custom problem set support (import your own questions)
- [ ] Topic-based filtering (arrays, graphs, dynamic programming, etc.)
- [ ] Export/import progress for backup
- [ ] Multi-device sync

## 🤝 Contributing

Contributions are welcome! Whether it's:
- 🐛 Bug reports
- 💡 Feature requests
- 📝 Documentation improvements
- 🔧 Pull requests

Please feel free to open an issue or submit a PR.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ⭐ Support

If this tool helps you ace your interviews, consider giving it a star! It helps others discover the project.

## 🏷️ Keywords

`leetcode` `neetcode` `coding-interview` `interview-prep` `spaced-repetition` `cli` `golang` `sqlite` `srs` `flashcards` `study-tool` `algorithm` `data-structures` `faang` `technical-interview`

---

**Happy studying! May your interview prep be efficient and your offers be plentiful.** 🎯
