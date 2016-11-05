package main

import (
	"bufio"
	"fmt"
	"os/exec"
)

func Git(args ...string) ([]string, error) {
	cmd := exec.Command("git", args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	ret := make([]string, 0, 20)
	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		ret = append(ret, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	err = cmd.Wait()
	return ret, err
}

func getLocalUserName() (string, error) {
	lines, err := Git("config", "user.name")
	if err != nil {
		return "", err
	}
	if len(lines) == 0 {
		return "", fmt.Errorf("failed to determine local user name")
	}
	return lines[0], nil
}

func getMergeBase(base, target string) (string, error) {
	lines, err := Git("merge-base", base, target)
	if err != nil {
		return "", err
	}
	if len(lines) == 0 {
		return "", fmt.Errorf("failed to determine merge base for %s and %s", base, target)
	}
	return lines[0], nil
}
