package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

var jsonOutput = flag.Bool("json", false, "Produce machine-readable JSON output")

func main() {
	flag.Parse()
	rand.Seed(time.Now().Unix())

	referenceCommit := flag.Arg(0)
	newCommit := flag.Arg(1)

	if referenceCommit == "" {
		referenceCommit = "master"
	}
	if newCommit == "" {
		newCommit = "HEAD"
	}

	// Need to find the common merge base of the two given commits,
	// so that we will use the diff from *where the branch diverged*,
	// rather than where the reference *currently* is. This is particularly
	// important when dealing with branches that diverged a long time
	// ago, where the reference branch may have changed considerably
	// in the mean time.
	baseCommit, err := getMergeBase(referenceCommit, newCommit)
	if err != nil {
		log.Fatalf("Failed to find common merge base between %s and %s: %s", referenceCommit, newCommit, err)
	}

	localUserName, err := getLocalUserName()
	if err != nil {
		log.Fatalf("Failed to determine your git full name: %s", err)
	}

	filesChanged, err := fileChanges(baseCommit, newCommit)
	if err != nil {
		log.Fatalf("Failed to obtain file change stats from %s to %s: %s", baseCommit, newCommit, err)
	}
	topFilesChanged := filesChanged.TopPercent(50, 0).Sorted()
	fileChangers, err := authorsForPaths(baseCommit, topFilesChanged)
	topFileChangers := fileChangers.TopPercent(50, 0)
	topFileChangers.Remove(localUserName)
	if err != nil {
		panic(err)
	}

	dirsChanged, err := dirChanges(baseCommit, newCommit)
	if err != nil {
		log.Fatalf("Failed to obtain directory change stats from %s to %s: %s", baseCommit, newCommit, err)
	}
	topDirsChanged := dirsChanged.TopPercent(50, 0).Sorted()
	dirChangers, err := authorsForPaths(baseCommit, topDirsChanged)
	topDirChangers := dirChangers.TopPercent(75, 2)
	topDirChangers.Remove(localUserName)
	if err != nil {
		panic(err)
	}

	fileTypeStats, err := authorsByFileType()
	if err != nil {
		log.Fatalf("Failed to look up file type modification statistics: %s", err)
	}
	extsChanged := fileExtStats(filesChanged).TopPercent(50, 0).Sorted()
	extChangers := map[string]StringSet{}

	for _, ext := range extsChanged {
		if fileTypeStats[ext] == nil {
			continue
		}
		topChangers := fileTypeStats[ext].TopPercent(60, 3)
		topChangers.Remove(localUserName)

		extChangers[ext] = topChangers
	}

	reviewers := make([]*Reviewer, 0, 1+len(extChangers))
	namesUsed := make(StringSet)
	namesUsed.Add(localUserName)

	topFileChangers = topFileChangers.Subtract(namesUsed)
	if len(topFileChangers) > 0 {
		reviewer := topFileChangers.Shuffled()[0]
		reviewers = append(reviewers, &Reviewer{
			Name:   reviewer,
			Reason: "has changed these files recently",
		})
		namesUsed.Add(reviewer)
	} else {
		topDirChangers = topDirChangers.Subtract(namesUsed)
		if len(topDirChangers) > 0 {
			reviewer := topDirChangers.Shuffled()[0]
			reviewers = append(reviewers, &Reviewer{
				Name:   reviewer,
				Reason: "has changed files in these directories recently",
			})
			namesUsed.Add(reviewer)
		}
	}

	for ext, changers := range extChangers {
		changers = changers.Subtract(namesUsed)

		// If possible we will try to use people who haven't worked
		// in similar files or dirs recently, to get the "knows this
		// language but might not know this specific subsystem" reviewer
		// that can point out when things might be confusing to new
		// maintainers.
		constrainedChangers := changers.Subtract(topFileChangers).Subtract(topDirChangers)
		if len(constrainedChangers) > 0 {
			changers = constrainedChangers
		}

		if len(changers) > 0 {
			reviewer := changers.Shuffled()[0]
			reviewers = append(reviewers, &Reviewer{
				Name:   reviewer,
				Reason: fmt.Sprintf("has made lots of %s changes recently", ext),
			})
			namesUsed.Add(reviewer)
		}
	}

	if *jsonOutput {
		printJSONResult(reviewers)
	} else {
		printMarkdownResult(reviewers)
	}
}

func printJSONResult(reviewers []*Reviewer) {
	buf, err := json.MarshalIndent(reviewers, "", "    ")
	if err != nil {
		// Should never happen
		panic(err)
	}
	os.Stdout.Write(buf)
	os.Stdout.Write([]byte{10})
}

func printMarkdownResult(reviewers []*Reviewer) {
	if len(reviewers) == 0 {
		fmt.Println("\nreviewhoo isn't sure who to suggest for this one. Sorry!\n")
		return
	}

	fmt.Println("\nreviewhoo suggests the following reviewers:\n")
	for _, reviewer := range reviewers {
		fmt.Printf("* **%s** %s\n", reviewer.Name, reviewer.Reason)
	}
	fmt.Println("")
}

type Reviewer struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}
