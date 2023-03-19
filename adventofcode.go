package adventofcode

import (
	"fmt"
	"os"
	"strings"
)

func LoadInput(p Puzzle) (Input, error) {
	path := fmt.Sprintf("inputs/d%02d.txt", p.Details().Day)
	if bytes, err := os.ReadFile(path); err != nil {
		return "", err
	} else {
		return Input(bytes[:]), nil
	}
}

type Details struct {
	Day         int
	Description string
}

func (d Details) String() string {
	return fmt.Sprintf("Day %02d: %s", d.Day, d.Description)
}

type Input string

func (i Input) Lines() []string {
	return strings.Split(string(i), "\n")
}

type Result struct {
	Part1, Part2 string
}

type Puzzle interface {
	Details() Details
	Solve(*Input) (Result, error)
}
