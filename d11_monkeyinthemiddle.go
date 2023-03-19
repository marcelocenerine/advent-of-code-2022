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
	part1, err := p.solve(20, p.divBy3Relief, input)
	if err != nil {
		return Result{}, err
	}
	part2, err := p.solve(10000, p.modByDivisorsRelief, input)
	if err != nil {
		return Result{}, err
	}
	return Result{part1, part2}, nil
}

type MonkeyId string
type WorryLevel int
type Inspect func(wl WorryLevel) WorryLevel
type DecideNext func(wl WorryLevel) MonkeyId
type ReliefMaker func(monkeys []*Monkey) Relief
type Relief func(wl WorryLevel) WorryLevel

type Monkey struct {
	id      MonkeyId
	items   []WorryLevel
	divisor int
	inspect Inspect
	next    DecideNext
}

func (p MonkeyInTheMiddle) solve(rounds int, rm ReliefMaker, input *Input) (string, error) {
	monkeys, err := p.parseNotes(input)
	if err != nil {
		return "", err
	}
	counts, err := p.processRounds(rounds, rm, monkeys)
	if err != nil {
		return "", err
	}
	monkeyBusinessLevel, err := p.calcMonkeyBusinessLevel(counts)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(monkeyBusinessLevel), nil
}

func (p MonkeyInTheMiddle) calcMonkeyBusinessLevel(inspections map[MonkeyId]int) (int, error) {
	if len(inspections) < 2 {
		return 0, fmt.Errorf("input contains %d monkeys; needs at least %d", len(inspections), 2)
	}
	var counts []int
	for _, count := range inspections {
		counts = append(counts, count)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(counts)))
	return counts[0] * counts[1], nil
}

func (p MonkeyInTheMiddle) processRounds(rounds int, rm ReliefMaker, monkeys []*Monkey) (map[MonkeyId]int, error) {
	monkeysById := map[MonkeyId]*Monkey{}
	inspections := map[MonkeyId]int{}
	for _, monkey := range monkeys {
		if _, ok := monkeysById[monkey.id]; ok {
			return nil, fmt.Errorf("duplicate monkey id: %s", monkey.id)
		}
		monkeysById[monkey.id] = monkey
		inspections[monkey.id] = 0
	}
	reliefFn := rm(monkeys)

	for round := 0; round < rounds; round++ {
		for _, monkey := range monkeys {
			items := monkey.items
			monkey.items = nil
			for _, wl := range items {
				newWl := reliefFn(monkey.inspect(wl))
				recipient := monkey.next(newWl)
				if nextMonkey, ok := monkeysById[recipient]; ok {
					nextMonkey.items = append(nextMonkey.items, newWl)
				} else {
					return nil, fmt.Errorf("invalid recipient id: %s", recipient)
				}
				inspections[monkey.id]++
			}
		}
	}
	return inspections, nil
}

var (
	monkeyRgx = regexp.MustCompile(`^Monkey (.+?):$`)
	itemsRgx  = regexp.MustCompile(`^\s+Starting items: (.*)`)
	opRgx     = regexp.MustCompile(`^\s+Operation: new = (old|\d+) ([+*]) (old|\d+)$`)
	testRgx   = regexp.MustCompile(`^\s+Test: divisible by (\d+)$`)
	trueRgx   = regexp.MustCompile(`^\s+If true: throw to monkey (.+)$`)
	falseRgx  = regexp.MustCompile(`^\s+If false: throw to monkey (.+)$`)
)

func (p MonkeyInTheMiddle) parseNotes(input *Input) ([]*Monkey, error) {
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
		divisor, _ := strconv.Atoi(testLineGroups[0][1])
		whenTrue := MonkeyId(trueLineGroups[0][1])
		whenFalse := MonkeyId(falseLineGroups[0][1])
		decideNext := func(wl WorryLevel) MonkeyId {
			if int(wl)%divisor == 0 {
				return whenTrue
			}
			return whenFalse
		}

		result = append(result, &Monkey{
			id:      id,
			items:   items,
			divisor: divisor,
			inspect: operation,
			next:    decideNext,
		})
		i += 6
	}
	return result, nil
}

func (p MonkeyInTheMiddle) divBy3Relief(monkeys []*Monkey) Relief {
	return func(wl WorryLevel) WorryLevel { return wl / 3 }
}

func (p MonkeyInTheMiddle) modByDivisorsRelief(monkeys []*Monkey) Relief {
	divisor := 1
	for _, m := range monkeys {
		divisor *= m.divisor
	}
	return func(wl WorryLevel) WorryLevel { return WorryLevel(int(wl) % divisor) }
}
