package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func main() {
	fmt.Println(part1(input1, 2022))
	fmt.Println(part1(input2, 2022))
	fmt.Println(part1(input1, 1000000000000))
	fmt.Println(part1(input2, 1000000000000))
}

func part1(jetstream string, n int) int {
	rocks := parse(rockRaw)
	var tower []rock

	makeRock := func() rock {
		if len(tower) == 0 {
			r := rocks[0]
			r = r.shift(2+1, 3+1) // +1 to include the wall.
			return r
		}

		maxY := math.MinInt
		for _, r := range tower {
			_, max := r.boundY()
			if max > maxY {
				maxY = max
			}
		}

		r := rocks[len(tower)%len(rocks)]
		r = r.shift(2+1, maxY+3+1) // +1 to include the wall.
		return r
	}

	isOverlap := func(r rock) bool {
		for i := len(tower) - 1; i > -1; i-- {
			if tower[i].overlap(r) {
				return true
			}
		}

		return false
	}

	isOOB := func(r rock) bool {
		minX, maxX := r.boundX()
		return minX < 1 || maxX > 7 || isOverlap(r)
	}

	hitsGroundOrRock := func(r rock) bool {
		// Touches the ground.
		minY, _ := r.boundY()
		if minY < 1 {
			return true
		}

		return isOverlap(r)
	}

	jet := 0
	stream := func(r rock) rock {
		switch jetstream[jet%len(jetstream)] {
		case '<':
			if isOOB(r.move(-1, 0)) {
				return r
			}
			return r.move(-1, 0)
		case '>':
			if isOOB(r.move(1, 0)) {
				return r
			}
			return r.move(1, 0)
		}

		return r
	}

	checkRepeating := func(n int) (string, bool) {
		if len(tower) < n {
			return "", false
		}
		target := tower[len(tower)-n:]
		if len(target) < n {
			return "", false
		}
		rocks := make([]rock, len(target))
		for i := range rocks {
			r := target[i]
			rocks[i] = r.norm(n * 10) // It is important to get the quotient of height, otherwise we can't detect the cycle.
		}

		return fmt.Sprint(rocks), true
	}

	type repeat struct {
		count, height int
	}
	cache := make(map[string]repeat)
	count := 0
	debug := false
	for count != n {
		r := makeRock()
		if debug {
			drawTower(append(tower, r))
		}

		for !hitsGroundOrRock(r) {
			r = stream(r)
			r = r.move(0, -1)
			jet++
		}
		// Undo.
		r = r.move(0, 1)

		tower = append(tower, r)
		count++
		if debug {
			drawTower(tower)
		}
		//Prune to save space.
		prune := 100
		if len(tower) > prune {
			tower = tower[prune/2:]
		}

		// This is tough, we are checking the repeating patterns in the list of rocks.
		if key, ok := checkRepeating(10); ok {
			_, maxY := tower[len(tower)-1].boundY()
			if prev, ok := cache[key]; ok {
				// We will eventually get the rate of change of height by count.
				dc := count - prev.count
				dh := maxY - prev.height
				mul := (n - count) / dc   // Multiplier is the amount we need to multiply to get `n=1000000000000`
				rem := n - mul*dc - count // However, there will be remainder, so we just keep iterating until the remainder is 0, then we can add up the height from existing height + multiplier * rate of change of height.
				if rem == 0 {
					return mul*dh + maxY
				}
			}
			cache[key] = repeat{count: count, height: maxY}
		}
	}

	maxY := math.MinInt
	for _, r := range tower {
		_, max := r.boundY()
		if max > maxY {
			maxY = max
		}
	}

	return maxY
}

func toInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	return n
}

