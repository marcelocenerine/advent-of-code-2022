// Solution to https://adventofcode.com/2022/day/13
package adventofcode

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type DistressSignal struct{}

func (p DistressSignal) Details() Details {
	return Details{Day: 13, Description: "Distress Signal"}
}

func (p DistressSignal) Solve(input *Input) (Result, error) {
	pairs, err := p.parse(input)
	if err != nil {
		return Result{}, err
	}

	return Result{
		Part1: strconv.Itoa(p.sumIndexOfPairsInTheRightOrder(pairs)),
		Part2: strconv.Itoa(p.decoderKey(pairs)),
	}, nil
}

func (p DistressSignal) sumIndexOfPairsInTheRightOrder(pairs []packetPair) int {
	var result int
	for i, pair := range pairs {
		if pair.isInOrder() {
			result += i + 1
		}
	}
	return result
}

func (p DistressSignal) decoderKey(pairs []packetPair) int {
	var packets []packet
	for _, pair := range pairs {
		packets = append(packets, pair.left, pair.right)
	}
	sort.Slice(packets, func(i, j int) bool { return packets[i].compare(packets[j]) < 0 })

	i, j := 1, len(packets)
	div1 := packet{listValue{intValue(2).asList()}}
	div2 := packet{listValue{intValue(6).asList()}}

	for ; i < j && div1.compare(packets[i-1]) > 0; i++ {
	}
	for ; j > i && div2.compare(packets[j-1]) < 0; j-- {
	}

	return i * (j + 2)
}

type packetData interface {
	asList() listValue
}

type listValue []packetData

func (l listValue) asList() listValue { return l }

type intValue int

func (i intValue) asList() listValue { return listValue{i} }

type packet struct {
	data listValue
}

func (p packet) compare(that packet) int {
	var recurse func(left, right listValue) int
	recurse = func(left, right listValue) int {
		for i := 0; i < len(left); i++ {
			if i == len(right) {
				return 1 // right ran out of items
			}

			li := left[i]
			ri := right[i]
			_, liListOk := li.(listValue)
			_, riListOk := ri.(listValue)

			if liListOk || riListOk {
				if order := recurse(li.asList(), ri.asList()); order != 0 {
					return order
				}
				continue
			}

			liInt, _ := li.(intValue)
			riInt, _ := ri.(intValue)

			if liInt == riInt {
				continue
			}

			return int(liInt - riInt)
		}

		// did left run out of items?
		return len(left) - len(right)
	}

	return recurse(p.data, that.data)
}

type packetPair struct {
	left, right packet
}

func (p packetPair) isInOrder() bool {
	return p.left.compare(p.right) < 1
}

func (p DistressSignal) parse(input *Input) ([]packetPair, error) {
	var result []packetPair
	var packets []packet
	lines := input.Lines()

	for idx, line := range lines {
		line = strings.TrimSpace(line)

		if line != "" {
			packet, err := p.parsePacketLine(line)
			if err != nil {
				return nil, err
			}
			packets = append(packets, packet)
		}

		if line == "" || idx == len(lines)-1 {
			if len(packets) != 2 {
				return nil, fmt.Errorf("invalid input: a pair of packets is expected before line %d", idx+1)
			}

			pair := packetPair{left: packets[0], right: packets[1]}
			result = append(result, pair)
			packets = nil
		}
	}

	return result, nil
}

func (p DistressSignal) parsePacketLine(line string) (packet, error) {
	data, end, err := p.parsePacketData(line, 0)
	if err != nil {
		return packet{}, err
	}

	if end < len(line)-1 {
		return packet{}, fmt.Errorf("invalid packet data: %s", line)
	}

	return packet{data}, nil
}

func (p DistressSignal) parsePacketData(line string, start int) (listValue, int, error) {
	if line[start] != '[' {
		return listValue{}, 0, fmt.Errorf("invalid packet data: %s", line)
	}

	var (
		result     listValue
		curInt     string
		closeFound bool
	)
	i := start + 1

	for ; i < len(line); i++ {
		curChar := line[i]

		if curChar >= '0' && curChar <= '9' {
			curInt += string(curChar)
			continue
		}

		if curChar == ',' {
			if curInt != "" {
				value, _ := strconv.Atoi(curInt)
				result = append(result, intValue(value))
				curInt = ""
			}
			continue
		}

		if curChar == '[' {
			data, end, err := p.parsePacketData(line, i)
			if err != nil {
				return listValue{}, 0, err
			}

			result = append(result, data)
			i = end
			continue
		}

		if curChar == ']' {
			if curInt != "" {
				value, _ := strconv.Atoi(curInt)
				result = append(result, intValue(value))
			}
			closeFound = true
			break
		}
	}

	if !closeFound {
		return listValue{}, 0, fmt.Errorf("invalid packet data: %s", line)
	}

	return result, i, nil
}
