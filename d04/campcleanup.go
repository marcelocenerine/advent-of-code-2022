// Solution to https://adventofcode.com/2022/day/3
package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	if assigments, err := parseInputFile("input.txt"); err != nil {
		fmt.Printf("Invalid input: %v\n", err)
	} else {
		count := 0
		for _, as := range assigments {
			if as.left.contains(as.right) || as.right.contains(as.left) {
				count++
			}
		}
		fmt.Println("How many assignment pairs does one range fully contain the other? ", count)
	}
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

func (s section) contains(that section) bool {
	return s.start <= that.start && s.end >= that.end
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
		return as, fmt.Errorf("rucksack contains invalid items: %s", line)
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
