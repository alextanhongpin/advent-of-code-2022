package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("pt1:", part1(input1))
	fmt.Println("pt1:", part1(input2))
	fmt.Println()
	fmt.Println("pt2:", part2(input1))
	fmt.Println("pt2:", part2(input2))
}

func oneOfCyle(cycle int) bool {
	switch cycle {
	case 20, 60, 100, 140, 180, 220:
		return true
	default:
		return false
	}
}

func part1(input string) int {
	rows := parse(input)
	var x, cycle, total int
	x = 1
	for _, row := range rows {
		parts := strings.Fields(row)
		switch parts[0] {
		case "addx":
			n := toInt(parts[1])
			cycle++
			if oneOfCyle(cycle) {
				total += cycle * x
			}

			cycle++
			if oneOfCyle(cycle) {
				total += cycle * x
			}
			x += n
		case "noop":
			cycle++
			if oneOfCyle(cycle) {
				total += cycle * x
			}
		default:
			panic("invalid command")
		}
		if cycle > 220 {
			break
		}
	}

	return total
}

func part2(input string) int {
	rows := parse(input)
	grids := make([][]rune, 6)
	for i := range grids {
		grids[i] = make([]rune, 40)
		for j := range grids[i] {
			grids[i][j] = '.'
		}
	}

	var x, cycle int
	x = 1

	refresh := func() {
		// x is the midpoint of the 3px sprite.
		// At every tick, as long as the pixel overlaps with the sprite position,
		// the pixel is lit `#`.
		if cycle%40 >= x-1 && cycle%40 <= x+1 {
			grids[cycle/40][cycle%40] = '#'
		}
	}

	for _, row := range rows {
		refresh()

		parts := strings.Fields(row)
		switch parts[0] {
		case "addx":
			n := toInt(parts[1])

			cycle++
			refresh()

			cycle++
			x += n
		case "noop":
			cycle++
			refresh()
		default:
			panic("invalid command")
		}
		if cycle > 240 {
			break
		}
	}

	draw(grids)

	return 0
}

func draw(grids [][]rune) {
	for _, grid := range grids {
		fmt.Println(string(grid))
	}
	fmt.Println()
}

func parse(input string) []string {
	return strings.Split(input, "\n")
}

func toInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

var input1 = `addx 15
addx -11
addx 6
addx -3
addx 5
addx -1
addx -8
addx 13
addx 4
noop
addx -1
addx 5
addx -1
addx 5
addx -1
addx 5
addx -1
addx 5
addx -1
addx -35
addx 1
addx 24
addx -19
addx 1
addx 16
addx -11
noop
noop
addx 21
addx -15
noop
noop
addx -3
addx 9
addx 1
addx -3
addx 8
addx 1
addx 5
noop
noop
noop
noop
noop
addx -36
noop
addx 1
addx 7
noop
noop
noop
addx 2
addx 6
noop
noop
noop
noop
noop
addx 1
noop
noop
addx 7
addx 1
noop
addx -13
addx 13
addx 7
noop
addx 1
addx -33
noop
noop
noop
addx 2
noop
noop
noop
addx 8
noop
addx -1
addx 2
addx 1
noop
addx 17
addx -9
addx 1
addx 1
addx -3
addx 11
noop
noop
addx 1
noop
addx 1
noop
noop
addx -13
addx -19
addx 1
addx 3
addx 26
addx -30
addx 12
addx -1
addx 3
addx 1
noop
noop
noop
addx -9
addx 18
addx 1
addx 2
noop
noop
addx 9
noop
noop
noop
addx -1
addx 2
addx -37
addx 1
addx 3
noop
addx 15
addx -21
addx 22
addx -6
addx 1
noop
addx 2
addx 1
noop
addx -10
noop
noop
addx 20
addx 1
addx 2
addx 2
addx -6
addx -11
noop
noop
noop`
var input2 = `addx 1
noop
addx 29
addx -24
addx 4
addx 3
addx -2
addx 3
addx 1
addx 5
addx 3
addx -2
addx 2
noop
noop
addx 7
noop
noop
noop
addx 5
addx 1
noop
addx -38
addx 21
addx 8
noop
addx -19
addx -2
addx 2
addx 5
addx 2
addx -12
addx 13
addx 2
addx 5
addx 2
addx -18
addx 23
noop
addx -15
addx 16
addx 7
noop
noop
addx -38
noop
noop
noop
noop
noop
noop
addx 8
addx 2
addx 3
addx -2
addx 4
noop
noop
addx 5
addx 3
noop
addx 2
addx 5
noop
noop
addx -2
noop
addx 3
addx 6
noop
addx -38
addx -1
addx 35
addx -6
addx -19
addx -2
addx 2
addx 5
addx 2
addx 3
noop
addx 2
addx 3
addx -2
addx 2
noop
addx -9
addx 16
noop
addx 9
addx -3
addx -36
addx -2
addx 11
addx 22
addx -28
noop
addx 3
addx 2
addx 5
addx 2
addx 3
addx -2
addx 2
noop
addx 3
addx 2
noop
addx -11
addx 16
addx 2
addx 5
addx -31
noop
addx -6
noop
noop
noop
noop
noop
addx 7
addx 30
addx -24
addx -1
addx 5
noop
noop
noop
noop
noop
addx 5
noop
addx 5
noop
addx 1
noop
addx 2
addx 5
addx 2
addx 1
noop
noop
noop
noop`
