// Solution to https://adventofcode.com/2022/day/4
package adventofcode

import (
	"fmt"
	"regexp"
	"strconv"
)

type CampCleanup struct{}

func (c CampCleanup) Details() Details {
	return Details{Day: 4, Description: "Camp Cleanup"}
}

func (c CampCleanup) Solve(input *Input) (Result, error) {
	assignments, err := parseAssignments(input)
	if err != nil {
		return Result{}, err
	}
	part1 := part1CountFullOverlaps(assignments)
	part2 := part2CountOverlaps(assignments)
	return Result{
		Part1: strconv.Itoa(part1),
		Part2: strconv.Itoa(part2),
	}, nil
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

type section struct {
	start, end int
}

func (s section) overlaps(o section) bool {
	return (s.start >= o.start && s.start <= o.end) || (s.end >= o.start && s.end <= o.end) || s.fullyOverlaps(o)
}

func (s section) fullyOverlaps(o section) bool {
	return (s.start <= o.start && s.end >= o.end) || (s.start >= o.start && s.end <= o.end)
}

func parseAssignments(input *Input) ([]assignment, error) {
	rgx := regexp.MustCompile(`^(\d+)-(\d+),(\d+)-(\d+)$`)
	lines := input.Lines()
	assignments := make([]assignment, len(lines))

	for i, line := range lines {
		as, err := parseAssignment(line, rgx)
		if err != nil {
			return nil, err
		}
		assignments[i] = as
	}
	return assignments, nil
}

func parseAssignment(line string, rgx *regexp.Regexp) (assignment, error) {
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
