// Solution to https://adventofcode.com/2022/day/2
package main

import (
	"fmt"
	"os"
	"strings"
)

type Shape string

const (
	Rock     Shape = "A"
	Paper    Shape = "B"
	Scissors Shape = "C"
)

var loserFor = map[Shape]Shape{
	Rock:     Scissors,
	Paper:    Rock,
	Scissors: Paper,
}
var winnerFor = map[Shape]Shape{
	Rock:     Paper,
	Paper:    Scissors,
	Scissors: Rock,
}
var shapeScore = map[Shape]int{
	Rock:     1,
	Paper:    2,
	Scissors: 3,
}

type Outcome string

const (
	Loss Outcome = "X"
	Draw Outcome = "Y"
	Win  Outcome = "Z"
)

var outcomeScore = map[Outcome]int{
	Loss: 0,
	Draw: 3,
	Win:  6,
}

type Round struct {
	player   Shape
	opponent Shape
	outcome  Outcome
}

func (r *Round) playerScore() int {
	return shapeScore[r.player] + outcomeScore[r.outcome]
}

func main() {
	if rounds, err := parseInputFile("input.txt"); err != nil {
		fmt.Printf("Invalid input: %v\n", err)
	} else {
		totalScore := calculatePlayerTotalScore(rounds)
		fmt.Printf("The player's total score is: %d\n", totalScore)
	}
}

func calculatePlayerTotalScore(rounds []Round) int {
	var totalScore int
	for _, round := range rounds {
		totalScore += round.playerScore()
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
	cols := strings.Split(line, " ")
	if len(cols) != 2 {
		return Round{}, fmt.Errorf("line must have exact 2 columns: %s", line)
	}
	opponent := Shape(cols[0])
	outcome := Outcome(cols[1])
	player := choosePlay(opponent, outcome)
	return Round{player: player, opponent: opponent, outcome: outcome}, nil
}

func choosePlay(opponent Shape, outcome Outcome) Shape {
	switch outcome {
	case Loss:
		return loserFor[opponent]
	case Win:
		return winnerFor[opponent]
	}
	return opponent // draw
}
