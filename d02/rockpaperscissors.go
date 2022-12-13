// Solution to https://adventofcode.com/2022/day/2
package main

import (
	"fmt"
	"os"
	"strings"
)

type Play int

const (
	Rock     Play = 1
	Paper    Play = 2
	Scissors Play = 3
)

type Outcome int

const (
	Win  Outcome = 6
	Draw Outcome = 3
	Lose Outcome = 0
)

func (p Play) outcomeAgainst(other Play) Outcome {
	switch p {
	case Rock:
		switch other {
		case Paper:
			return Lose
		case Scissors:
			return Win
		}
	case Paper:
		switch other {
		case Rock:
			return Win
		case Scissors:
			return Lose
		}
	case Scissors:
		switch other {
		case Rock:
			return Lose
		case Paper:
			return Win
		}
	}

	return Draw
}

type Round struct {
	player   Play
	opponent Play
}

func (r *Round) playersScore() int {
	return int(r.player) + int(r.player.outcomeAgainst(r.opponent))
}

func main() {
	if rounds, err := parseInputFile("input.txt"); err != nil {
		fmt.Printf("Invalid input: %v\n", err)
	} else {
		totalScore := calculatePlayersTotalScore(rounds)
		fmt.Printf("The player's total score is: %d\n", totalScore)
	}
}

func calculatePlayersTotalScore(rounds []Round) int {
	var totalScore int

	for _, round := range rounds {
		totalScore += round.playersScore()
	}

	return totalScore
}

func parseInputFile(path string) ([]Round, error) {
	if bytes, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		return parseInput(string(bytes[:]))
	}
}

func parseInput(input string) ([]Round, error) {
	var rounds []Round

	for _, line := range strings.Split(input, "\n") {
		round, err := parseLine(line)
		if err != nil {
			return nil, err
		}

		rounds = append(rounds, round)
	}

	return rounds, nil
}

func parseLine(line string) (Round, error) {
	var round Round
	cols := strings.Split(line, " ")

	if len(cols) != 2 {
		return round, fmt.Errorf("line must have exact 2 columns: %s", line)
	}

	opponent, err := str2Play(cols[0])
	if err != nil {
		return round, err
	}

	player, err := str2Play(cols[1])
	if err != nil {
		return round, err
	}

	round.player = player
	round.opponent = opponent
	return round, nil
}

func str2Play(str string) (Play, error) {
	switch str {
	case "A", "X":
		return Rock, nil
	case "B", "Y":
		return Paper, nil
	case "C", "Z":
		return Scissors, nil
	default:
		return 0, fmt.Errorf("invalid play: %s", str)
	}
}
