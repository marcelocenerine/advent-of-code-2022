// Solution to https://adventofcode.com/2022/day/2
package adventofcode

import (
	"fmt"
	"regexp"
	"strconv"
)

type RockPaperScissors struct{}

func (r RockPaperScissors) Details() Details {
	return Details{Day: 2, Description: "Rock Paper Scissors"}
}

func (r RockPaperScissors) Solve(input *Input) (Result, error) {
	rounds, err := parseRounds(input)
	if err != nil {
		return Result{}, err
	}
	part1 := calculatePlayerTotalScore(rounds, strategy1())
	part2 := calculatePlayerTotalScore(rounds, strategy2())
	return Result{
		Part1: strconv.Itoa(part1),
		Part2: strconv.Itoa(part2),
	}, nil
}

type Shape string

const (
	Rock     Shape = "A"
	Paper    Shape = "B"
	Scissors Shape = "C"
)

func (p Shape) scoreAgainst(other Shape) int {
	if p == other { // draw
		return 3 + shapePoints[p]
	}
	if loserFor[p] == other { // win
		return 6 + shapePoints[p]
	}
	return shapePoints[p] // loss
}

var shapePoints = map[Shape]int{
	Rock:     1,
	Paper:    2,
	Scissors: 3,
}
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

type EncPlay string

type Round struct {
	opponent Shape
	player   EncPlay
}

type Strategy func(Round) Shape

func strategy1() Strategy {
	mapping := map[EncPlay]Shape{
		"X": Rock,
		"Y": Paper,
		"Z": Scissors,
	}
	return func(r Round) Shape { return mapping[r.player] }
}

func strategy2() Strategy {
	return func(r Round) Shape {
		switch r.player {
		case "X":
			return loserFor[r.opponent]
		case "Z":
			return winnerFor[r.opponent]
		}
		return r.opponent // Y = draw
	}
}

func calculatePlayerTotalScore(rounds []Round, strategy Strategy) int {
	var totalScore int
	for _, round := range rounds {
		player := strategy(round)
		totalScore += player.scoreAgainst(round.opponent)
	}
	return totalScore
}

func parseRounds(input *Input) ([]Round, error) {
	rgx := regexp.MustCompile(`^([ABC]) ([XYZ])$`)
	lines := input.Lines()
	rounds := make([]Round, len(lines))
	for i, line := range input.Lines() {
		round, err := parseRound(line, rgx)
		if err != nil {
			return nil, err
		}
		rounds[i] = round
	}
	return rounds, nil
}

func parseRound(line string, rgx *regexp.Regexp) (Round, error) {
	if !rgx.MatchString(line) {
		return Round{}, fmt.Errorf("invalid line: %s", line)
	}
	groups := rgx.FindAllStringSubmatch(line, -1)
	opponent := Shape(groups[0][1])
	player := EncPlay(groups[0][2])
	return Round{player: player, opponent: opponent}, nil
}
