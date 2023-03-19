// Solution to https://adventofcode.com/2022/day/10
package adventofcode

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
)

type CathodeRayTube struct{}

func (p CathodeRayTube) Details() Details {
	return Details{Day: 10, Description: "Cathode-Ray Tube"}
}

func (p CathodeRayTube) Solve(input *Input) (Result, error) {
	instructions, err := p.parseInstructions(input)
	if err != nil {
		return Result{}, err
	}
	part1 := p.sumOfSignalStrengths(instructions, 20, 40)
	part2 := p.crtDraw(instructions, 40, 3)
	if err != nil {
		return Result{}, err
	}
	return Result{
		Part1: strconv.Itoa(part1),
		Part2: part2,
	}, nil
}

func (p CathodeRayTube) sumOfSignalStrengths(instructions []Add, firstCheckpoint, checkpointIntervals int) int {
	checkpoint := firstCheckpoint
	register := 1
	result := 0
	for cycle := 1; cycle <= len(instructions); cycle++ {
		instr := instructions[cycle-1]
		if cycle == checkpoint {
			result += cycle * register
			checkpoint += checkpointIntervals
		}
		register += instr.Value
	}
	return result
}

func (p CathodeRayTube) crtDraw(instructions []Add, screenWidth, spriteWidth int) string {
	register := 1
	var buffer bytes.Buffer

	for cycle, instr := range instructions {
		pixel := cycle % screenWidth
		spritePos := register - int(spriteWidth/2)

		if pixel >= spritePos && pixel < spritePos+spriteWidth {
			buffer.WriteRune('#')
		} else {
			buffer.WriteRune('.')
		}
		if pixel == screenWidth-1 && cycle < len(instructions)-1 {
			buffer.WriteRune('\n')
		}
		register += instr.Value
	}
	return buffer.String()
}

type Add struct {
	Value int
}

var (
	noopRgx = regexp.MustCompile(`^noop$`)
	addxRgx = regexp.MustCompile(`^addx (-?\d+)$`)
)

func (p CathodeRayTube) parseInstructions(input *Input) ([]Add, error) {
	var result []Add

	for _, line := range input.Lines() {
		result = append(result, Add{0})
		switch {
		case noopRgx.MatchString(line):
		case addxRgx.MatchString(line):
			groups := addxRgx.FindAllStringSubmatch(line, -1)
			value, _ := strconv.Atoi(groups[0][1])
			result = append(result, Add{value})
		default:
			return nil, fmt.Errorf("invalid line: %s", line)
		}
	}
	return result, nil
}
