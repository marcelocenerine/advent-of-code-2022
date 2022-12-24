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
	for _, crane := range []crane{cm9000{}, cm9001{}} {
		stacks, arrangements, err := parseInputFile("input.txt")
		if err != nil {
			fmt.Printf("Invalid input: %v\n", err)
			os.Exit(1)
		}
		rearranged, err := rearrange(stacks, arrangements, crane)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Result: ", topCrates(rearranged))
	}
}

func topCrates(stacks []*stack) string {
	var buffer bytes.Buffer
	for _, s := range stacks {
		if top, err := s.peek(); err == nil {
			buffer.WriteRune(rune(top))
		}
	}
	return buffer.String()
}

type crate rune

func (c crate) String() string {
	return string(c)
}

type crane interface {
	move(n int, from, to *stack) error
}

type cm9000 struct{}

func (cm cm9000) move(n int, from, to *stack) error {
	if from.size() < n {
		return fmt.Errorf("from.size < n: size=%d, n=%d", from.size(), n)
	}
	for i := 0; i < n; i++ {
		top, _ := from.pop()
		to.push(top)
	}
	return nil
}

type cm9001 struct{}

func (cm cm9001) move(n int, from, to *stack) error {
	if from.size() < n {
		return fmt.Errorf("from.size < n: size=%d, n=%d", from.size(), n)
	}
	crates := from.crates[from.size()-n:]
	to.crates = append(to.crates, crates...)
	from.crates = from.crates[:from.size()-n]
	return nil
}

type stack struct {
	id     int
	crates []crate
}

func (s *stack) push(c crate) {
	s.crates = append(s.crates, c)
}

func (s *stack) peek() (crate, error) {
	if s.size() == 0 {
		return 0, fmt.Errorf("stack %d is empty", s.id)
	}
	return s.crates[s.size()-1], nil
}

func (s *stack) pop() (crate, error) {
	top, err := s.peek()
	if err != nil {
		return 0, err
	}
	s.crates = s.crates[:s.size()-1]
	return top, nil
}

func (s *stack) size() int {
	return len(s.crates)
}

func (c *stack) String() string {
	return fmt.Sprintf("%d - %v", c.id, c.crates)
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
	result := make([]step, len(lines))

	for i, line := range lines {
		if !rgx.MatchString(line) {
			return result, fmt.Errorf("invalid step line: %s", line)
		}
		groups := rgx.FindAllStringSubmatch(line, -1)
		n, _ := strconv.Atoi(groups[0][1])
		from, _ := strconv.Atoi(groups[0][2])
		to, _ := strconv.Atoi(groups[0][3])
		result[i] = step{n: n, from: from, to: to}
	}
	return result, nil
}

func rearrange(stacks []*stack, arrangement []step, cm crane) ([]*stack, error) {
	for _, step := range arrangement {
		if step.from > len(stacks) {
			return nil, fmt.Errorf("invalid 'from' in step: %v", step)
		}
		if step.to > len(stacks) {
			return nil, fmt.Errorf("invalid 'to' in step: %v", step)
		}
		from := stacks[step.from-1]
		to := stacks[step.to-1]
		err := cm.move(step.n, from, to)
		if err != nil {
			return nil, err
		}
	}
	return stacks, nil
}
