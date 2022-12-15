package main

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var re = regexp.MustCompile(`([-\d]+)`)

func main() {
	fmt.Println(part1(input1, 10))           // 26
	fmt.Println(part1(input2, 2_000_000))    // 6124805
	fmt.Println(part2(input1, 0, 20))        // 56000011
	fmt.Println(part2(input2, 0, 4_000_000)) // 12555527364986
}

func part1(input string, y int) int {
	data := parse(input)
	result := make(map[int]bool)
	for _, d := range data {
		if !d.contains(y) {
			continue
		}

		n := d.dist - abs(d.sensor.y-y)
		for x := d.sensor.x - n; x <= d.sensor.x+n; x++ {
			if (d.beacon == point{x: x, y: y}) {
				continue
			}
			result[x] = true
		}
	}

	return len(result)
}

func part2(input string, minD, maxD int) int {
	data := parse(input)

	// Sort to ensure we are searching from left to right.
	sort.Slice(data, func(i, j int) bool {
		lhs, rhs := data[i].sensor, data[j].sensor
		if lhs.x == rhs.x {
			return lhs.y < rhs.y
		}

		return lhs.x < rhs.x
	})
	//grid := make(map[point]bool)

	for y := minD; y <= maxD; y++ {
		var intervals []interval
		for _, d := range data {
			if !d.contains(y) {
				continue
			}

			n := d.dist - abs(d.sensor.y-y)
			curr := interval{
				min: d.sensor.x - n,
				max: d.sensor.x + n,
			}
			// DEBUG
			//for x := d.sensor.x - n; x <= d.sensor.x+n; x++ {
			//grid[point{x: x, y: y}] = true
			//}

			intervals = append(intervals, curr)
		}

		intervals = reduceIntervals(intervals)
		if len(intervals) != 2 {
			continue
		}
		if intervals[1].min-intervals[0].max > 2 || intervals[0].max+1 < minD || intervals[0].max+1 > maxD {
			continue
		}

		b := point{x: intervals[0].max + 1, y: y}
		return tuningFrequency(b.x, b.y)
	}
	//draw(grid, data)

	return -1
}

func reduceIntervals(intervals []interval) []interval {
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].min < intervals[j].min
	})
	reduced, intervals := intervals[:1], intervals[1:]

	for len(intervals) > 0 {
		var left, right interval
		reduced, left = reduced[:len(reduced)-1], reduced[len(reduced)-1]
		right, intervals = intervals[0], intervals[1:]
		if left.adjacent(right) || left.intersect(right) {
			reduced = append(reduced, interval{
				min: min(left.min, right.min),
				max: max(left.max, right.max),
			})
		} else {
			reduced = append(reduced, left, right)
		}
	}

	return reduced
}

func draw(grid map[point]bool, data []data) {
	minX, maxX, minY, maxY := math.MaxInt, 0, math.MaxInt, 0
	for _, p := range data {
		s, b := p.sensor, p.beacon
		if s.x < minX {
			minX = s.x
		}
		if s.x > maxX {
			maxX = s.x
		}
		if s.y < minY {
			minY = s.y
		}
		if s.y > maxY {
			maxY = s.y
		}

		if b.x < minX {
			minX = b.x
		}
		if b.x > maxX {
			maxX = b.x
		}
		if b.y < minY {
			minY = b.y
		}
		if b.y > maxY {
			maxY = b.y
		}
	}
	for p := range grid {
		if p.x < minX {
			minX = p.x
		}
		if p.x > maxX {
			maxX = p.x
		}
		if p.y < minY {
			minY = p.y
		}
		if p.y > maxY {
			maxY = p.y
		}
	}
	dx := maxX - minX
	dy := maxY - minY
	output := make([][]rune, dy+1)
	for i := range output {
		output[i] = make([]rune, dx+1)
		for j := range output[i] {
			output[i][j] = '.'
		}
	}
	for p := range grid {
		output[p.y-minY][p.x-minX] = '#'
	}
	for _, d := range data {
		s, b := d.sensor, d.beacon
		output[s.y-minY][s.x-minX] = 'S'
		output[b.y-minY][b.x-minX] = 'B'
	}
	for i := range output {
		fmt.Println(string(output[i]))
	}
	fmt.Println()
}

