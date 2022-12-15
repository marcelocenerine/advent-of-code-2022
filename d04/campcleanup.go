// Solution to https://adventofcode.com/2022/day/4
package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	if assignments, err := parseInputFile("input.txt"); err != nil {
		fmt.Printf("Invalid input: %v\n", err)
	} else {
		fmt.Println("How many assignment pairs does one range fully contain the other? ", part1CountFullOverlaps(assignments))
		fmt.Println("How many assignment pairs does one range overlaps the other? ", part2CountOverlaps(assignments))
	}
}

func part1CountFullOverlaps(assignments []assignment) int {
	count := 0
	for _, as := range assignments {
		if as.left.fullyOverlaps(as.right) {
			count++
		}
	}
	return count
}

func part2CountOverlaps(assignments []assignment) int {
	count := 0
	for _, as := range assignments {
		if as.left.overlaps(as.right) {
			count++
		}
	}
	return count
}

type assignment struct {
	left, right section
}

func (a assignment) String() string {
	return fmt.Sprintf("[%s,%s]", a.left, a.right)
}

type section struct {
	start, end int
}

func (s section) overlaps(o section) bool {
	return (s.start >= o.start && s.start <= o.end) || (s.end >= o.start && s.end <= o.end) || s.fullyOverlaps(o)
}

func (s section) fullyOverlaps(o section) bool {
	return (s.start <= o.start && s.end >= o.end) || (s.start >= o.start && s.end <= o.end)
}

func (s section) String() string {
	return fmt.Sprintf("%d-%d", s.start, s.end)
}

func parseInputFile(path string) ([]assignment, error) {
	if bytes, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		return parseInput(string(bytes[:]))
	}
}

func parseInput(input string) ([]assignment, error) {
	rgx, err := regexp.Compile(`^(\d+)-(\d+),(\d+)-(\d+)$`)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(input, "\n")
	assignments := make([]assignment, 0, len(lines))

	for _, line := range lines {
		as, err := parseLine(line, rgx)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, as)
	}
	return assignments, nil
}

func parseLine(line string, rgx *regexp.Regexp) (assignment, error) {
	as := assignment{}
	if !rgx.MatchString(line) {
		return as, fmt.Errorf("invalid line: %s", line)
	}

	groups := rgx.FindAllStringSubmatch(line, -1)
	ls, _ := strconv.Atoi(groups[0][1])
	le, _ := strconv.Atoi(groups[0][2])
	rs, _ := strconv.Atoi(groups[0][3])
	re, _ := strconv.Atoi(groups[0][4])
	return assignment{
		left:  section{start: ls, end: le},
		right: section{start: rs, end: re},
	}, nil
}
