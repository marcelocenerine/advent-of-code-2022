// Solution to https://adventofcode.com/2022/day/11
package adventofcode

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type MonkeyInTheMiddle struct{}

func (p MonkeyInTheMiddle) Details() Details {
	return Details{Day: 11, Description: "Monkey in the Middle"}
}

func (p MonkeyInTheMiddle) Solve(input *Input) (Result, error) {
	monkeys, err := p.parseNotes(input)
	if err != nil {
		return Result{}, err
	}
	counts, err := p.processRounds(20, monkeys)
	if err != nil {
		return Result{}, err
	}
	monkeyBusinessLevel, err := p.calcMonkeyBusinessLevel(counts)
	if err != nil {
		return Result{}, err
	}
	return Result{
		Part1: strconv.Itoa(monkeyBusinessLevel),
		Part2: "",
	}, nil
}

type MonkeyId string
type WorryLevel int
type Operation func(wl WorryLevel) WorryLevel
type ThrowNext func(wl WorryLevel) MonkeyId

type Monkey struct {
	Id    MonkeyId
	Items []WorryLevel
	Op    Operation
	Test  ThrowNext
}

func (p MonkeyInTheMiddle) calcMonkeyBusinessLevel(inspectCounts map[MonkeyId]int) (int, error) {
	if len(inspectCounts) < 2 {
		return 0, fmt.Errorf("input contains %d monkeys; needs at least %d", len(inspectCounts), 2)
	}
	var counts []int
	for _, count := range inspectCounts {
		counts = append(counts, count)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(counts)))
	return counts[0] * counts[1], nil
}

func (p MonkeyInTheMiddle) processRounds(n int, monkeys []*Monkey) (map[MonkeyId]int, error) {
	monkeysById := map[MonkeyId]*Monkey{}
	inspectCountById := map[MonkeyId]int{}
	for _, monkey := range monkeys {
		if _, ok := monkeysById[monkey.Id]; ok {
			return nil, fmt.Errorf("duplicate monkey id: %s", monkey.Id)
		}
		monkeysById[monkey.Id] = monkey
	}

	for round := 0; round < n; round++ {
		for _, monkey := range monkeys {
			items := monkey.Items
			monkey.Items = nil
			for _, item := range items {
				newWl := WorryLevel(monkey.Op(item) / 3)
				next := monkey.Test(newWl)
				if nextMonkey, ok := monkeysById[next]; ok {
					nextMonkey.Items = append(nextMonkey.Items, newWl)
				} else {
					return nil, fmt.Errorf("invalid monkey id: %s", next)
				}
				inspectCountById[monkey.Id]++
			}
		}
	}

	return inspectCountById, nil
}

func (p MonkeyInTheMiddle) parseNotes(input *Input) ([]*Monkey, error) {
	monkeyRgx := regexp.MustCompile(`^Monkey (.+?):$`)
	itemsRgx := regexp.MustCompile(`^\s+Starting items: (.*)`)
	opRgx := regexp.MustCompile(`^\s+Operation: new = (old|\d+) ([+*]) (old|\d+)$`)
	testRgx := regexp.MustCompile(`^\s+Test: divisible by (\d+)$`)
	trueRgx := regexp.MustCompile(`^\s+If true: throw to monkey (.+)$`)
	falseRgx := regexp.MustCompile(`^\s+If false: throw to monkey (.+)$`)
	lines := input.Lines()
	var result []*Monkey

	for i := 0; i < len(lines); {
		if lines[i] == "" {
			i++
			continue
		}
		// id
		if !monkeyRgx.MatchString(lines[i]) {
			return nil, fmt.Errorf("invalid line: %s", lines[i])
		}
		monkeyLineGroups := monkeyRgx.FindAllStringSubmatch(lines[i], -1)
		id := MonkeyId(monkeyLineGroups[0][1])
		// items
		if !itemsRgx.MatchString(lines[i+1]) {
			return nil, fmt.Errorf("invalid line: %s", lines[i+1])
		}
		itemsLineGroups := itemsRgx.FindAllStringSubmatch(lines[i+1], -1)
		var items []WorryLevel
		for _, sitem := range strings.Split(itemsLineGroups[0][1], ", ") {
			item, err := strconv.Atoi(sitem)
			if err != nil {
				return nil, fmt.Errorf("invalid line: %s", lines[i])
			}
			items = append(items, WorryLevel(item))
		}
		// operation
		if !opRgx.MatchString(lines[i+2]) {
			return nil, fmt.Errorf("invalid line: %s", lines[i+2])
		}
		opLineGroups := opRgx.FindAllStringSubmatch(lines[i+2], -1)
		leftOperand := opLineGroups[0][1]
		operator := opLineGroups[0][2]
		rightOperand := opLineGroups[0][3]
		convOperand := func(operand string, wl WorryLevel) int {
			switch operand {
			case "old":
				return int(wl)
			default:
				n, _ := strconv.Atoi(operand)
				return n
			}
		}
		operation := func(wl WorryLevel) WorryLevel {
			lhs := convOperand(leftOperand, wl)
			rhs := convOperand(rightOperand, wl)
			switch operator {
			case "+":
				return WorryLevel(lhs + rhs)
			case "*":
				return WorryLevel(lhs * rhs)
			default:
				panic(fmt.Sprintf("unexpected operator: %s", operator))
			}
		}
		// test
		if !testRgx.MatchString(lines[i+3]) {
			return nil, fmt.Errorf("invalid line: %s", lines[i+3])
		}
		if !trueRgx.MatchString(lines[i+4]) {
			return nil, fmt.Errorf("invalid line: %s", lines[i+4])
		}
		if !falseRgx.MatchString(lines[i+5]) {
			return nil, fmt.Errorf("invalid line: %s", lines[i+5])
		}
		testLineGroups := testRgx.FindAllStringSubmatch(lines[i+3], -1)
		trueLineGroups := trueRgx.FindAllStringSubmatch(lines[i+4], -1)
		falseLineGroups := falseRgx.FindAllStringSubmatch(lines[i+5], -1)
		divisibleBy, _ := strconv.Atoi(testLineGroups[0][1])
		whenTrue := MonkeyId(trueLineGroups[0][1])
		whenFalse := MonkeyId(falseLineGroups[0][1])
		test := func(wl WorryLevel) MonkeyId {
			if int(wl)%divisibleBy == 0 {
				return whenTrue
			}
			return whenFalse
		}

		result = append(result, &Monkey{
			Id:    id,
			Items: items,
			Op:    operation,
			Test:  test,
		})
		i += 6
	}
	return result, nil
}
