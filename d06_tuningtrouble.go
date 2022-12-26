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
	part1 := charCountUntilEndOfMarker(input, PacketMarkerLength)
	part2 := charCountUntilEndOfMarker(input, MessageMarkerLength)

	return Result{
		Part1: strconv.Itoa(part1),
		Part2: strconv.Itoa(part2),
	}, nil
}

const PacketMarkerLength = 4
const MessageMarkerLength = 14

func charCountUntilEndOfMarker(input *Input, markerLength int) int {
	runes := []rune(*input)
	runeCount := map[rune]int{}

	for hi := 0; hi < len(runes); hi++ {
		right := runes[hi]
		runeCount[right]++

		if hi < markerLength {
			continue
		}

		lo := hi - markerLength
		left := runes[lo]
		if leftCount := runeCount[left] - 1; leftCount == 0 {
			delete(runeCount, left)
		} else {
			runeCount[left] = leftCount
		}

		if len(runeCount) == markerLength {
			return hi + 1
		}
	}

	return -1
}