func drawTower(tower []rock) {
	maxX, maxY := 8, math.MinInt
	for _, r := range tower {
		_, y := r.boundY()
		if y > maxY {
			maxY = y
		}
	}

	res := make([][]rune, maxY+1)
	for y := range res {
		res[y] = make([]rune, maxX+1)
		for x := range res[y] {
			if (x == 0 && y == maxY) || (x == maxX && y == maxY) {
				res[y][x] = '+'
			} else if x == 0 || x == maxX {
				res[y][x] = '|'
			} else if y == maxY {
				res[y][x] = '-'
			} else {
				res[y][x] = '.'
			}
		}
	}
	for i, r := range tower {
		for _, p := range r.parts {
			res[maxY-p.y][p.x] = rune(fmt.Sprint(i % 10)[0])
		}
	}
	for _, r := range res {
		fmt.Println(string(r))
	}
	fmt.Println()
}

type point struct {
	x, y int
}

type rock struct {
	parts []point
}

func newRock(input string) rock {
	r := rock{}
	parts := strings.Split(input, "\n")
	for y, p := range parts {
		for x, c := range p {
			if c == '#' {
				r.parts = append(r.parts, point{
					x: x,
					y: len(parts) - 1 - y, // The y-axis starts from bottom.
				})
			}
		}
	}
	return r
}

func (r rock) draw() {
	_, maxX := r.boundX()
	_, maxY := r.boundY()
	res := make([][]rune, maxY+1)
	for y := range res {
		res[y] = make([]rune, maxX+1)
		for x := range res[y] {
			res[y][x] = '.'
		}
	}
	for _, p := range r.parts {
		// y-axis is from bottom.
		res[maxY-p.y][p.x] = '#'
	}
	for _, r := range res {
		fmt.Println(string(r))
	}
}

func (r rock) overlap(o rock) bool {
	if len(r.parts) > len(o.parts) {
		return o.overlap(r)
	}
	cache := make(map[point]bool)
	for _, p := range r.parts {
		cache[p] = true
	}
	for _, p := range o.parts {
		if cache[p] {
			return true
		}
	}

	return false
}

func (r rock) shift(x, y int) rock {
	minX, _ := r.boundX()
	minY, _ := r.boundY()
	dx := x - minX
	dy := y - minY

	rr := rock{parts: make([]point, len(r.parts))}
	for i, p := range r.parts {
		p.x += dx
		p.y += dy
		rr.parts[i] = p
	}
	return rr
}

func (r rock) norm(n int) rock {
	rr := rock{parts: make([]point, len(r.parts))}
	for i, p := range r.parts {
		p.y %= n
		rr.parts[i] = p
	}
	return rr
}

func (r rock) move(x, y int) rock {
	rr := rock{parts: make([]point, len(r.parts))}
	for i, p := range r.parts {
		p.x += x
		p.y += y
		rr.parts[i] = p
	}
	return rr
}

func (r rock) boundX() (min, max int) {
	min, max = math.MaxInt, math.MinInt
	for _, p := range r.parts {
		if p.x < min {
			min = p.x
		}
		if p.x > max {
			max = p.x
		}
	}

	return
}

func (r rock) boundY() (min, max int) {
	min, max = math.MaxInt, math.MinInt
	for _, p := range r.parts {
		if p.y < min {
			min = p.y
		}
		if p.y > max {
			max = p.y
		}
	}

	return
}

func parse(input string) []rock {
	rows := strings.Split(input, "\n\n")
	rocks := make([]rock, len(rows))
	for i, row := range rows {
		rocks[i] = newRock(row)
	}

	return rocks
}

