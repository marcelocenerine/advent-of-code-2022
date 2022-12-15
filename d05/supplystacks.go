// Solution to https://adventofcode.com/2022/day/5
package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	stacks, arrangements, err := parseInputFile("input.txt")
	if err != nil {
		fmt.Printf("Invalid input: %v\n", err)
		os.Exit(1)
	}
	err = rearrange(stacks, arrangements)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Cranes at the top of the stacks: %s\n", part1CranesOnTop(stacks))
}

func part1CranesOnTop(stacks []*stack) string {
	var buffer bytes.Buffer
	for _, stk := range stacks {
		if top, err := stk.peek(); err == nil {
			buffer.WriteRune(rune(top))
		}
	}
	return buffer.String()
}

type crate rune

func (c crate) String() string {
	return fmt.Sprintf("[%s]", string(c))
}

type stack struct {
	id     int
	crates []crate
}

func (s *stack) push(c crate) {
	s.crates = append(s.crates, c)
}

func (s *stack) pop() (crate, error) {
	top, err := s.peek()
	if err != nil {
		return 0, err
	}
	s.crates = s.crates[:s.size()-1]
	return top, nil
}

func (s *stack) peek() (crate, error) {
	if s.size() == 0 {
		return 0, fmt.Errorf("stack %d is empty", s.id)
	}
	return s.crates[s.size()-1], nil
}

func (s *stack) size() int {
	return len(s.crates)
}

func (c *stack) String() string {
	return fmt.Sprintf("{id=%d, crates=%v}", c.id, c.crates)
}

type step struct {
	n, from, to int
}

func (s step) String() string {
	return fmt.Sprintf("{n=%d, from=%d to=%d}", s.n, s.from, s.to)
}

func parseInputFile(path string) ([]*stack, []step, error) {
	if bytes, err := os.ReadFile(path); err != nil {
		return nil, nil, err
	} else {
		return parseInput(string(bytes[:]))
	}
}

func parseInput(input string) ([]*stack, []step, error) {
	rgx, err := regexp.Compile(`^move (\d+) from (\d+) to (\d+)$`)
	if err != nil {
		return nil, nil, err
	}
	lines := strings.Split(input, "\n")
	if len(lines) == 0 {
		return nil, nil, nil
	}

	for i := 0; i < len(lines); i++ {
		if lines[i] == "" {
			stks, err := parseStacks(lines[:i])
			if err != nil {
				return nil, nil, err
			}
			arrs, err := parseArrangement(lines[i+1:], rgx)
			if err != nil {
				return nil, nil, err
			}
			return stks, arrs, nil
		}
	}

	return nil, nil, errors.New("empty line separating stacks and arrangement not found in the input")
}

func parseStacks(lines []string) ([]*stack, error) {
	if len(lines) < 2 {
		return nil, errors.New("invalid stacks")
	}
	length := len(lines[0])
	var result []*stack
	for i := len(lines) - 2; i >= 0; i-- { // discards the bottom line
		line := lines[i]
		if len(line) != length {
			return nil, errors.New("lines in the stack has inconsistent lengths")
		}
		for c, s := 1, 0; c < length; c, s = c+4, s+1 {
			if len(result) == s {
				result = append(result, &stack{id: s + 1})
			}
			if crt := crate(line[c]); crt != ' ' {
				result[s].push(crt)
			}
		}
	}
	return result, nil
}

func parseArrangement(lines []string, rgx *regexp.Regexp) ([]step, error) {
	result := make([]step, 0, len(lines))

	for _, line := range lines {
		if !rgx.MatchString(line) {
			return result, fmt.Errorf("invalid step line: %s", line)
		}
		groups := rgx.FindAllStringSubmatch(line, -1)
		n, _ := strconv.Atoi(groups[0][1])
		from, _ := strconv.Atoi(groups[0][2])
		to, _ := strconv.Atoi(groups[0][3])
		result = append(result, step{n: n, from: from, to: to})
	}
	return result, nil
}

func rearrange(stacks []*stack, arrangement []step) error {
	for _, step := range arrangement {
		if step.from > len(stacks) {
			return fmt.Errorf("invalid 'from' in step: %v", step)
		}
		if step.to > len(stacks) {
			return fmt.Errorf("invalid 'to' in step: %v", step)
		}
		from := stacks[step.from-1]
		to := stacks[step.to-1]
		if step.n > from.size() {
			return fmt.Errorf("invalid 'n' in step: %v", step)
		}
		for i := 0; i < step.n; i++ {
			top, _ := from.pop()
			to.push(top)
		}
	}
	return nil
}
