package main

import (
	"container/heap"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	fmt.Println(part1(input1)) // 18
	fmt.Println(part1(input2)) // 308
	fmt.Println(part2(input1)) // 54
	fmt.Println(part2(input2)) // 908
}

var moves = []point{
	{0, 0},
	{1, 0},
	{-1, 0},
	{0, -1},
	{0, 1},
}

// makeCache stores the position of blizzard at time t.
// Since the blizzard will always circle back, we can easily find the position
// of the blizzard moving left and right by using (t % maxX - 2)
// and up and down by using (t % maxY - 2).
func makeCache(grid map[point]rune, maxX, maxY int) (map[point]bool, map[int]map[point]bool, map[int]map[point]bool) {
	leftright := make(map[int]map[point]bool)
	for i := 0; i < maxX-2; i++ {
		leftright[i] = make(map[point]bool)
	}

	updown := make(map[int]map[point]bool)
	for i := 0; i < maxY-2; i++ {
		updown[i] = make(map[point]bool)
	}

	// Wall's position does not change with time.
	walls := make(map[point]bool)
	for k, v := range grid {
		switch v {
		case '#':
			walls[k] = true
		}
	}

	for i := 0; i < maxX-2; i++ {
		for p, ch := range grid {
			if ch != '<' && ch != '>' {
				continue
			}
			dx := 1
			if ch == '<' {
				dx = -1
			}
			pp := p

			x := p.x      // original position
			x--           // minus 1 since the wall doesn't count
			x += dx * i   // at time i, we already move distance dx
			x += maxX - 2 // add maxX - 2, in case the value becomes negative
			x %= maxX - 2 // modulo by maxX - 2, in case the value becomes greater than maxX - 2
			x++           // add back the 1 to get the offset after the wall

			pp.x = x
			leftright[i][pp] = true
		}
	}

	for i := 0; i < maxY-2; i++ {
		for p, ch := range grid {
			if ch != '^' && ch != 'v' {
				continue
			}
			dy := 1
			if ch == '^' {
				dy = -1
			}
			pp := p
			pp.y = (((p.y-1)+dy*i)+maxY-2)%(maxY-2) + 1
			updown[i][pp] = true
		}
	}

	return walls, leftright, updown
}

func part1(input string) int {
	maxX, maxY, grid := parse(input)
	start := point{x: 1, y: 0}
	end := point{x: maxX - 2, y: maxY - 1}

	// Close the up and bottom hole so that we don't go out of boundary.
	// This modifies the original grid, so be careful especially since it affects
	// the maxX and maxY.
	grid[point{x: start.x, y: start.y - 1}] = '#'
	grid[point{x: end.x, y: end.y + 1}] = '#'

	minutes := 0

	return solve(minutes, start, end, maxX, maxY, grid)
}

func part2(input string) int {
	maxX, maxY, grid := parse(input)
	start := point{x: 1, y: 0}
	end := point{x: maxX - 2, y: maxY - 1}

	grid[point{x: start.x, y: start.y - 1}] = '#'
	grid[point{x: end.x, y: end.y + 1}] = '#'

	t0 := solve(0, start, end, maxX, maxY, grid)
	t1 := solve(t0, end, start, maxX, maxY, grid)
	t2 := solve(t1, start, end, maxX, maxY, grid)
	return t2
}

func solve(minutes int, start, end point, maxX, maxY int, grid map[point]rune) int {
	walls, leftright, updown := makeCache(grid, maxX, maxY)
	cache := make(map[Item]bool)

	q := make(PriorityQueue, 0, 1_000_000)
	heap.Push(&q, Item{point: start, minutes: minutes})

	i := 0
	for q.Len() > 0 {
		s := heap.Pop(&q).(Item)
		if s.point == end {
			return s.minutes
		}

		if cache[s] {
			continue
		}
		cache[s] = true

		if i%1_000 == 0 {
			fmt.Println(i, s.minutes)
		}
		i++

		for _, m := range moves {
			p := s.point
			p.x += m.x
			p.y += m.y

			if walls[p] {
				continue
			}

			// Check blizzards on the left or right.
			// Note that we start at minute + 1.
			// maxX - 2 excludes the walls on both side.
			dx := (s.minutes + 1) % (maxX - 2)
			if leftright[dx][p] {
				continue
			}

			dy := (s.minutes + 1) % (maxY - 2)
			if updown[dy][p] {
				continue
			}

			ns := s
			ns.minutes++
			ns.point = p

			heap.Push(&q, ns)
		}
	}

	return 0
}

func toInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

type point struct {
	x, y int
}

func parse(input string) (maxX, maxY int, res map[point]rune) {
	res = make(map[point]rune)
	rows := strings.Split(input, "\n")

	maxY = len(rows)
	for y, row := range rows {
		maxX = len(row)

		for x, r := range row {
			if r == '.' {
				continue
			}
			res[point{x: x, y: y}] = r
		}
	}

	return
}

type Item struct {
	minutes int
	point   point
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []Item

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].minutes < pq[j].minutes }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x any) {
	item := x.(Item)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

var input1 = `#.######
#>>.<^<#
#.<..<<#
#>v.><>#
#<^v^^>#
######.#`
var input2 = `#.########################################################################################################################
#>v<^^><<>^>^>v<>>^^.v>^vv<^v<v><vv.^^^^v<v<v><<^v.<>^^^>>>>><^<.vvvvv.^v>vv<>v<v<^>^v<^v.<>>v<><v^v>.v>^vv^>><v>v>.^>><>#
#<.v<v>>>.>>>>^<vv><^^v>>v>^><v^>^<.<v^<>>><v<vv<<^^<^^v^<.>v.v>>.v<^><v<v<>^><^>^>><.<<<<^>^^^v.v.>v><^<^v^>.>^>>vv.v^<<#
#<<>v>^<<^v^<<<vv.>^^>^^v^>>v<>v>v<><^^v<^>><^<>^^v^vv<v<<<>^v^v<.<^>^<>>>.<><<<v<v<<<vv<>>v.^^<v.vv<^<^><^<vv^v<v<v>v.^<#
#.v^v<v<>>>.v^vv<<v^<vv<<vvv^>^^v.v^v^^v<>^^^^^><.v^v>>v>>^>vvv^^<<^>^v^.v<^<<^v^v>>vv.<<^<<v<v>>^v>^>v.v^><><vvv^>v<>^<<#
#<v<^<v<^<>v>>v^>>^<<>>v.^^<>vv^<><.v^>><<^<v><<^>>^^^<vv>^<v>.>><<.v^v^<v>^<v>>>>vvv<><^><<<v>v>v.^<<.<v^<>vv<v<<>v>^>v>#
#.^.<>v>v^v<v>>v<>v>vv<.^><^v<<><>^.^>>v<^v.v^^^^<vv>^<vv^>vv^>^.<>.<^^v<>><^^^.><^v>v^^<>v>>v><>>v.^^<>>v<<^<^<v<vv.>^<<#
#>^v^>^v>...^<>^.<v>.v<v<>^^v<<<^.v>v>^.><<.v<v<>.>v>>v>>><<v<><^.v>.>>>^>^>.v<>>>v<^^>^v<<^v>^v^<v.v>><>><vvv^v>v^v<v<<<#
#<^.vv><<><>v^^<^.<^<><..v.^vv<v.vvv><v<^v>v<^v>v^v^^^^<>^^^.^^^v^v>>v^>><>v>^<>v^^.^<^^<v>>v^vvvv>>v^><^<<>>>^<<^v^^<<><#
#<>><^<^<<^v^>>>v<.>..v>.^>>^v^<<v^<<vv<^v^><<.><^^><.vvv^>>.>>^<^vv><.v>^^>>>.>^vv<^<<>vv<>vv^^^v.>>>><..v>.^<^v<>v.<vv<#
#<<.<><>.<^<<<<^>^<>^vv^>^>^>vv>>v>>v^>.^.>vv<><.v><<>>vvv.v^^vv>>>v<v<><^<^<>^.v<^^>^><.>>^^^^^vv<<^>^<<v<^v^<.v<<v<<.<<#
#>>v>.^.<><>^.vvv<vv>^>.^^.v>^v.v>^<vv^v>v>^v^vv<<>><v.<><<v><>^<^<v<>.<^<^^<>^<>>v^v.^<v^^<v^>v<<^.^>>^v>^^^>vv.v.>^v>^>#
#>^<>v^^^v<v^vv>>v.^>>>^v>><>^v^><v^.v><<>>><^<<.<>v<><><^v^>^^<^^>>^v^>^>^vv^<>>>^v^<.v^^^v<v<v...>^>^<^^v<v^^v>^^>.<^^>#
#<v<v^^v<>><<<^vvv<vv>>.^.<>>.^v^^<>>v^^>v><>v><v<^^.v^>>.vv^v><>^^v<>>^>^<<^.>v<v^vvv^.><<<v.^^<vvv^<.>>^.<<vv^<<<>.^>><#
#>.>^.v.>>^<v><><^<.>^v<><>v^^<.>vvv<><>.<v<>.><<.^v^^v^.>vvvv^>.<>><<^v<<>^<<<><<>>^<^>^v<^v<>v^>v<v<.>^<^vv<>^<^>>^>^<<#
#<.^.>^vv<><>^^<^<.<>vv.>^<<<>v><.<<v<>^<>>v<^^<v^v^^>v<v<^<vv<.v^<>v<^>>.<>^<^.vv><.><v<^<<^<<>^<>v<vvv^v^^>><v>v>..v<<>#
#<^v>><<v.^^.<^>^<^vvv.^v^v<>vv^v<^vvvvvv^v><>>.^<<^v>v.v.>.^<<^^>>>v^><v^>>>.v^^>.>><^.^^.>.>v^<v^>^<<v<<^><>v<^^>^<v<v<#
#.v><v>>>v^>^^v<>^<<>>^<v>^><<v^><^.<<v<vv<<^v^^><>><^^><><^^^vv>.^>.<v.^v<v<^<<.v.>>>><<<^v<<.<^<>v.^<^<vvv<.v<.v><^>v<>#
#<v<v^<^v><<^<>^>>>>^^v<>^^>v<v.^^^^^.^<^<<^><<v^><^.><^><v^>.^>v^v>v.><^<^>v<>^>.<>^..<<>^^>^<><^<^..vv>v^^^v>^^<..^.v<>#
#>v<<v<<<^.^vvv>>v.>v^v>v<><v.^><^v<>^v<v^v<v<.^^.>v>.v>vvvv^>vv^>^^v><v^>>>.<v^>><.<v><<.>>^^v^vvv<vv^^^<^>v^>vv>^>^vv>>#
#><<v>^>>>>vv<<><>^>^v>^vv^<^v><>v^.<.^.v>^<<^v>v<.^v.>^v<^^^>v^^.^>vv.<v^>>^.><<>>^<<>.>v.v>>^..>><<..v^>><>v<>.<^v.>>^<#
#.<<^^<.^v>>>>.>^vv<^<<>>>>^v>v<^.v^<^>>v^<^>^^>>v<>^<v^^^^>><^<^^>>v<.<<v.^v<>>.vv^^^<v^^.^v.<>^>^>><><^.vv<v>v>v<^^>.<>#
#<<<v>.^v>v>v<^^>vv^<vv<vv><v>v<vv<.>vv<^^^<.v<>..^<<<<<^^^>><v^^>>^<<^v^>><.^<^.<^.<v<vv<..^.><<>vv<^<v<v.<>.<<^v<<^>vv<#
#.>^^.<<>.>>>.<^<>.<<v<^>v>v>v^v^<v>v^<^>^<<>v^^v^v^^^<>^vvv^vv<.^^>.v^<.>>^>vvv^>.v^>v^v^>.^^<v<<vv^^>v>^vvvvv<vvv^<<>>>#
#<^v^.<<<><^<<v^^^v^^^<^^<<>.^<^v<.>.^^..^v<<^>^.<<>>.<^>v.^<v>.v^^.<vv^^>v^><v^><^.>>>>>^v<v<^vv^^>^v^.^.^^^>v^<.><>^^>>#
#.v^vv.^v>vv>>v<vv^^v<v>><.<>><^^>^v.vvvv>><^^>v^v>.v^<.<^.<^<^v.><v>^v<<^^<>vv>>>^v^v^<^<^^^>..<><^<^.vv<<<v>v^.^>^<>^><#
########################################################################################################################.#`
