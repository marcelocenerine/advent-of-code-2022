// Solution to https://adventofcode.com/2022/day/6
package adventofcode

import (
	"strconv"
)

type TuningTrouble struct{}

func (s TuningTrouble) Details() Details {
	return Details{Day: 6, Description: "Tuning Trouble"}
}

func (s TuningTrouble) Solve(input *Input) (Result, error) {
	part1 := charCountBeforeStartOfPacket(input)

	return Result{
		Part1: strconv.Itoa(part1),
		Part2: "",
	}, nil
}

const MarkerLength = 4

func charCountBeforeStartOfPacket(input *Input) int {
	runes := []rune(*input)
	runeCount := map[rune]int{}

	for hi, right := range runes {
		runeCount[right]++

		if hi < MarkerLength {
			continue
		}

		lo := hi - MarkerLength
		left := runes[lo]
		if leftCount := runeCount[left] - 1; leftCount == 0 {
			delete(runeCount, left)
		} else {
			runeCount[left] = leftCount
		}

		if len(runeCount) == MarkerLength {
			return hi + 1
		}
	}

	return -1
}
