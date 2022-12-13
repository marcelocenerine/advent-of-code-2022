package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if max, err := caloriesCountFromFile("input.txt"); err != nil {
		fmt.Printf("Invalid input: %v", err)
	} else {
		fmt.Printf("The Elf carrying the most calories carries %d calories in total\n", max)
	}
}

func caloriesCountFromFile(path string) (int, error) {
	if bytes, err := os.ReadFile(path); err != nil {
		return 0, err
	} else {
		return caloriesCount(string(bytes[:]))
	}
}

func caloriesCount(input string) (int, error) {
	var curCalories, maxCalories int

	for _, line := range strings.Split(input, "\n") {
		if len(line) == 0 {
			if curCalories > maxCalories {
				maxCalories = curCalories
			}

			curCalories = 0
			continue
		}

		calories, err := strconv.Atoi(line)
		if err != nil {
			return 0, err
		}
		curCalories += calories
	}

	return maxCalories, nil
}
