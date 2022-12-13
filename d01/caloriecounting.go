package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	if caloriesPerElf, err := caloriesCountFromFile("input.txt"); err != nil {
		fmt.Printf("Invalid input: %v", err)
	} else if len(caloriesPerElf) < 3 {
		fmt.Printf("The number of elfs in the input is %d; the required is %d", len(caloriesPerElf), 3)
	} else {
		top3Sum := caloriesPerElf[0] + caloriesPerElf[1] + caloriesPerElf[2]
		fmt.Printf("The top 3 Elfs carrying the most calories carries %d calories in total\n", top3Sum)
	}
}

func caloriesCountFromFile(path string) ([]int, error) {
	if bytes, err := os.ReadFile(path); err != nil {
		return nil, err
	} else {
		return caloriesCount(string(bytes[:]))
	}
}

func caloriesCount(input string) ([]int, error) {
	var caloriesPerElf []int
	var curElfCalories int

	for _, line := range strings.Split(input, "\n") {
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
