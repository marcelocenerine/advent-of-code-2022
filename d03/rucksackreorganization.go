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
		part1 := part1SumOfPriorities(rucksacks)
		fmt.Printf("Part I: %d\n", part1)

		if part2, err := part2SumOfPriorities(rucksacks); err != nil {
			fmt.Printf("Invalid input: %v\n", err)
		} else {
			fmt.Printf("Part II: %d\n", part2)
		}
	}
}

func part1SumOfPriorities(rucksacks []rucksack) int {
	sumOfPriorities := 0
	for _, rs := range rucksacks {
		// fmt.Println(rs)
		sumOfPriorities += rs.common.priority()
	}
	return sumOfPriorities
}

func part2SumOfPriorities(rucksacks []rucksack) (int, error) {
	groups, err := mkGroups(rucksacks)
	if err != nil {
		return 0, err
	}

	sumOfPriorities := 0
	for _, gr := range groups {
		// fmt.Println(gr)
		sumOfPriorities += gr.badge.priority()
	}
	return sumOfPriorities, nil
}

const groupSize = 3

type item rune

func (it item) priority() int {
	if 'a' <= it && it <= 'z' {
		return 1 + int(it-'a')
	}
	return 27 + int(it-'A')
}

type compartment map[item]bool

type rucksack struct {
	items  []item
	left   compartment
	right  compartment
	common item
}

type group struct {
	elfs  []rucksack
	badge item
}

func (gr group) String() string {
	return fmt.Sprintf("elfs=%v, badge=%v, priority=%d", gr.elfs, string(gr.badge), gr.badge.priority())
}

func (rs rucksack) contains(it item) bool {
	return rs.left[it] || rs.right[it]
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
	lines := strings.Split(input, "\n")
	rucksacks := make([]rucksack, 0, len(lines))

	for _, line := range lines {
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

	rs.left = mkCompartment(line[:rcount/2])
	rs.right = mkCompartment(line[rcount/2:])

	for lr, _ := range rs.left {
		if _, ok := rs.right[lr]; ok {
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

func mkCompartment(line string) compartment {
	res := compartment{}
	for _, run := range line {
		res[item(run)] = true
	}
	return res
}

func mkGroups(rucksacks []rucksack) ([]group, error) {
	length := len(rucksacks)
	if length%groupSize != 0 {
		return nil, fmt.Errorf("%d rucksacks cannot be put into groups of %d", length, groupSize)
	}
	groups := make([]group, 0, length/groupSize)
gloop:
	for i := 0; i < length; i += groupSize {
		elfs := rucksacks[i : i+groupSize]
		first, others := elfs[0], elfs[1:]
	iloop:
		for _, it := range first.items {
			for _, other := range others {
				if !other.contains(it) {
					continue iloop
				}
			}
			groups = append(groups, group{elfs: elfs, badge: it})
			continue gloop
		}

		return nil, fmt.Errorf("group doesn't have a badge: %v", elfs)
	}
	return groups, nil
}
