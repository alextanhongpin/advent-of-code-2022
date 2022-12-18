package main

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

func main() {
	fmt.Println(part1(input1)) // 64
	fmt.Println(part1(input2)) // 4308
	fmt.Println(part2(input1)) // 58
	fmt.Println(part2(input2)) // 2540
}

func part1(input string) int {
	cache := make(map[string]int)
	rows := strings.Split(input, "\n")
	set := func(key string, x, y, z string) {
		cache[fmt.Sprintf("%s:%s:%s:%d", key, x, y, toInt(z)-1)]++
		cache[fmt.Sprintf("%s:%s:%s:%s", key, x, y, z)]++
	}

	for _, row := range rows {
		xyz := strings.Split(row, ",")
		x, y, z := xyz[0], xyz[1], xyz[2]
		set("xy", x, y, z)
		set("yz", y, z, x)
		set("xz", x, z, y)
	}

	var c int
	for _, n := range cache {
		if n != 1 {
			continue
		}
		c++
	}

	return c
}

type point struct {
	x, y, z int
}

func (p point) String() string {
	return fmt.Sprintf("%d,%d,%d", p.x, p.y, p.z)
}

func part2(input string) int {
	rows := strings.Split(input, "\n")

	cubes := make(map[point]bool)
	minX, maxX := math.MaxInt, math.MinInt
	minY, maxY := math.MaxInt, math.MinInt
	minZ, maxZ := math.MaxInt, math.MinInt
	for _, row := range rows {
		xyz := strings.Split(row, ",")
		x, y, z := xyz[0], xyz[1], xyz[2]
		p := point{
			x: toInt(x),
			y: toInt(y),
			z: toInt(z),
		}
		cubes[p] = true

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
		if p.z < minZ {
			minZ = p.z
		}
		if p.z > maxZ {
			maxZ = p.z
		}
	}

	isOOB := func(p point) bool {
		return false ||
			p.x < minX ||
			p.x > maxX ||
			p.y < minY ||
			p.y > maxY ||
			p.z < minZ ||
			p.z > maxZ
	}

	filled := make(map[point]bool)
	var spaceSurroundedByCubes func(e point, exclude map[point]bool) (map[point]bool, bool)
	// e represents the empty point.
	// All 6 sides must be a cube, otherwise, it could be a space that is
	// surrounded by cubes.
	spaceSurroundedByCubes = func(e point, exclude map[point]bool) (map[point]bool, bool) {
		spaces := make(map[point]bool)
		sides := adj(e)

		for _, side := range sides {
			if exclude[side] {
				continue
			}

			if cubes[side] {
				continue
			}

			if isOOB(side) {
				return nil, false
			}

			exclude[side] = true
			neighbors, ok := spaceSurroundedByCubes(side, exclude)
			if !ok {
				return nil, false
			}

			for k, v := range neighbors {
				spaces[k] = v
			}
		}

		spaces[e] = true

		return spaces, true
	}

	found := make(map[string]bool)
	for c := range cubes {
		sides := adj(c)
		for _, side := range sides {
			if filled[side] {
				continue
			}
			if cubes[side] || isOOB(side) {
				continue
			}

			spaces, ok := spaceSurroundedByCubes(side, make(map[point]bool))
			if !ok {
				break
			}
			var res []string
			for s := range spaces {
				filled[s] = true
				res = append(res, s.String())
			}
			sort.Strings(res)
			key := strings.Join(res, "\n")
			found[key] = true
		}
	}
	var exclude int
	for in := range found {
		exclude += part1(in)
	}

	return part1(input) - exclude
}

func adj(p point) []point {
	x, y, z := p.x, p.y, p.z
	return []point{
		{x + 1, y, z},
		{x - 1, y, z},
		{x, y + 1, z},
		{x, y - 1, z},
		{x, y, z + 1},
		{x, y, z - 1},
	}
}

func copyMap[K comparable, V any](m map[K]V) map[K]V {
	res := make(map[K]V)
	for k, v := range m {
		res[k] = v
	}
	return res
}

func toInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

