// Solution to https://adventofcode.com/2022/day/3
package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

func main() {
	if rucksacks, err := parseInputFile("input.txt"); err != nil {
		fmt.Printf("Invalid input: %v\n", err)
	} else {
		sumOfPriorities := 0
		for _, rs := range rucksacks {
			fmt.Println(rs)
			sumOfPriorities += rs.common.priority()
		}
		fmt.Printf("Sum of priorities of common items: %d\n", sumOfPriorities)
	}
}

type item rune

func (it item) priority() int {
	if 'a' <= it && it <= 'z' {
		return 1 + int(it-'a')
	}
	return 27 + int(it-'A')
}

type rucksack struct {
	items  []item
	common item
}

func (rs rucksack) String() string {
	return fmt.Sprintf("items=%s, common=%v, priority=%d", string(rs.items), string(rs.common), rs.common.priority())
}

func parseInputFile(path string) ([]rucksack, error) {
	if bytes, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		return parseInput(string(bytes[:]))
	}
}

func parseInput(input string) ([]rucksack, error) {
	rgx, err := regexp.Compile("[a-zA-Z]")
	if err != nil {
		return nil, err
	}
	var rucksacks []rucksack

	for _, line := range strings.Split(input, "\n") {
		rs, err := parseLine(line, rgx)
		if err != nil {
			return nil, err
		}
		rucksacks = append(rucksacks, rs)
	}
	return rucksacks, nil
}

func parseLine(line string, rgx *regexp.Regexp) (rucksack, error) {
	rs := rucksack{items: []item(line)}
	rcount := utf8.RuneCountInString(line)
	if rcount == 0 || rcount%2 != 0 {
		return rs, fmt.Errorf("rucksack is empty or contains an odd number of items: %s", line)
	}
	if !rgx.MatchString(line) {
		return rs, fmt.Errorf("rucksack contains invalid items: %s", line)
	}

	left := charMap(line[:rcount/2])
	right := charMap(line[rcount/2:])

	for lr, _ := range left {
		if _, ok := right[lr]; ok {
			if rs.common > 0 {
				return rs, fmt.Errorf("rucksack contains multiple common items in the two compartments: %s", line)
			}
			rs.common = item(lr)
		}
	}

	if rs.common == 0 {
		return rs, fmt.Errorf("rucksack contains no common item in the two compartments: %s", line)
	}
	return rs, nil
}

func charMap(line string) map[rune]bool {
	res := make(map[rune]bool)
	for _, run := range line {
		res[run] = true
	}
	return res
}
