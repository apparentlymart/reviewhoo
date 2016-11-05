package main

import (
	"path"
)

// Returns the number of commits where each author touched a file
// with each distinct file extension.
//
// This is not an exact science, but is intended to approximate
// an answer to finding e.g. the people who tend to write JavaScript in
// a codebase that is a mixture of JavaScript and CSS, or similar.
func authorsByFileType() (map[string]StringRanking, error) {
	tally := make(map[string]map[string]int)

	lines, err := Git("log", "--no-merges", "--name-only", "--pretty=format:%aN", "--since=2 months ago")
	if err != nil {
		return nil, err
	}

	currentAuthor := ""
	for _, line := range lines {
		if line == "" {
			// There's a blank line between each commit, which is then
			// followed by the author email from our --pretty=format: above
			currentAuthor = ""
			continue
		}

		// If we don't currently have a current author then this is either
		// the very first line or it's the line after the blank line separating
		// commits, and in either case it must be an author email.
		if currentAuthor == "" {
			currentAuthor = line
			continue
		}

		// Otherwise it's a path to a file that has changed.
		ext := path.Ext(line)
		if ext == "" {
			// Can't do anything with files that don't have extensions
			continue
		}

		if _, ok := tally[ext]; !ok {
			tally[ext] = make(map[string]int)
		}
		if _, ok := tally[ext][currentAuthor]; !ok {
			tally[ext][currentAuthor] = 0
		}
		tally[ext][currentAuthor]++
	}

	ret := make(map[string]StringRanking)
	for ext, authorMap := range tally {
		ret[ext] = make(StringRanking, 0, len(authorMap))

		for author, lines := range authorMap {
			score := &StringScore{
				String: author,
				Score:  lines,
			}
			ret[ext] = append(ret[ext], score)
		}

		ret[ext].Sort()
	}

	return ret, nil
}
