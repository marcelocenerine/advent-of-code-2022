// Solution to https://adventofcode.com/2022/day/7
package adventofcode

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"

	"golang.org/x/exp/maps"
)

type NoSpaceLeftOnDevice struct{}

func (s NoSpaceLeftOnDevice) Details() Details {
	return Details{Day: 7, Description: "No Space Left On Device"}
}

func (s NoSpaceLeftOnDevice) Solve(input *Input) (Result, error) {
	root, err := parseCommandsOutput(input)
	if err != nil {
		return Result{}, err
	}
	part1 := part1SumOfDirSizesUpTo100000(root)
	part2, err := part2SpaceToBeFreedUp(root)
	if err != nil {
		return Result{}, err
	}

	return Result{
		Part1: strconv.Itoa(part1),
		Part2: strconv.Itoa(part2),
	}, nil
}

func part1SumOfDirSizesUpTo100000(root *Dir) int {
	totalSize := 0
	for _, size := range dirSizes(root) {
		if size <= 100000 {
			totalSize += size
		}
	}
	return totalSize
}

func part2SpaceToBeFreedUp(root *Dir) (int, error) {
	sizesByDir := dirSizes(root)
	totalDiskSpace := 70000000
	unused := totalDiskSpace - sizesByDir["/"]
	toFreeUp := 30000000 - unused
	dirSizes := maps.Values(sizesByDir)
	sort.Ints(dirSizes)
	for _, size := range dirSizes {
		if size >= toFreeUp {
			return size, nil
		}
	}
	return 0, fmt.Errorf("can't free up %d of disk space", toFreeUp)
}

func dirSizes(root *Dir) map[string]int {
	sizes := map[string]int{}
	var dfs func(string, *Dir) int
	dfs = func(path string, curr *Dir) int {
		size := 0
		for _, entry := range curr.Entries {
			if f, ok := entry.(*File); ok {
				size += f.Size
			}
			if d, ok := entry.(*Dir); ok {
				path := fmt.Sprintf("%s%s/", path, d.Name)
				size += dfs(path, d)
			}
		}
		sizes[path] = size
		return size
	}
	dfs("/", root)
	return sizes
}

type FsEntry any

type File struct {
	Name   string
	Size   int
	Parent *Dir
}

type Dir struct {
	Name    string
	Entries map[string]FsEntry
	Parent  *Dir
}

func (d *Dir) cd(dest string) (*Dir, error) {
	switch dest {
	case "/":
		curr := d
		for curr.Parent != nil {
			curr = curr.Parent
		}
		return curr, nil
	case "..":
		if d.Parent != nil {
			return d.Parent, nil
		}
		return nil, errors.New("cannot go up another level")
	default:
		if entry, ok := d.Entries[dest]; ok {
			if dir, ok := entry.(*Dir); ok {
				return dir, nil
			}
			return nil, fmt.Errorf("%s is not a directory", dest)
		}
		return nil, fmt.Errorf("dir %s doesn't exit", dest)
	}
}

func parseCommandsOutput(input *Input) (*Dir, error) {
	cdRgx := regexp.MustCompile(`^\$ cd (.+)$`)
	lsRgx := regexp.MustCompile(`^\$ ls$`)
	dirRgx := regexp.MustCompile(`^dir (.+)$`)
	fileRgx := regexp.MustCompile(`^(\d+) (.+)$`)

	root := &Dir{Name: "/", Entries: map[string]FsEntry{}}
	curr := root
	lines := input.Lines()
	i := 0

	for i < len(lines) {
		line := lines[i]
		i++

		if cdRgx.MatchString(line) {
			groups := cdRgx.FindAllStringSubmatch(line, -1)
			dest := groups[0][1]
			if dir, err := curr.cd(dest); err != nil {
				return nil, err
			} else {
				curr = dir
				continue
			}
		}

		if lsRgx.MatchString(line) {
			for i < len(lines) {
				lsOutputLine := lines[i]
				if dirRgx.MatchString(lsOutputLine) {
					groups := dirRgx.FindAllStringSubmatch(lsOutputLine, -1)
					dir := groups[0][1]
					curr.Entries[dir] = &Dir{Name: dir, Parent: curr, Entries: map[string]FsEntry{}}
					i++
					continue
				}
				if fileRgx.MatchString(lsOutputLine) {
					groups := fileRgx.FindAllStringSubmatch(lsOutputLine, -1)
					size, _ := strconv.Atoi(groups[0][1])
					name := groups[0][2]
					curr.Entries[name] = &File{Name: name, Size: size, Parent: curr}
					i++
					continue
				}
				break
			}
			continue
		}
		return nil, fmt.Errorf("line %d is invalid: %s", i, line)
	}

	return root, nil
}
