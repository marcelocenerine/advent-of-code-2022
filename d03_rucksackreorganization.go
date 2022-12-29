// Solution to https://adventofcode.com/2022/day/3
package adventofcode

import (
	"fmt"
	"regexp"
	"strconv"
	"unicode/utf8"
)

type RucksackReorganization struct{}

func (r RucksackReorganization) Details() Details {
	return Details{Day: 3, Description: "Rucksack Reorganization"}
}

func (r RucksackReorganization) Solve(input *Input) (Result, error) {
	rucksacks, err := parseRucksacks(input)
	if err != nil {
		return Result{}, err
	}
	part1 := part1SumOfPriorities(rucksacks)
	part2, err := part2SumOfPriorities(rucksacks)
	if err != nil {
		return Result{}, err
	}
	return Result{
		Part1: strconv.Itoa(part1),
		Part2: strconv.Itoa(part2),
	}, nil
}

func part1SumOfPriorities(rucksacks []rucksack) int {
	sumOfPriorities := 0
	for _, rs := range rucksacks {
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

func (rs rucksack) contains(it item) bool {
	return rs.left[it] || rs.right[it]
}

func parseRucksacks(input *Input) ([]rucksack, error) {
	rgx := regexp.MustCompile("[a-zA-Z]")
	lines := input.Lines()
	rucksacks := make([]rucksack, 0, len(lines))

	for _, line := range lines {
		rs, err := parseRucksack(line, rgx)
		if err != nil {
			return nil, err
		}
		rucksacks = append(rucksacks, rs)
	}
	return rucksacks, nil
}

func parseRucksack(line string, rgx *regexp.Regexp) (rucksack, error) {
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
