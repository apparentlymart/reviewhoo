package main

import (
	"path"
	"strconv"
	"strings"
)

func fileChanges(base, target string) (StringRanking, error) {
	lines, err := Git("diff", "--numstat", base, target)
	if err != nil {
		return nil, err
	}

	changes := make(StringRanking, 0, len(lines))

	for _, line := range lines {
		fields := strings.Fields(line)
		added, _ := strconv.Atoi(fields[0])
		removed, _ := strconv.Atoi(fields[1])
		filename := fields[2]

		// Removed lines get twice the weight of added ones, both because
		// removing is often more destructive and because most diffs
		// (unfortunately) have a lot more adding than removing.
		score := added + (2 * removed)

		changes.Add(filename, score)
	}

	changes.Sort()

	return changes, nil
}

// Takes a ranking of file statistics and turns it into a ranking of
// file extension statistics.
func fileExtStats(fileStats StringRanking) StringRanking {
	tally := map[string]int{}

	for _, score := range fileStats {
		ext := path.Ext(score.String)
		if ext == "" {
			continue
		}

		tally[ext] = tally[ext] + score.Score
	}

	return NewRanking(tally)
}

func dirChanges(base, target string) (StringRanking, error) {
	lines, err := Git("diff", "--dirstat=lines", base, target)
	if err != nil {
		return nil, err
	}

	changes := make(StringRanking, 0, len(lines))

	for _, line := range lines {
		fields := strings.Fields(line)
		percent, _ := strconv.Atoi(fields[0][:len(fields[0])-3])
		filename := fields[1]

		changes.Add(filename, percent)
	}

	changes.Sort()

	return changes, nil
}