var input1 = `2,2,2
1,2,2
3,2,2
2,1,2
2,3,2
2,2,1
2,2,3
2,2,4
2,2,6
1,2,5
3,2,5
2,1,5
2,3,5`
var input2 = `16,8,6
19,11,14
10,12,4
7,4,14
9,4,5
11,4,17
14,17,13
4,10,15
10,4,15
8,12,18
16,16,6
10,10,3
9,8,20
11,1,10
3,4,9
7,10,5
3,13,5
9,5,18
7,13,6
15,17,11
2,11,13
15,5,6
15,15,7
4,8,15
1,14,9
4,15,9
17,13,8
15,19,11
17,12,5
16,12,2
6,15,15
7,4,11
9,19,6
12,18,14
16,9,3
9,3,17
18,10,8
10,5,16
20,9,12
4,15,14
17,7,13
10,6,16
3,4,12
13,10,19
5,4,8
14,15,16
17,11,3
6,12,18
12,9,19
15,17,15
10,2,12
11,6,3
3,10,4
8,11,19
9,2,9
4,11,5
6,4,11
14,13,16
11,8,18
16,11,5
9,2,14
13,8,16
17,9,5
15,5,7
18,6,9
6,5,8
18,13,12
12,4,4
15,6,7
17,13,6
15,3,6
9,12,20
7,7,4
4,5,7
12,12,2
17,9,17
7,5,5
3,14,8
13,18,11
16,17,7
6,6,15
16,7,4
5,13,4
7,19,6
8,15,18
10,13,19
12,19,6
3,14,17
17,11,16
10,5,3
13,3,14
6,8,6
10,19,13
15,15,9
3,11,6
2,11,8
6,13,17
15,17,13
2,8,11
6,7,6
18,9,9
12,5,15
10,3,8
7,14,4
10,4,4
6,11,2
7,11,2
14,4,13
20,10,13
4,15,8
15,4,14
9,20,13
7,4,6
8,10,17
11,14,18
11,11,19
3,15,10
6,8,17
6,11,4
6,17,8
13,3,7
7,6,5
12,5,19
5,7,13
10,19,14
13,2,9
16,9,6
3,13,11
14,4,14
4,8,4
5,12,16
2,12,8
17,14,15
10,12,19
3,7,10
12,8,19
7,19,9
6,18,6
10,6,3
12,1,10
5,8,4
16,5,10
8,14,19
16,13,4
13,2,10
15,14,15
16,15,16
4,7,10
11,5,15
15,6,14
6,20,11
7,16,18
12,3,15
6,13,18
7,10,2
14,9,17
19,9,9
13,4,15
7,2,10
7,3,14
12,13,2
17,11,5
2,10,6
18,12,15
11,14,20
4,11,17
19,15,11
16,4,14
5,9,17
16,18,11
4,6,16
16,12,4
11,6,17
17,11,17
9,3,13
12,13,4
15,19,10
7,3,11
15,17,6
11,5,5
6,4,15
8,11,17
4,8,16
6,14,18
10,20,8
10,19,11
4,8,7
2,7,7
15,16,17
11,16,18
12,8,18
7,5,6
7,3,7
14,11,19
13,12,18
10,18,11
20,10,12
14,6,6
1,14,10
11,17,18
13,5,17
17,5,9
5,15,9
8,3,12
7,18,15
18,15,11
13,5,18
15,16,6
3,17,12
17,7,18
13,17,6
3,8,6
18,14,10
13,14,4
18,13,6
12,19,13
14,17,14
16,3,9
6,16,17
6,10,3
18,14,15
8,20,12
16,16,17
11,7,19
19,10,5
18,11,10
2,14,10
11,2,14
15,6,5
14,17,17
14,15,17
10,18,15
16,5,13
3,7,8
5,11,19
6,3,9
14,18,6
11,4,6
18,15,12
20,9,10
18,14,8
11,3,9
10,16,5
12,12,5
13,17,7
7,12,2
11,7,16
19,16,11
16,9,17
16,8,4
2,13,13
16,15,10
11,18,11
12,16,5
8,13,18
14,5,8
19,8,6
14,13,18
14,16,16
9,10,2
5,9,6
5,7,6
8,6,4
11,18,15
7,14,17
8,2,13
14,15,18
15,7,13
7,2,7
5,14,9
18,10,7
10,6,18
3,6,10
1,13,13
8,3,6
10,10,19
8,4,14
5,5,9
16,15,5
6,2,9
17,6,14
11,5,4
10,18,12
9,14,4
6,7,18
18,13,8
13,7,18
11,13,20
11,3,13
19,8,12
17,12,6
6,12,3
8,10,2
19,6,10
7,18,11
8,11,3
12,9,18
13,12,20
19,13,8
9,2,7
11,8,19
16,8,5
7,2,9
17,12,8
19,9,13
15,9,18
14,10,18
13,10,3
15,16,15
10,7,1
19,10,9
9,0,9
12,17,18
4,4,13
18,13,5
3,9,13
9,18,9
16,17,9
12,16,4
13,14,18
19,11,9
2,11,12
13,8,18
6,18,10
8,18,10
17,17,14
12,18,7
14,14,5
12,11,2
18,10,10
17,6,8
7,12,19
16,12,17
11,7,3
7,9,3
6,11,16
9,3,11
17,4,9
4,14,13
3,9,7
4,11,16
13,9,3
8,19,10
4,7,16
15,17,10
12,3,7
14,7,4
11,9,19
13,16,17
14,17,7
15,14,16
8,11,20
7,17,6
4,8,5
18,5,8
14,15,19
7,8,17
9,6,15
1,11,11
14,6,3
6,3,8
13,7,2
18,13,9
11,12,2
7,18,10
13,5,15
19,12,8
18,5,7
5,16,15
10,16,18
18,11,13
13,4,8
18,11,9
7,4,9
1,13,11
15,18,11
14,16,17
13,14,3
14,3,14
17,8,6
6,11,3
10,2,11
18,13,7
11,17,7
6,2,11
6,7,16
5,6,5
17,6,15
16,14,6
11,2,12
7,15,4
12,7,19
9,2,11
12,18,16
4,17,12
11,3,16
14,5,7
14,1,9
7,16,5
17,17,15
13,3,9
17,13,14
8,6,16
19,13,9
11,2,9
16,6,16
12,5,6
9,17,18
9,2,15
18,8,10
4,8,18
4,16,12
5,5,7
9,18,8
8,12,1
15,9,5
14,17,16
3,17,10
11,18,12
11,5,3
2,8,12
19,6,11
3,16,7
5,11,18
7,4,12
16,4,7
4,4,12
13,2,12
2,11,14
12,19,14
16,15,4
10,2,10
16,13,17
15,2,10
11,15,2
9,20,9
7,17,13
14,18,11
3,14,14
18,11,7
16,15,14
3,14,10
8,3,11
1,9,8
6,10,19
15,13,4
17,10,18
9,7,18
8,2,7
18,16,11
15,7,14
5,9,18
10,2,7
16,12,16
10,18,6
6,4,14
4,12,16
16,12,18
9,17,4
3,15,11
15,14,5
11,18,16
4,14,6
7,19,13
10,5,18
5,16,9
11,15,6
9,19,10
2,11,7
17,16,11
6,16,5
5,14,6
14,4,11
15,13,3
11,2,8
3,16,8
17,14,13
16,7,16
19,9,11
13,19,11
8,18,14
12,1,6
4,10,4
10,4,5
1,10,10
13,15,3
12,17,15
8,17,6
17,5,12
3,10,7
12,6,18
8,8,5
6,7,5
7,2,14
8,16,5
10,15,18
4,5,10
3,11,17
7,5,3
9,3,6
12,7,3
10,21,12
6,16,9
8,9,19
9,17,15
11,19,9
15,19,8
8,9,18
2,10,14
6,17,10
16,3,11
8,7,15
13,3,11
17,5,6
8,3,13
4,5,13
11,18,14
9,3,10
3,4,11
4,14,17
18,6,11
10,3,11
11,19,12
12,2,11
19,8,13
8,1,13
9,18,10
11,4,15
8,14,18
14,14,3
9,18,14
8,16,16
7,12,18
9,13,2
15,3,11
13,7,19
5,9,7
16,16,10
4,15,15
6,4,9
4,14,4
11,12,19
4,3,10
16,17,6
9,18,5
7,18,7
11,10,19
7,6,16
16,2,10
5,2,12
13,8,20
14,17,4
18,11,14
10,16,17
17,16,7
16,10,17
9,8,5
15,17,12
7,8,4
10,7,17
14,5,6
10,14,19
5,6,17
6,12,5
5,7,5
13,9,4
13,8,3
17,16,13
19,11,7
10,2,16
2,8,6
17,7,15
19,9,8
3,16,11
11,17,17
16,9,16
14,7,19
15,9,3
6,3,12
13,15,2
15,15,15
8,3,5
3,10,16
19,13,11
11,19,11
8,7,19
12,4,16
16,4,9
2,13,11
20,14,11
4,12,3
8,3,7
18,15,9
19,10,7
4,8,14
15,14,4
12,14,3
12,6,16
8,3,16
9,3,14
10,16,4
15,6,10
6,7,3
13,17,16
8,9,2
5,17,10
8,14,6
12,18,12
15,15,4
12,16,3
11,2,10
10,17,12
9,5,17
12,7,2
5,16,7
18,12,13
4,16,7
17,4,12
12,5,5
9,5,16
13,18,12
5,14,16
20,10,11
19,7,14
19,8,9
15,5,16
16,4,13
5,17,12
4,5,6
17,15,7
17,7,16
4,16,13
12,3,6
12,10,19
18,15,14
7,8,3
6,12,19
18,6,12
13,18,16
13,5,7
7,5,15
6,16,15
16,16,8
11,19,6
18,8,14
11,6,18
16,7,15
16,11,15
13,16,16
19,10,11
11,20,9
8,18,9
13,11,18
8,4,6
9,19,13
14,9,3
11,16,17
18,5,13
5,5,10
12,17,6
1,10,12
19,9,6
12,10,20
12,13,19
4,9,4
11,10,5
17,9,6
15,18,12
12,14,5
6,4,7
11,2,15
12,3,5
6,8,18
15,11,4
13,17,14
9,13,3
15,10,4
11,5,16
10,20,11
7,15,3
16,15,7
19,13,12
8,17,17
7,4,16
7,11,17
10,9,3
14,9,19
19,10,8
1,8,8
17,15,14
7,10,18
14,3,16
1,9,15
7,8,18
16,9,18
9,7,4
14,3,10
6,6,17
2,8,10
3,11,9
13,3,10
3,13,12
18,8,7
18,8,17
4,13,4
13,17,15
11,1,13
13,7,17
11,9,3
16,3,10
7,19,11
11,13,18
17,6,9
3,9,6
5,10,18
19,13,13
4,7,14
17,13,16
13,16,4
11,14,16
4,10,6
16,14,10
17,15,11
18,7,15
17,9,13
12,6,17
7,6,4
19,12,13
14,5,5
2,11,5
12,19,11
18,9,16
16,5,12
18,14,13
12,5,17
6,15,5
5,15,17
7,11,3
6,4,6
16,11,18
2,14,8
11,13,2
18,10,9
5,6,14
2,13,12
2,11,6
9,10,20
1,13,7
16,14,5
10,3,7
2,14,13
18,9,11
17,12,7
18,10,4
13,1,9
10,11,20
16,17,11
12,18,8
12,3,16
19,7,11
19,11,13
10,15,5
14,3,12
3,13,13
13,4,11
6,5,16
6,4,16
4,14,7
12,3,13
5,18,9
18,12,12
8,17,15
9,18,15
16,4,11
5,4,9
17,4,11
15,16,9
10,4,16
19,13,10
14,7,17
3,10,9
8,11,2
3,7,15
17,13,13
8,15,2
2,9,12
2,15,11
13,12,16
14,5,18
15,4,12
15,4,5
20,11,12
7,5,7
8,3,14
11,6,5
18,7,12
14,8,3
12,3,12
17,10,6
18,13,14
6,5,13
8,13,17
4,13,15
12,5,14
19,12,7
1,9,10
6,4,8
5,5,5
15,4,13
12,14,20
4,13,16
8,17,4
9,4,6
13,9,2
11,14,2
11,9,2
16,17,14
17,6,7
10,2,9
11,13,21
5,7,15
10,5,5
10,16,6
12,2,12
3,11,5
6,3,11
12,17,14
11,14,17
18,10,5
18,6,10
16,5,9
5,6,6
5,3,11
16,16,9
13,3,13
6,14,17
13,18,6
14,6,18
10,5,15
13,3,4
9,9,20
6,10,18
5,16,5
17,12,16
7,9,1
17,9,14
19,8,10
15,5,10
4,13,18
3,11,13
4,12,5
6,6,4
5,13,16
9,11,17
18,16,9
7,11,5
9,13,20
14,12,18
5,16,6
6,17,6
15,20,10
2,14,15
6,15,4
15,8,4
7,13,19
7,6,6
9,20,6
14,2,11
17,10,5
9,15,15
5,17,9
15,4,11
9,1,12
12,16,18
10,7,5
4,17,9
16,6,7
13,12,17
5,8,2
13,4,12
10,10,17
8,12,19
5,9,3
17,7,6
5,9,4
5,12,17
20,10,8
4,13,17
10,13,2
1,8,9
4,9,7
16,5,7
7,15,16
16,11,16
15,15,5
6,18,7
13,3,16
11,19,8
2,9,15
5,10,17
13,2,13
8,12,2
10,14,17
4,9,15
16,17,8
11,15,4
7,9,18
12,4,6
16,4,8
9,4,15
10,17,11
2,10,12
4,11,7
12,12,19
4,6,9
14,11,4
11,15,3
17,13,3
6,15,16
20,12,11
3,13,14
17,5,8
20,13,11
10,17,4
19,11,12
2,12,10
9,14,2
14,10,5
16,4,6
9,18,6
5,17,15
15,7,4
18,5,9
12,9,2
10,2,15
17,11,18
10,5,17
18,16,10
8,18,11
15,5,11
16,4,10
16,17,5
5,4,16
10,11,2
12,17,16
9,3,12
18,8,13
2,11,15
13,18,5
20,10,10
14,16,15
8,19,8
5,10,3
18,7,9
8,4,7
15,17,8
5,16,8
16,18,10
5,16,13
18,7,14
12,19,8
4,15,5
6,17,12
19,9,5
17,10,7
9,15,18
6,18,14
10,3,13
14,18,13
5,8,3
8,14,5
9,18,11
4,11,4
10,7,2
17,5,10
16,6,6
8,6,3
13,9,1
7,16,16
2,8,8
7,6,7
7,3,15
3,12,14
3,13,6
6,2,10
17,11,7
11,17,5
15,17,16
12,10,3
14,9,4
16,10,4
18,7,13
17,16,10
13,20,12
12,16,16
5,5,14
14,11,3
3,9,14
11,15,19
20,7,8
3,8,9
12,3,8
17,8,5
3,8,14
9,9,16
17,13,4
13,19,7
4,11,14
18,10,13
6,14,5
18,12,11
18,12,14
4,12,6
19,14,8
5,12,14
15,6,17
3,6,14
15,16,7
17,15,13
2,13,14
4,5,11
8,5,5
8,4,10
11,4,7
4,6,13
5,17,8
7,18,9
15,9,17
11,10,2
4,6,7
4,7,11
7,14,18
20,13,8
16,14,8
13,3,15
11,2,6
14,3,5
11,16,3
18,14,12
14,19,10
14,15,3
11,8,2
13,19,8
6,18,11
11,3,6
10,8,2
16,7,17
13,11,1
15,7,17
9,5,9
5,15,15
1,11,9
5,14,11
11,16,6
3,6,6
14,17,5
3,11,8
4,7,7
9,18,12
2,7,9
13,3,8
11,16,4
3,7,6
10,20,12
9,18,16
8,13,4
6,2,8
12,2,8
9,12,1
15,13,2
9,12,19
10,16,3
18,9,6
15,3,8
5,4,13
9,5,4
10,3,5
17,16,9
9,14,6
3,13,8
6,17,11
11,2,13
4,12,7
13,11,3
17,13,15
9,19,7
8,5,4
17,10,4
14,2,15
10,17,14
1,13,12
6,11,5
9,3,16
12,4,14
10,9,2
13,8,4
6,15,17
6,17,9
11,1,12
15,14,3
7,2,13
15,17,14
4,16,9
13,15,18
13,13,1
8,5,17
10,17,9
16,6,15
3,5,9
2,11,11
11,14,3
13,19,13
17,10,10
5,14,15
14,12,2
6,5,15
15,4,16
2,13,8
18,6,13
4,14,14
6,18,9
10,18,14
4,14,9
10,13,18
8,4,16
4,4,11
11,3,5
7,6,18
7,19,10
3,9,12
7,19,8
12,18,9
9,5,5
12,6,3
18,4,12
14,13,2
15,15,18
19,11,10
13,3,12
14,2,6
12,3,14
7,9,4
3,12,5
16,14,4
13,18,14
11,18,5
14,4,17
15,10,18
13,4,5
11,6,4
15,14,18
15,18,6
8,18,16
4,13,6
12,19,10
15,6,4
6,13,16
8,18,12
7,9,19
11,3,7
13,6,5
11,13,4
6,8,4
8,13,20
9,4,3
7,6,3
17,17,11
17,12,4
2,12,7
8,13,19
6,4,5
16,6,14
18,7,11
16,16,14
13,16,18
6,13,5
2,12,13
18,8,11
7,7,18
14,5,16
18,12,8
14,7,16
12,2,7
16,14,17
8,19,13
17,8,16
5,13,18
11,5,17
16,4,12
16,7,18
19,12,12
11,16,5
12,14,18
13,17,4
12,2,9
6,13,4
5,6,16
14,12,1
7,11,20
9,6,17
17,13,17
6,17,13
19,12,6
4,15,16
8,10,1
4,12,14
9,17,14
5,12,18
17,14,8
3,7,13
12,13,17
18,9,12
9,14,3
4,16,11
17,8,15
17,9,16
13,17,5
14,4,5
3,7,12
12,3,9
19,15,9
12,17,11
11,5,18
14,14,17
18,9,10
19,11,8
4,6,8
15,3,7
4,17,11
18,12,7
4,5,8
10,9,17
14,4,15
5,17,7
2,9,10
1,11,12
5,10,19
6,5,5
10,16,7
7,6,14
6,9,3
11,19,13
18,15,8
18,9,8
14,18,9
4,14,11
6,4,13
16,10,5
2,12,14
3,12,17
12,3,11
9,19,12
6,10,4
14,5,17
9,1,11
4,13,12
6,10,2
18,7,7
8,17,8
7,16,4
1,9,9
5,15,14
14,18,16
15,9,4
13,14,17
14,10,3
15,5,3
8,4,5
15,13,18
9,9,19
7,11,18
8,5,6
11,10,20
3,8,5
3,6,9
11,18,6
19,10,6
7,18,12
3,5,15
13,2,7
10,8,20
4,13,5
11,6,2
4,16,10
4,16,15
7,5,8
8,19,14
18,14,14
10,4,14
3,13,7
10,4,18
7,11,19
2,8,15
10,14,1
15,8,17
5,10,16
4,17,13
13,11,4
7,13,20
3,14,15
10,17,7
17,6,12
13,6,4
1,6,9
14,14,7
5,4,10
10,6,17
10,2,6
15,3,5
9,19,8
14,18,10
18,15,5
9,3,5
13,5,6
7,7,19
7,4,7
4,11,8
5,9,11
8,15,19
14,15,6
15,8,6
15,12,17
10,8,3
16,7,5
1,9,11
8,9,1
16,11,19
2,14,11
10,9,19
9,8,2
14,18,7
15,16,16
10,4,6
2,12,5
8,17,14
10,17,15
13,14,16
11,8,20
6,11,17
13,13,3
11,16,2
9,6,6
9,4,16
9,9,2
3,10,14
5,7,17
7,7,3
1,14,11
5,14,4
7,2,8
15,3,15
16,8,8
5,7,16
6,3,10
14,6,16
18,12,5
11,19,10
12,13,18
14,12,3
5,11,15
20,8,9
12,18,5
6,14,8
10,18,7
10,11,0
14,2,7
6,16,14
4,6,15
6,16,13
4,6,10
4,13,14
13,12,3
10,2,8
10,4,8
5,11,17
6,5,11
5,5,15
13,6,18
9,4,14
14,5,14
16,3,8
4,12,17
15,3,12
6,9,17
10,4,3
16,7,6
8,4,13
10,5,2
5,13,15
12,12,1
17,7,14
14,2,14
17,15,12
7,18,14
19,9,10
17,9,15
9,1,10
2,7,13
9,3,8
15,15,16
2,13,6
3,5,6
15,12,4
13,18,9
19,7,8
14,9,16
15,18,8
17,3,8
16,6,5
6,16,16
13,11,20
10,12,3
6,9,2
15,17,7
14,3,7
2,12,9
19,7,13
4,10,14
5,2,13
6,6,5
17,5,15
9,16,5
10,1,11
16,19,9
12,4,17
7,17,15
18,8,9
20,11,9
17,16,15
10,15,16
13,7,14
6,4,10
10,14,18
15,7,3
17,7,12
11,4,5
19,13,14
10,7,19
14,12,4
16,5,8
14,14,19
11,17,16
18,8,6
8,4,11
13,6,16
9,1,13
6,6,16
1,8,10
12,6,19
8,20,9
16,18,9
4,6,11
12,7,4
8,2,14
8,20,8
10,19,9
8,2,11
3,6,13
12,14,19
17,14,14
15,4,7
13,15,5
12,12,3
8,2,12
2,9,14
6,16,4
4,18,12
5,4,14
14,15,15
5,11,3
15,8,18
5,14,17
5,6,9
13,16,14
10,6,4
12,20,11
5,19,13
12,17,3
4,14,15
3,15,8
13,11,19
15,4,10
16,13,6
11,7,2
5,5,13
3,8,15
20,11,11
3,10,13
15,18,10
5,18,8
2,15,13
6,3,6
11,11,20
8,16,14
14,14,6
5,16,16
10,15,17
3,10,6
19,8,11
2,7,8
7,10,3
16,12,5
12,14,4
9,6,19
3,8,10
3,15,6
11,4,18
10,19,12
9,4,18
2,10,11
18,17,11
11,20,14
15,10,19
17,14,16
7,14,3
16,11,17
2,7,12
14,18,14
4,5,9
3,5,13
7,7,5
8,6,18
7,17,5
8,13,2
7,10,17
14,2,8
18,7,8
13,4,17
13,1,12
10,16,16
19,7,12
20,15,9
13,5,13
5,12,4
5,7,14
15,13,17
17,15,9
16,9,5
9,3,7
17,5,11
17,14,10
16,16,11
13,18,10
2,7,11
2,10,15
15,15,17
3,6,11
9,11,2
14,2,10
2,6,12
8,17,11
16,5,14
15,3,13
20,13,10
12,2,13
13,18,7
10,7,18
12,19,12
17,12,10
9,2,10
8,16,15
6,6,6
11,12,4
7,16,14
3,17,11
8,15,3
2,12,12
14,14,16
8,8,18
8,15,4
17,16,14
8,17,7
12,8,3
7,13,18
14,16,4
5,13,8
5,7,18
1,7,8
4,12,4
14,11,2
15,4,17
3,7,16
10,14,6
5,13,17
14,3,8
16,11,3
19,11,11
3,7,9
8,5,15
7,15,17
6,16,6
18,12,10
13,2,6
20,8,12
13,16,5
13,4,16
11,19,7
2,9,13
10,4,13
12,20,8
11,13,17
5,6,15
14,7,18
7,3,5
16,13,14
13,5,5
13,19,10
13,17,8
9,12,0
8,5,3
3,8,13
5,17,11
6,5,4
13,13,6
5,4,7
12,17,5
14,13,19
14,14,4
11,14,19
13,15,4
10,9,18
3,8,7
13,13,18
6,10,6
7,20,10
8,14,17
17,5,14
13,18,13
10,10,20
9,16,18
14,2,9
13,6,3
15,13,16
19,9,12
18,11,6
15,18,13
16,7,13
4,15,7
14,4,7
6,5,6
8,15,17
16,16,12
15,5,13
17,7,7
15,7,15
3,15,13
7,5,14
16,14,3
19,8,8
18,15,10
9,14,1
17,17,9
4,15,10
18,11,11
12,18,15
18,12,16
3,14,12
18,16,12
6,17,14
11,5,6
12,16,7
5,3,9
3,11,11
18,10,15
5,10,4
9,6,3
13,2,11
8,19,7
14,2,13
15,2,11
9,11,19
20,11,13
8,3,8
1,10,8
7,14,14
12,11,19
9,4,13
5,4,12
9,13,19
9,11,18
10,3,16
11,20,13
18,5,11
2,9,8
7,17,7
15,3,9
7,19,14
15,4,8
9,1,9
6,4,12
6,6,19
9,1,8
16,18,14
2,8,14
14,4,4
2,12,16
8,10,3
2,16,10
14,4,9
16,11,11
7,15,5
13,17,9
16,10,15
9,16,17
9,15,4
3,14,7
17,3,9
5,14,7
3,5,7
3,11,14
19,12,10
9,20,11
1,7,9
10,3,14
13,14,19
15,5,12
8,17,9
18,14,11
14,4,6
16,19,8
9,4,4
13,4,14
9,7,19
5,3,12
12,10,1
2,11,10
17,8,12
17,6,11
3,15,12
14,18,5
2,9,9
20,9,11
8,4,9
9,8,18
19,10,10
18,12,9
4,7,12
9,8,3
11,17,9
5,8,5
14,11,14
11,17,13
7,6,17
12,18,10
14,9,2
8,3,15
6,15,3
12,16,2
7,4,5
6,15,14
17,17,12
10,13,3
8,16,7
12,8,4
10,19,8
11,20,7
3,6,8
8,18,13
6,17,7
11,18,9
16,17,16
8,5,16
7,4,15
3,16,15
15,8,9
18,6,8
14,9,18
7,6,15
4,13,3
16,13,5
9,10,3
15,18,9
2,13,9
14,17,8
10,17,16
9,5,15
10,12,2
10,7,4
3,12,9
10,10,2
12,4,15
9,15,3
12,17,7
3,10,12
7,17,10
10,13,4
11,9,20
19,12,9
16,15,8
19,11,15
13,17,17
14,8,4
7,8,2
17,5,7
17,8,8
10,12,1
2,13,10
6,11,18
7,14,19
19,14,9
5,6,13
6,14,13
10,17,17
12,7,16
12,6,2
3,9,10
10,2,5
5,5,6
12,18,4
14,17,15
9,14,17
6,19,10
13,19,12
7,17,16
17,17,13
7,12,16
9,3,15
7,17,17
8,19,9
14,20,12
15,17,5
8,18,6
10,0,13
15,16,14
1,12,11
14,16,5
4,12,18
12,20,12
20,14,12
13,7,4
14,6,17
14,5,10
6,3,13
14,7,13
10,18,16
8,1,11
9,19,14
17,13,12
6,14,3
12,7,17
12,5,3
8,6,17
10,2,14
11,7,18
14,11,17
2,6,6
17,14,7
3,12,11
2,10,10
15,18,7
13,17,13
17,8,13
3,15,9
7,5,11
17,7,4
15,6,13
16,6,17
3,12,8
11,7,17
9,19,9
1,13,10
10,6,6
18,10,14
7,13,17
17,12,15
7,2,6
12,1,8
9,3,9
4,15,13
8,17,13
14,9,5
13,11,2
15,6,16
13,8,5
7,3,12
12,18,17
7,4,10
3,5,11
9,17,10
12,9,1
3,10,8
17,7,5
15,4,6
15,7,16
2,6,11
17,18,13
8,14,3
7,5,17
3,12,7
5,7,9
5,5,8
18,5,10
2,6,10
16,14,16
8,2,8
14,11,18
14,8,19
6,16,7
6,9,4
18,8,5
3,8,8
16,13,18
18,16,13
14,3,9
13,2,8
17,6,10
17,16,12
11,17,4
12,9,20
3,4,10
12,3,17
5,8,17
7,10,19
16,5,16
5,9,19
4,13,10
13,8,2
3,9,11
4,7,5
4,14,16
10,17,18
12,15,6
16,14,7
11,18,13
15,9,20
8,8,3
3,11,16
18,15,15
3,7,11
3,13,4
10,11,1
2,10,9
3,4,7
14,8,17
2,6,9
11,3,10
4,15,6
6,7,7
8,9,3
1,10,11
18,3,11
5,19,12
5,18,10
19,14,10
12,18,6
13,15,16
14,12,19
4,5,12
19,12,14
8,3,9
13,19,9
14,11,6
12,0,12
9,17,8
4,8,13
17,9,7
17,15,16
10,6,19
7,16,7
19,15,10
14,16,7
12,7,18
18,11,8
3,14,13
12,16,17
6,9,18
2,8,9
5,15,13
18,10,6
13,9,17
9,9,1
8,17,5
6,5,7
18,12,6
1,9,14
15,10,17
11,11,18
11,17,6
20,12,9
20,12,12
9,17,5
19,6,13
18,10,16
14,3,6
14,19,9
17,8,14
8,7,3
18,9,7
8,10,18
9,15,17
14,7,3
17,5,13
11,10,1
11,4,14
14,4,16
9,14,19
8,18,7
19,6,8
20,13,9
16,7,14
2,16,11
11,20,12
14,17,9
4,15,11
12,9,17
4,6,17
9,16,16
17,7,10
16,3,13
13,6,17
0,12,12
10,12,18
4,16,8
11,15,18
15,11,18
4,9,16
11,8,1
17,7,8
18,6,7
13,5,16
14,6,5
6,9,19
16,10,16
5,11,6
11,2,7
14,11,1
14,12,17
9,19,11
17,15,5
19,9,7
7,17,14
20,12,10
5,3,13
10,3,6
5,16,10
11,17,15
18,8,15
13,13,19
8,11,16
10,8,4
13,18,8
16,9,4
3,9,4
7,3,13
2,10,7
16,10,3
11,18,10
10,7,20
4,7,3
9,17,7
18,13,13
8,16,8
1,13,14
4,17,14
5,7,7
3,18,11
7,9,17
12,12,18
17,9,12
4,4,9
11,10,17
6,13,3
8,4,15
10,1,8
3,8,12
4,18,9
2,10,8
10,19,6
12,4,5
13,8,17
5,8,16
14,4,3
17,13,10
12,3,10
15,14,6
16,2,13
1,10,9
7,2,11
15,4,15
9,16,3
13,10,4
8,14,2
7,2,12
7,17,9
6,16,8
19,7,15
14,10,20
8,14,4
15,5,15
16,9,2
7,17,12
1,10,13
14,4,12
9,14,20
11,11,3
9,2,12
8,16,4
3,12,10
19,9,14
16,8,3
7,12,1
4,7,4
7,4,17
13,16,10
14,19,11
12,19,15
6,8,15
2,9,7
2,14,9
2,14,12
4,7,15
14,8,2
17,17,10
14,4,10
13,9,19
10,21,9
15,16,5
14,10,17
11,12,1
4,5,16
9,2,8
19,14,11
13,5,4
5,4,11
6,7,4
12,11,3
4,8,10
16,7,8
8,1,7
3,10,5
13,16,15
15,10,3
9,14,16
19,14,13
5,6,4
19,7,7
4,12,11
5,12,3
17,15,6
7,13,3
7,3,8
6,10,7
17,6,6
11,5,19
7,5,16
16,17,10
10,1,12
12,15,17
15,11,17
10,18,8
12,5,4
3,12,16
10,4,17
4,10,7
6,17,5
8,4,4
12,14,1
4,6,12
4,6,14
8,6,19
5,14,5
3,13,15
18,12,17
3,14,5
15,14,17
3,12,15
11,19,14
2,5,8
14,5,15
14,6,4
10,18,17
7,7,16
11,3,15
4,4,16
7,5,4
16,5,11
10,11,18
8,19,11
4,9,17
17,4,10
11,6,20
8,7,2
13,4,7
18,14,16
5,6,8
12,17,17
13,4,9
9,6,2
17,8,7
9,8,1
2,14,6
11,1,9
11,18,7
11,4,16
5,16,14
6,5,12
3,12,4
10,15,3
11,3,8
5,12,19
19,10,13
12,15,4
9,2,17
18,9,5
17,10,16
15,11,1
15,16,8
8,13,3
3,16,10
10,9,4
5,15,7
9,16,15
7,6,19
13,10,1
10,14,4
18,16,14
11,3,12
16,14,12
5,7,11
1,8,11
13,13,2
17,9,9
8,7,17
12,5,16
8,2,10
11,11,2
16,15,3
5,5,11
4,7,9
8,18,5
15,19,13
3,14,9
1,11,10
19,12,5
11,12,20
3,9,15
4,9,14
3,11,7
19,10,12
8,7,5
6,6,7
10,6,2
6,5,18
12,16,6
18,8,12
9,2,16
18,7,10
4,10,17
13,14,5
8,12,17
3,6,7
17,8,17
2,15,10
5,7,10
9,4,7
8,13,1
16,10,18
17,8,11
6,5,10
8,7,18
12,15,19
15,4,9
15,5,5
12,11,20
16,16,7
19,15,7
16,5,6
18,8,8
6,17,15
14,10,2
12,6,5
13,1,11
9,17,9
7,12,4
13,10,2
14,5,4
10,17,5
13,14,1
10,13,1
8,16,6
19,11,16
13,11,5
14,8,16
11,16,16
6,18,12
16,16,13
5,12,13
12,19,9
12,15,18
12,5,13
8,17,16
12,15,3
17,6,5
15,9,16
10,11,19
13,7,3
12,12,4
7,16,17
3,14,6
10,19,10
18,13,10
10,10,18
12,6,1
11,15,5
12,15,2
19,14,7
9,4,10
9,11,20
7,16,8
17,7,9
11,13,3
5,5,16
16,6,9
16,13,16
7,19,7
8,2,15
17,3,11
7,18,13
17,12,3
4,15,12
6,10,16
17,10,17
17,11,13
16,13,15
20,9,8
17,14,5
15,2,9
15,15,14
5,17,14
19,7,10
9,8,19
8,5,8
9,12,18
5,17,13
8,11,18
2,11,16
1,7,11
12,5,7
4,14,5
4,13,7
17,14,9
5,15,5
11,17,12
3,11,10
17,9,18
13,16,6
20,11,14
15,18,14
4,8,17
14,19,7
13,6,19
3,15,7
4,7,6
10,14,3
18,15,7
9,17,13
2,11,9
7,16,15
13,9,16
6,13,19
12,13,3
18,9,13
5,6,7
18,11,15
6,14,4
9,18,7
13,12,2
17,10,14
16,8,17
4,9,5
4,17,7
7,15,6
15,11,5
2,12,6
14,15,14
10,10,0
2,6,7
4,13,13
6,10,17
10,6,20
10,18,9
5,11,16
4,10,5
10,14,20
4,12,10
11,9,16
3,13,10
17,6,13
6,1,9
7,3,9
13,20,11
11,13,19
15,12,3
17,14,6
4,9,6
16,15,11
19,15,8
18,7,5
11,1,11
5,10,7
3,9,5
9,5,3
2,8,7
16,5,5
16,15,6
5,9,5
16,3,16
7,15,2
3,7,7
14,10,19
2,12,11
5,6,3
9,2,6
7,12,3
4,4,10
4,13,8
5,5,4
2,12,15
12,8,20
6,7,17
2,9,11
13,15,17
3,12,12
12,15,5
12,2,14
16,17,13
18,12,4
18,5,14
3,7,14
10,4,7
11,1,14
5,13,6
7,13,2
5,9,16
9,7,17
10,3,15
14,17,11
7,10,4
8,8,1
13,9,18
15,2,12
13,12,19
4,12,13
7,7,17
19,7,9
16,6,4
14,6,7
1,8,13
14,6,2
18,8,16
12,6,6
5,15,4
11,13,5
10,8,19
11,12,3
4,10,16
4,11,6
11,3,17
9,10,19
12,6,4
10,19,15
3,9,16
11,17,3
2,7,10
11,2,11
6,18,8
15,12,18
11,7,4
19,8,7
18,11,16
1,12,10
6,19,9
12,14,2
14,10,4
15,12,2
18,9,14
14,8,6
8,16,18
1,8,12
20,13,13
17,7,17
9,11,3
9,13,4
1,12,12
9,9,4
14,12,16
8,4,8
20,13,12
3,9,9
13,19,14
7,16,6
18,10,11
4,7,8
3,7,5
16,19,11
15,8,2
18,11,12
2,15,9
13,3,5
6,9,5
4,9,13
5,16,12
8,8,19
5,17,6
15,6,19
3,5,8
15,12,19
9,18,4
15,11,15
9,16,4
19,5,10
13,4,6
18,5,12
1,9,13
15,11,19
8,17,12
3,9,8
4,11,18
13,13,17
4,16,14
4,17,10
10,3,9
3,10,18
17,15,10
14,19,8
11,20,8
13,8,19
14,17,12
8,15,5
11,10,3
16,12,15
15,17,9
12,5,18
13,3,6
8,2,9
14,8,18
16,19,12
4,7,13
17,12,17
18,15,6
7,15,18
4,3,12
15,7,6
17,4,8
8,17,10
9,8,4
9,4,12
16,13,3
6,15,18
3,6,12
12,10,2
11,11,1
14,13,3
12,13,1
17,16,6
14,16,11
18,7,16
16,4,15
14,16,18
14,15,13
15,7,8
20,11,5
6,2,12
9,15,19
10,8,18
18,9,3
4,17,8
4,8,6
12,4,7
7,11,4
13,1,10
3,13,16
21,12,11
8,20,14
4,4,8
12,1,13
11,4,4
11,6,16
15,6,8
19,13,7
5,8,7
10,20,9
19,6,12
9,18,13
10,15,2
9,2,13
9,17,16
14,17,6
6,12,2
10,5,4
7,5,13
7,9,2
16,14,18
9,7,16
12,10,17
11,19,15
16,5,15
17,13,7
16,3,12
4,14,12
14,3,13
12,10,5
19,13,16
9,6,18
14,18,12
5,6,12
9,9,3
5,12,15
0,10,11
4,6,6
10,8,1
14,5,12
2,6,14
3,10,15
14,15,4
5,7,4
12,16,8
3,17,8
13,2,15
19,14,14
15,11,3
5,15,11
16,7,9
20,13,14
16,15,15
5,13,3
8,5,2
7,8,16
12,5,2
10,18,13
13,14,14
6,8,3
15,5,17
14,8,15
18,6,5
10,19,16
15,17,4
6,5,14
4,4,14
14,10,16
2,6,13
9,7,6
15,9,15
15,10,16
1,13,8
7,13,5
9,20,7
5,14,14
7,3,10
12,2,16
10,6,1
6,17,16
6,12,20
13,13,4
6,7,10
7,13,4
10,2,13
16,3,15
5,16,11
8,8,4
17,15,8
14,3,11
17,11,6
6,15,6
5,5,17
7,14,5
14,7,2
12,16,14
8,18,15
18,9,15
11,9,1
19,12,11
6,6,3
9,4,8
14,20,10
20,10,9
7,15,15
15,11,16
11,18,17
10,13,20
9,16,2
18,17,10
10,15,11
20,10,7
16,19,6
14,15,5
16,17,15
9,7,2
18,11,17
16,16,16
10,17,6
9,14,18
11,18,8
9,17,17
13,16,7
15,5,14
2,14,7
7,18,16
9,16,6
15,16,11
9,6,16
9,6,4
12,1,9
19,16,12
13,4,4
12,17,4
6,5,17
14,13,4
3,8,11
17,4,13
15,11,14
3,11,4
8,8,17
11,8,4
7,13,1
6,8,16
17,18,12
11,17,14
18,7,6
11,6,19
18,13,11
4,16,6
6,12,16
5,15,12
15,2,15
6,18,13
5,18,13
10,16,15
5,13,5
7,18,6
18,17,8
10,12,20
6,11,19
5,11,11
9,2,5
7,4,13
16,16,15
15,16,4
15,3,14
7,1,10
1,14,12
10,11,3
14,19,14
11,16,7
20,7,11
14,19,12
6,6,18
9,13,17
9,5,19
10,20,10
8,16,3
5,8,6
5,8,11
13,5,3
9,4,17
20,9,13
4,8,9
6,15,8
18,11,5
5,9,14
16,17,12
12,20,9
11,15,16
11,20,10
5,2,11
4,12,8
13,1,8
6,3,15
11,14,4
5,12,5
16,8,18
19,14,12
4,8,8
12,16,9
3,14,11
15,12,6
19,8,14
15,16,13
8,20,13
5,9,15
14,16,13
7,3,16
16,6,12
10,6,5
2,16,9
7,14,15
19,6,9
12,8,1
13,17,11
13,5,14
7,16,12
16,8,12
14,16,6
15,16,12
5,15,6
17,18,10
4,10,8
17,8,4
8,15,16
5,10,6
15,6,18
12,14,17
19,16,10
16,6,13
16,19,10
1,12,13
3,11,12
8,5,13
7,17,11
12,12,17
8,18,8
16,18,7
17,14,4
3,11,15
16,13,8
3,12,6
2,9,6
20,7,10
15,7,18
9,20,12
6,18,15
17,12,9
6,19,11
9,15,5
8,4,3
7,5,19
17,12,14
5,11,4
9,12,2
4,14,8
13,15,6
16,19,13
15,19,7
18,13,17
4,8,11
10,9,1
14,1,11
17,18,11
16,11,4
15,7,5
10,14,2
12,18,13
17,11,12
8,17,3
16,8,15
6,2,14
5,10,14
5,10,13
18,11,4
5,15,16
10,19,7
12,19,7
3,15,16
7,8,19
16,7,12
16,2,12
9,20,10
6,14,15
8,11,4
5,3,8
10,19,5
17,18,9
11,17,8
4,5,15
5,11,7
18,15,13
7,18,8
4,13,9
13,7,16
13,20,14
10,1,10
6,6,8
6,3,14
11,3,14
8,9,20
17,6,17
20,9,9
17,9,8
14,11,5
5,10,5
8,20,10`