func tuningFrequency(x, y int) int {
	return x*4_000_000 + y
}

type interval struct {
	min int
	max int
}

func (i interval) intersect(o interval) bool {
	if o.min < i.min {
		return o.intersect(i)
	}

	return o.min <= i.max && i.min <= o.max
}

func (i interval) adjacent(o interval) bool {
	return o.min-i.max == 1
}

type data struct {
	sensor point
	beacon point
	dist   int
}

func (d data) contains(y int) bool {
	return abs(d.sensor.y-y) <= d.dist
}

type point struct {
	x, y int
}

func (p point) manhattan(o point) int {
	return abs(p.x-o.x) + abs(p.y-o.y)
}

func parse(input string) []data {
	rows := strings.Split(input, "\n")
	res := make([]data, len(rows))
	for i, row := range rows {
		m := re.FindAllString(row, -1)
		s := point{
			x: toInt(m[0]),
			y: toInt(m[1]),
		}
		b := point{
			x: toInt(m[2]),
			y: toInt(m[3]),
		}

		res[i] = data{
			sensor: s,
			beacon: b,
			dist:   s.manhattan(b),
		}
	}

	return res
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func toInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

var input1 = `Sensor at x=2, y=18: closest beacon is at x=-2, y=15
Sensor at x=9, y=16: closest beacon is at x=10, y=16
Sensor at x=13, y=2: closest beacon is at x=15, y=3
Sensor at x=12, y=14: closest beacon is at x=10, y=16
Sensor at x=10, y=20: closest beacon is at x=10, y=16
Sensor at x=14, y=17: closest beacon is at x=10, y=16
Sensor at x=8, y=7: closest beacon is at x=2, y=10
Sensor at x=2, y=0: closest beacon is at x=2, y=10
Sensor at x=0, y=11: closest beacon is at x=2, y=10
Sensor at x=20, y=14: closest beacon is at x=25, y=17
Sensor at x=17, y=20: closest beacon is at x=21, y=22
Sensor at x=16, y=7: closest beacon is at x=15, y=3
Sensor at x=14, y=3: closest beacon is at x=15, y=3
Sensor at x=20, y=1: closest beacon is at x=15, y=3`
var input2 = `Sensor at x=3907621, y=2895218: closest beacon is at x=3790542, y=2949630
Sensor at x=1701067, y=3075142: closest beacon is at x=2275951, y=3717327
Sensor at x=3532369, y=884718: closest beacon is at x=2733699, y=2000000
Sensor at x=2362427, y=41763: closest beacon is at x=2999439, y=-958188
Sensor at x=398408, y=3688691: closest beacon is at x=2275951, y=3717327
Sensor at x=1727615, y=1744968: closest beacon is at x=2733699, y=2000000
Sensor at x=2778183, y=3611924: closest beacon is at x=2275951, y=3717327
Sensor at x=2452818, y=2533012: closest beacon is at x=2733699, y=2000000
Sensor at x=88162, y=2057063: closest beacon is at x=-109096, y=390805
Sensor at x=2985370, y=2315046: closest beacon is at x=2733699, y=2000000
Sensor at x=2758780, y=3000106: closest beacon is at x=3279264, y=2775610
Sensor at x=3501114, y=3193710: closest beacon is at x=3790542, y=2949630
Sensor at x=313171, y=1016326: closest beacon is at x=-109096, y=390805
Sensor at x=3997998, y=3576392: closest beacon is at x=3691556, y=3980872
Sensor at x=84142, y=102550: closest beacon is at x=-109096, y=390805
Sensor at x=3768533, y=3985372: closest beacon is at x=3691556, y=3980872
Sensor at x=2999744, y=3998031: closest beacon is at x=3691556, y=3980872
Sensor at x=3380504, y=2720962: closest beacon is at x=3279264, y=2775610
Sensor at x=3357940, y=3730208: closest beacon is at x=3691556, y=3980872
Sensor at x=1242851, y=838744: closest beacon is at x=-109096, y=390805
Sensor at x=3991401, y=2367688: closest beacon is at x=3790542, y=2949630
Sensor at x=3292286, y=2624894: closest beacon is at x=3279264, y=2775610
Sensor at x=2194423, y=3990859: closest beacon is at x=2275951, y=3717327`
