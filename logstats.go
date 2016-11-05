package main

func authorsForPaths(commit string, paths []string) (StringRanking, error) {
	args := []string{
		"log",
		"--no-merges",
		"--pretty=format:%aN",
		"--since=2 months ago",
		commit,
		"--",
	}
	args = append(args, paths...)
	lines, err := Git(args...)
	if err != nil {
		return nil, err
	}

	tally := map[string]int{}

	for _, line := range lines {
		tally[line]++
	}

	ranking := make(StringRanking, 0, len(tally))
	for name, score := range tally {
		ranking.Add(name, score)
	}

	ranking.Sort()

	return ranking, nil
}
