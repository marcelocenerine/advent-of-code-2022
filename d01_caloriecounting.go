// Solution to https://adventofcode.com/2022/day/1
package adventofcode

import (
	"fmt"
	"sort"
	"strconv"
)

type CalorieCounting struct{}

func (c CalorieCounting) Details() Details {
	return Details{Day: 1, Description: "Calorie Counting"}
}

func (c CalorieCounting) Solve(input *Input) (Result, error) {
	caloriesPerElf, err := caloriesCount(input)
	if err != nil {
		return Result{}, err
	}
	elfCount := len(caloriesPerElf)
	if elfCount < 3 {
		return Result{}, fmt.Errorf("the number of elfs in the input is %d; the required is %d", elfCount, 3)
	}
	top1Sum := caloriesPerElf[0]
	top3Sum := caloriesPerElf[0] + caloriesPerElf[1] + caloriesPerElf[2]
	return Result{
		Part1: strconv.Itoa(top1Sum),
		Part2: strconv.Itoa(top3Sum),
	}, nil
}

func caloriesCount(input *Input) ([]int, error) {
	var caloriesPerElf []int
	var curElfCalories int

	for _, line := range input.Lines() {
		if len(line) == 0 {
			caloriesPerElf = append(caloriesPerElf, curElfCalories)
			curElfCalories = 0
			continue
		}

		calories, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}
		curElfCalories += calories
	}

	sort.Sort(sort.Reverse(sort.IntSlice(caloriesPerElf)))
	return caloriesPerElf, nil
}