var rockRaw = `####

.#.
###
.#.

..#
..#
###

#
#
#
#

##
##`
var input1 = `>>><<><>><<<>><>>><<<>>><<<><<<>><>><<>>`
var input2 = `>>>><<><>><<<<>>>><<>>><>><><<>>>><<<<>>><<<>>>><<<>><<>>><<<<>>><<>>>><><><<>>>><<<>>><<<><<>>><<<<>><<<<>>>><<<<>>>><<<>>>><<<><><<<>>><<<>>><<<<>>>><<<<>>>><>><<>><<<>><<>><>>><<<<>><<<<><<<<><<>>><<<>><<<<>><<<><>>><<<<>>>><<<>>><<<>><>>><<<<><>>>><<<><<<>>><<<><<>><<<<><<<>>>><<<<>><<<<>>><<<<><<<>>><<><>>><>>><>>>><><<<<>>><<>>><><>>><<<>>><<<<>>><><<<<>>><<<<>>>><<><>>><>><<>><<>>>><<<<>><<<<><<<<>>><>>>><<>>><<><<<>>>><<>>><<<<>><<<>>><<<<>>>><<>>>><<>>><<<<>>>><<>>>><<><>><<<<>>><<<<>><<<><<<><<<><<>>><<<<>>>><<<<>>>><<<><<<>><<>>><>><<>><>>>><>><<<<>>><<>>><<<>>><<>>>><>>><>>><<><<<>>>><<<<>>><<<<>>><<>>>><<<<>><<>>>><<<>>><>>><<<><>>><>>>><<>>>><<>><><<<<><<<>>><<>><<><<>>>><<>>>><><<<<>>>><<<<><<<><>>><<<>><<<>>>><>><<<<>>><><<><<<<>><<<<>><<<<>>><><<<><<<>>><<<<>>><>><<<<><><<<>>>><<<<><>>>><<<>>>><<><<<><>>><<<><<<<><>>><<<>><<<>><<<><><<<>>><<><<<<>>><<<>>>><<>><>><<<>>><>>>><<<<><>>><<>>>><>>>><<>>>><>><><>>><<<>>><<<>><<>>>><<<>>>><><<>>>><<>>><<<<>>><<<<>>><<<<>><<<>><<<>>><<>><<<>><>><<<>><<>>><<>>>><<<>>><<>>>><>>>><<<<>><<<><<><<<>>><<<><<<<>>>><<<<><<<>>>><<<>>>><<<>>>><>>>><>><<<<>><<<<>>>><<<<>><<<>><><>>>><<>>>><<>><<<>>>><<<<>><<<<>><<>>>><<>><<<<>>><<<<><<<<>>><>>><><<<>>><>><<>>>><<<<><<>><<>>><<<<>>>><<>>><>>>><<><<>>><<>>><<><<>><><<<<>><<<>>><>><<<>>><><<<<>>>><>>>><<<>>><<<>><<<>>><<<<>>><>>>><<<>>><<<>><<<>>>><<>>>><>><>><<<>><<<<><<<<>><<<<>>><>><<<<>><<>><><<>>><<<>>>><<>>>><<<>>><<>>><>><>><<><<<>>><<<><<<>><><><<<>><>>>><<>>>><<<<>>><<<><>>><<<><<<<>><<<><<>>><<<<>><><<<>>>><<>>>><>>>><<>>>><<<<>>><><>>>><>><>>><<<<><<<>><<<<><<<>>>><<<<>>>><>>>><>>><<<<><<<>>>><<<>><<<><<<><<<<><<<<>>>><<><<>>>><<<>><<<<>>><<<>>><><<<><<<>>>><<<<>>><<<><>>>><<<<>><<<<>>>><>><>>>><<<<>>>><<<>>>><<>>><<<>>>><>>><<<>>><<<>>>><>>>><<<>>><<<>><><<<>>>><<<>><<<<>>><>>><<<>>>><<<<>>>><<<><<<><<>><>>><<<>>><<<><<<>><<><<<<>>>><<<<><<>>><>><>><><<<>>>><<<<>><<<<>><<<<>>>><><>>><<<>><<>>><<<<>>>><<<><<<><<<>>>><<<><<>>><<<<>>>><<><>>>><<<>>><<<<>>><<<<>><<<>><<<<><<<<>><>>><>>>><>><<<>>>><><<<<>>><>>><<>>>><<>>><<<<>>>><<<><<><<<<>>><><<>>><>>><<<>><<<>><<>>><<<>>>><<<<>>><<<<><<>>>><<<>>><><>>>><>>><><<<<><<>>>><><<<>>><<>>><><>>><<<<>>>><<<<>>>><<<>>><<<<>>><<<>>><><>><<<><<>>><<>>><<>>><<<<>>>><><<<>>>><<>><<>>>><><<<>>><<><<>>><<<<>>>><<<>><<<><<><>>>><<<>><>>>><<<<>>>><<>>>><>>><>>><<<>><<<<>>><<<><<<<>><<<<>><<<>><><<>>>><<><<<>>>><<<<>><<<<>>><<<>>>><<<>>><<<><<<<>><<<<>>>><<><<>>><>>><<<>><>><<<<>>>><<<>>>><<<><<>>>><<>>><<<<><<<<>>><<<>>><<>><<<<>><<<<><<<>><<<>><>><<>><>>><<<<>><<>>>><<>><<>><<>><<<<>>>><<>><<<<><<<>><><><<>><<<><>>>><<<>>><>>><><>><<<>>>><<>><<>><<>>><<<<><<><>><<<<>>>><>>>><<>>><><<<<>><<<<><>><<<>><<><>>><<>>><<>>>><<<>><<<>>>><<<<>>><<><<<>>><<>>>><<<>>><<<<>><<>>><<<<>><>>><<<<><<<>>>><>><<><><<<<>>>><>>>><<>>>><<<>>><<<>><>>><<>>><<<>><><<<>>>><>>><>>><>><<<<>>>><<><<<<><<>><<>><<><>><<>>><<<>>><<<>>>><<<<>><><<><><<<><<<>>><<>><>>><<>>><<<<>><<<>><<>>>><>>>><<<>>>><>>><<<>>><<<<><<>>><<<<>>><<>><<>>>><>>><<>><<>>><<<><<>>><<>>><<<>>><<<>><>><<<>><<<<>>><>>>><<<<>>><>>><<<>><<<<>>>><>>><<><<<<>><><<><<><>>><><>><<<<>><<<<><>>>><<>>><<<<>>><>><<>><>>><<<<>><<>>>><><<<<>>>><<<<>><<<>>><><<><<>>>><<>><>>><<>><>><><<><<<>><<<<>>><><<>><>><<>>>><<>><>><<<<>>>><<<<>><<><>><<<<>>><<<>>>><<<>><<><<>>><<>>><<<>>>><<><<<<>>><<><<<<>>>><<<<><<<>>><><<>>><<<<><>><<>><<<<>>><<<<><<<>>>><<<>>><<<><<<>><>><<><<<>>>><<>><<><>>>><<>>>><<<<><<<>>>><<>>>><<><<<><<<>>>><<<<>>>><<><<<><<>><<<<>>><<<>>><<>>><<<>>>><>>>><>>><>><<<>>>><<<>>><<><<<>><<<><<<><<<<><>><>>><<<>>><<<<>>><<<><<<>><>><<>>><<>>>><<<<>>>><<>>><<<><<<>>><<<<>><><<>>>><<>>>><<<<>>>><<<>><<<>>>><>>><<<>><><<<<>><<>>><<>>><>>><<<<>>>><<><<<<>>>><<<<><<<<>>>><><<>><<<<>>>><<<<><<<>><<<><>>>><<<><<<>>>><<>>><><<>>><<>>>><<<<>><<<>><<<<>><<<<>><>><<<<>>><>><<<<><<><>>><><<<<>>><>><>>><<<>><<<>><>>>><>>>><<>><<<<>><<<>><<<>>><<<>>>><><<><<><<<<><<>>><<>>>><<<>><<<>><>>><<<<>><<><>>><>>><<<<>>>><>>>><>><<<>><<>><<<>><<<<>>>><<<>>>><<>>><<>><<<><<<>>>><<<>><>>><<><>>>><<<>>>><>><>>><><>>>><<<<>><>><<<>>><<<><>><<<<>>><>>>><<<><<<<><<<<>><<>>><><<<>>><<>>><<<>>>><<<<><<<<><>>><<<<>><>>><<<<>>>><>>><<<<>>><>>>><>>>><<<<>>><<<<>>>><<<<>>><>>><<>>>><<<<>><>>><<<>><<<<>>>><<<<>>>><<<><><<<><><<><>><<<><<>>><<<><<>>><<<<>>>><<>>>><><<>>>><<>><<<<>><>>>><<<<>>><<<>>>><>><<<<>>>><<<>>>><<<>>><<<<><<>><<<>>>><>>>><>>>><<>><<>>>><><>>><<><<>><<<<>>><><<<<><<>>><<<>>><<<>>>><<<<>><<<><>><<<>>>><<<>>><<<><<<<>>>><<<<><><<<>>>><<<>>><<<<>>>><>><<<<>>><<<>>>><<>>><>><<><><>>><<<>>>><<<<><<>>>><<<>>><<<>>>><<>><<<>>><<<<>>>><<>>>><<>>>><>><<<>>>><<>>>><<>>><<<>>><<>><<><<<<><<>>>><<<<>>><><<<><<<<>>><<<>>>><<>>>><><<>>><<>>><<<>>>><><>>><<<><<<>><<<<><>>><<<<>><<<>>>><<<>>>><<<>><<>><<<>><<<<>><<>><<>>>><<<<><<<>><>><><<><>>>><<<<>>><<<>>>><<<<>>><<<><<>><<<<>>><<<<>>>><<<<>>>><<<>>>><<<<>><<<<>>><<>>><>>><<<><<<<>>><>>><<><<<><<<<>><<<>>>><<<>>>><<<<>><<<>><>>><<<<>>>><><><<>>>><<><<<>>><<<>><<>><>>><<>><<<>><<<>>>><<<>>><<>><<<>><<<<>>>><<<<>>>><<<>>>><><><>>><<>>><<>>>><<>>><<<<>>>><<>>><<>><<<>>><<<>>>><<>>><>>>><>><<>>><<<>><<<><>>>><>>><<>>>><>><<<<><<<>>>><<<>><<<>><<<><><><<<<><>>><<<>>>><<<<>><<<<>>><>><<<>>><<<><<<<>>><<>><<>>>><><>><<><<><<<>><<><>>><<><<>>><<<<>>>><<<<><<<>><<><<<<>><<<<>>><<>>><<>>>><<<<>>>><<<>>>><<>>><<>>>><<<<>><>>><<><<><<>><<<><>><<<<><<<>>>><<>>><<<>>><<>>>><<<<>>>><<>>><<<<>><<<>>>><<<>>>><<<<>>>><<<>>>><<<<>>>><>><<<<>>><>>>><>>><>><<>>><<<<>>>><<<>>>><>>>><<<><<<<>><<<>>>><><<>>>><>><<<>>><<<>>>><<><<>>>><>>><>>><>>><>><<<><>>>><<<<>><<>><<<>><>>>><<<<>><<<>>>><<<>><<<<><<<<>>><<<<>>><<<<>>><<<>><<<<>>>><<<<><<<<>>>><>>>><<<><<<<>>>><<>><<<>>><<<>>>><>>><<<<>>>><<>>><>><<><<<>><<<>><<<>>><>>>><<<<>>><<>>>><<<<><<<<>>>><<<<>>>><<<<>>>><<<>><<>><>>><<<<>>>><<<<>>><><<<<>>>><>>><<>>>><<>>>><><<<<>>>><>><>>>><<>>>><>><<>>><<<>><<>>><<><<<<>>>><<<<>><><<<>>><<<>><>><<<<>><>>><<<<>>>><<<<>>>><><><<<>>>><<<>>><<<<><<<>>>><>>><<>>><<<<><>>>><>>><<<<>>>><>>><<<<>>>><<<<><<<<>>><<<<>>>><<<>>>><<<>><<<<>>>><>><<<<>>><>>><<<<>>>><>><>><<<<>><<<>>><><>>><<><<<<><<>>>><>>><>>><>><>>><<<><><>>>><>>>><<>><<<<>><<<>><<<<>><>>>><<<<>>>><<<>>><<>>>><><<>><<<>><<<>>>><<><>>><<>><<>>>><<<>>><<>>><<>>><>><<<>>><<>>><<<<>>><>>>><<<><><<<<>><<<<>>><>>>><><>>>><<<>>>><<<>>><<<<>>>><<<<>>>><<>><<<>><>>><>>>><<<<>><<>>><<>><<<>>>><<<>><<<<><<>>>><>><>><<<<>>><<<><<>>>><<<<>>><<<><>>><<<><<<>>>><<><>>>><<<<><>>><><<<<>>>><<>>>><>>>><<>>>><>>>><>><<<<>>>><>><>>>><<>>><<<>>>><<<>><><<<<><<<<>>>><<>>>><<<<>>><<>>>><<>><<<>>>><<>><>>><<<<>>><>>>><><<<><<<>>><<<<><><<<>>><><<>>><<><<<<>>>><>>>><<<>>>><<>>><>>>><<<<>>>><<<>><>><<<><>><<<>>><>>><>>><><<<>>>><><<<<>>><<<<>>>><<<<>>><<<<>>><>>>><<<<>>><<>><<<<>><>>><<><<<><<<><<<<>>>><<<><<<<>>><<>>>><<>><<<<>>>><<<<>>>><>>><<<<>>>><>>>><<>>><><<<>><>><>>>><><>><<<<>><<<<>>>><<><<<<><<>>><>>>><<><<<<>><><<<>>><<<><>>><<><<>>><<<<>>><<<>>><<<<><<<>>>><<<<><<>>>><<<><<>>><<<><><<<<>>>><<<>>><>>><<><<<>>><<<>><<>>><>><<>><>>><>><>><>>>><<<>><<<<><>>>><<<>>><<<>><<<><><<>><<<>>><<<<>>>><>><<<><<<><<<<><<>><>><<<<>>>><<<>>><<><<<>>><<<><>><>>><>>>><><>>><<<>><><<<<>><<<><<<<>>><<<<>>><<<>>>><<>><<>>>><<><<<>>><<<<>><<>>>><>>>><<<<>><<>><<<>>>><<<>>><<<>><<<<><<<><<>>>><<<>>>><<<>>>><<>>><<<<>>><<>>>><<>>><><>>>><<<>>>><<<><<>>><<<><<<<>><<<>><<>>>><<>>>><<<>>>><<>><>>>><<<<>><>><<>>>><<<<>>><<<<>>>><<<<>>><<><<<>>>><<<<>>><<>><<><<><<>><<<<><<<<>>><>>><<>>><><<<<>><>>><<<<><<<<><<<>>>><<><<<>>>><<<<>>>><>>><>><><><<<<><>><<<<>><<<<>>><<>>>><<<><<<<><<<<>>><<>><<>>>><>><<>><<><<>>>><>>>><><<<<>>>><<<><>>>><<<<>><<>>><<<<><<<<>>><<<<><>><<<<>>>><<>>><<>>>><<<><<<>><>>><<<>>>><>><<<>><><<<>>>><><>>>><>>><<<>>>><>>>><>>>><<<><<>>><<><<<<>>><>><<><<<><<<>>><>>><>>><<>><<><>><><>>><<><<<<><<>>>><<<>>><>>><<<<>>>><<<>>><<<><<<>>><<>><<<>><>><<>><<<><<<<>><<>>>><<<<>>>><<><<>>>><<<<><>>>><<<<><><>>><<<>><<<<><<>>>><<<>><<<>>><<<<><>><<<><<>>>><<<<>>><<<>>><>>>><<><<>>><<>><<<><<<<>>><>>><<>><>>>><>><<<>>>><><<<<>>>><<<>>>><>><<<>>><>>><<>>>><<<<>>>><<<>>><<><<<<><<<<>>><>>><><<<<>>>><<<><<<>>>><<<<>><<>>><<<<>>>><<>>><<<<>>><>><<<>>>><>>><<<<><<<<>>>><<<<>><<><>>>><<<<>>><<<<>><<>>>><><<<>><>>>><<<><>><<<>>>><>><<>>>><<<>>><>><<<>><>><<<><<>>><<<<>>><<<><<<>><><<<>><><<<<>>><<<<><<<>><<>>><<>><<<><<<><<<<>>>><<<><<<><<<<><>>><><<<>>><<<<>>><<<>>><<<>>>><<<<>>><<<><<>><>>>><<>><<<>><<<<>>><>><<>><<>>>><>><>>><<<<><<<><>><>>><><<>><<<><<>>>><>><<<>><<<<><<<>>>><<>><<<<>><<<<><<<>>><<<<><<<<>><<>><<<>>>><<>>>><<<>>><<><<>>>><<<><<><<<<><<<>><<>>><<<<>><<>><<<>>>><<>>>><<<<>>><<<<>>><<<<>><>><<<<>>>><<>>>><<<<><<>>>><<><><><<<<>><<<>><><<<<>>><>>>><<><<<<>><<<<>>><><<<<>><<<>><<>><<<<><>><<<>>>><<<><<<<>>>><>>>><<<<>>><<<>>><<<<><<<>>><<<<><<>><<>>>><>><<<<><<<<>>>><<<<>>>><<>>>><<<>>><<>><<<>>><<<<>>>><<<<><>>><<>><<<<>><><<<<>>><>>><<<<>><>><<><<>>><>>><<<<>><><<<>><>><<<>><<<<>><<>><<<>><>><<>><<<>><>>>><<<<>>>><>>>><<<>><<<>>><<>><>><<<>>>><<<<>>>><>><>><<><<<<>>>><<><>>>><>>><<<>>>><<<<><<<<>><><<<<>>>><>>>><<>><<<<><<>>><<<<><>><><<<>>><<>>><>>>><><<<<><<<<>>><>>><>>>><<<>>>><>>>><<><<<>>><<>><<<>>><>><>><<<<>><<>>>><<<><><<<<><<>>><<<<><<<><<<>><>>><<<>>>><>><<<><<<<>>><<<<>>><<>><<>><>>>><<<<>>>><>><>>><>>>><<<<>><<<<>>><<<><<<>>><><>>>><<<>>><>>>><<<>>>><<<<>><<<<><<<<>>>><<>>>><>>><<<>>><<<>><<<<>><<<>><>><><>>><<><>>>><<<<>><<>><<>>><<>>>><<<<>>>><<<<><><<<<>><<<>>><<><<><<<<>>>><<>>>><<<>>>><>>><<>>>><><<<<>><<>>><<>>>><<><<>><>><<<>>><<<>>>><>><>><>>><<<<><<<>><<<<>><<<<>>><><>>><>><<>><>><>>><><<>>><>>><<<>>><>>><>><<<<>>>><<>>><>><<<>>>><><<>>>><<<>>>><>>>><<<>><<<>><<>>><<<>><<<<><<<><<<<>><<>>><>>><<<>>>><<<>><<<>><<<><<>><<<<><<<>><<<<><<<<><<<<><<<<><>>>><<<>>>><<<<>>>><<<><>>><<>><<<<>>>><>><<<<><>>>><<<<>><<><<<><>><>>>><<<><>>>><<<><<<><<<<>><<<><<<<><<<>>><<>>>><<>>><<<<><<<<><><<<<><<<>><>>><><>>><><<<<>>><<>>>><<<<>>><<<<>>>><<<>><>>><<<<>><>>><<<><<<><<>><<<><<>><<>>><<><>>><<<<><>>>><<<>><<<>>>><<<<>>><<<<>>><<<<>>><<>>>><>><>><>><<>><<<>>><>>><<><<>>><>>><<<>>>><<<><>><>><<>>><<<>><>>>><<<<>><<<<>><<<<><<<>><<>>><>>>><<<<>><>>><<<<><<<<>><>>><<>>><<<<>>><<<>>><><<<<>`
