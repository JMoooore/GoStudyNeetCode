# GoStudyNeetCode
## Seeding NeetCode 150

On first run, the app will attempt to seed the `problems` table with the NeetCode 150 list from `neetcode_150.json` in the project root.

- Expected JSON format:

```json
[
  { "title": "Two Sum", "difficulty": "Easy" },
  { "title": "Valid Anagram", "difficulty": "Easy" },
  { "title": "Group Anagrams", "difficulty": "Medium" }
]
```

- If the file is missing, seeding is skipped and you can add it later and rerun with a fresh database.

## Description
- The purpose of this app is to assist in NeetCode utilizing spaced repetion algorithms.
## CLI Flags
- review -> Amount of review cards to study
- new -> Amount of new cards to study

