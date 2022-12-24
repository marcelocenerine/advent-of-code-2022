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
	part1 := charCountBeforeEndOfMarker(input, 0, PacketMarkerLength)
	part2 := charCountBeforeEndOfMarker(input, part1, MessageMarkerLength)

	return Result{
		Part1: strconv.Itoa(part1),
		Part2: strconv.Itoa(part2),
	}, nil
}

const PacketMarkerLength = 4
const MessageMarkerLength = 14

func charCountBeforeEndOfMarker(input *Input, offset, markerLength int) int {
	runes := []rune(*input)
	runeCount := map[rune]int{}

	for hi := offset; hi < len(runes); hi++ {
		right := runes[hi]
		runeCount[right]++

		if hi-offset < markerLength {
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
