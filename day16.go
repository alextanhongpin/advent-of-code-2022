package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	fmt.Println(part1(input1)) // 1651
	fmt.Println(part1(input2)) // 1850
	fmt.Println(part2(input1)) // 1707
	// This works now, but it is very slow.
	fmt.Println(part2(input2)) // 2306
}

type state struct {
	valve   string
	minutes int
	rate    int
	visited map[string]bool
}

type state2 struct {
	valves  [2]string
	minutes [2]int
	rates   [2]int
	visited map[string]bool
}

func (s state2) rate() int {
	return s.rates[0] + s.rates[1]
}

func checkMove(h state, valves map[string]valve, timeTakenToVisitValveByOrigin map[string]map[string]int) int {
	var res int

	if h.minutes == 0 {
		return h.rate
	}

	visitable := func(h state) map[string]int {
		res := make(map[string]int)
		for v, t := range timeTakenToVisitValveByOrigin[h.valve] {
			// Don't revisit the same valve.
			if h.visited[v] {
				continue
			}

			// Not enough time to visit (1 minute) and open (another 1 minute)
			if h.minutes-t < 1 {
				continue
			}

			res[v] = t
		}

		return res
	}

	var move func(h state) int
	move = func(h state) int {
		for v, t := range visitable(h) {
			visited := copyMap(h.visited)
			visited[v] = true

			// The t+1 refers to the additional time taken to open the valve after visiting.
			mins := h.minutes - (t + 1)
			rate := mins*valves[v].rate + h.rate
			next := state{
				valve:   v,
				rate:    rate,
				minutes: mins,
				visited: visited,
			}

			res = max(res, max(rate, move(next)))
		}
		return max(res, h.rate)
	}

	return move(h)
}

func checkMove2(h state2, valves map[string]valve, timeTakenToVisitValveByOrigin map[string]map[string]int) int {
	cache := make(map[int]int)

	cacheKey := func(h state2) int {
		return len(h.visited)
	}

	visitable := func(h state2, at int) map[string]int {
		res := make(map[string]int)
		for v, t := range timeTakenToVisitValveByOrigin[h.valves[at]] {
			// Don't revisit the same valve.
			if h.visited[v] {
				continue
			}

			// Not enough time to visit (1 minute) and open (another 1 minute)
			if h.minutes[at]-t < 1 {
				continue
			}
			res[v] = t
		}

		return res
	}

	var move func(h state2) int
	move = func(h state2) int {
		var res int

		cache[cacheKey(h)] = max(cache[cacheKey(h)], h.rate())

		visit0 := visitable(h, 0)
		visit1 := visitable(h, 1)

		if len(visit0) != 0 && len(visit1) != 0 {
			// Visit both.
			for v0, t0 := range visit0 {
				visited0 := copyMap(h.visited)
				visited0[v0] = true
				mins0 := h.minutes[0] - (t0 + 1)
				rate0 := mins0*valves[v0].rate + h.rates[0]

				for v1, t1 := range visit1 {
					if visited0[v1] {
						continue
					}

					visited1 := copyMap(visited0)
					visited1[v1] = true
					mins1 := h.minutes[1] - (t1 + 1)
					rate1 := mins1*valves[v1].rate + h.rates[1]

					next := state2{
						valves:  [2]string{v0, v1},
						rates:   [2]int{rate0, rate1},
						minutes: [2]int{mins0, mins1},
						visited: visited1,
					}
					cache[cacheKey(next)] = max(cache[cacheKey(next)], next.rate())
					if n, ok := cache[cacheKey(next)-1]; ok && next.rate() < n {
						continue
					}
					res = max(res, move(next))
				}
			}
		} else if len(visit0) == 0 {
			// Visit 1 only.
			for v1, t1 := range visit1 {
				visited1 := copyMap(h.visited)
				visited1[v1] = true
				mins1 := h.minutes[1] - (t1 + 1)
				rate1 := mins1*valves[v1].rate + h.rates[1]

				next := state2{
					valves:  [2]string{h.valves[0], v1},
					rates:   [2]int{h.rates[0], rate1},
					minutes: [2]int{h.minutes[0], mins1},
					visited: visited1,
				}
				cache[cacheKey(next)] = max(cache[cacheKey(next)], next.rate())
				if n, ok := cache[cacheKey(next)-1]; ok && next.rate() < n {
					continue
				}
				res = max(res, move(next))
			}
		} else if len(visit1) == 0 {
			// Visit 0 only.
			for v0, t0 := range visit0 {
				visited0 := copyMap(h.visited)
				visited0[v0] = true
				mins0 := h.minutes[0] - (t0 + 1)
				rate0 := mins0*valves[v0].rate + h.rates[0]

				next := state2{
					valves:  [2]string{v0, h.valves[1]},
					rates:   [2]int{rate0, h.rates[1]},
					minutes: [2]int{mins0, h.minutes[1]},
					visited: visited0,
				}
				cache[cacheKey(next)] = max(cache[cacheKey(next)], next.rate())
				if n, ok := cache[cacheKey(next)-1]; ok && next.rate() < n {
					continue
				}
				res = max(res, move(next))
			}
		} else {
			return h.rate()
		}

		cache[cacheKey(h)] = max(res, h.rate())

		return max(res, h.rate())
	}

	return move(h)
}

func part1(input string) int {
	valves := parse(input)
	timeTakenToVisitValveByOrigin := computeTimeTakenToVisitValveByOrigin(valves)

	h := state{
		valve:   "AA",
		minutes: 30,
		rate:    0,
		visited: make(map[string]bool),
	}

	return checkMove(h, valves, timeTakenToVisitValveByOrigin)
}

func part2(input string) int {
	valves := parse(input)
	timeTakenToVisitValveByOrigin := computeTimeTakenToVisitValveByOrigin(valves)

	h := state2{
		valves:  [2]string{"AA", "AA"},
		minutes: [2]int{26, 26},
		rates:   [2]int{0, 0},
		visited: make(map[string]bool),
	}

	return checkMove2(h, valves, timeTakenToVisitValveByOrigin)
}

func copyMap[K comparable, V any](m map[K]V) map[K]V {
	c := make(map[K]V)
	for k, v := range m {
		c[k] = v
	}
	return c
}

type valve struct {
	valves []string
	rate   int
}

func computeTimeTakenToVisitValveByOrigin(valves map[string]valve) map[string]map[string]int {
	visit := make(map[string]bool)
	for v := range valves {
		visit[v] = true
	}

	res := make(map[string]map[string]int)
	for v := range visit {
		visited := make(map[string]bool)
		res[v] = make(map[string]int)
		mins := 0
		q := []string{v}
		for len(q) > 0 {
			var toVisit []string
			for _, vv := range q {
				if visited[vv] {
					continue
				}
				visited[vv] = true
				toVisit = append(toVisit, vv)

				// Only visit valves with positive flow rate.
				if valves[vv].rate > 0 {
					res[v][vv] = mins
				}
			}
			mins++
			q = []string{}
			for _, vv := range toVisit {
				q = append(q, valves[vv].valves...)
			}
		}
	}
	return res
}

func parse(input string) map[string]valve {
	rows := strings.Split(input, "\n")
	g := make(map[string]valve)
	for _, row := range rows {
		parts := strings.Fields(row)

		from := parts[1]

		toRaw := parts[9:]
		to := make([]string, len(toRaw))
		for j, t := range toRaw {
			t = strings.TrimSuffix(t, ",")
			to[j] = t
		}

		rate := strings.TrimSuffix(strings.TrimPrefix(parts[4], "rate="), ";")

		g[from] = valve{
			valves: to,
			rate:   toInt(rate),
		}
	}

	return g
}

func toInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// An Item is something we manage in a priority queue.
type Item struct {
	priority int // The priority of the item in the queue.

	yvalve, evalve     string
	yminutes, eminutes int
	yrate, erate       int
	visited            map[string]bool
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

// We want Pop to give us the highest, not lowest, priority so we use greater than here.
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].priority > pq[j].priority }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x any) {
	item := x.(*Item)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var input1 = `Valve AA has flow rate=0; tunnels lead to valves DD, II, BB
Valve BB has flow rate=13; tunnels lead to valves CC, AA
Valve CC has flow rate=2; tunnels lead to valves DD, BB
Valve DD has flow rate=20; tunnels lead to valves CC, AA, EE
Valve EE has flow rate=3; tunnels lead to valves FF, DD
Valve FF has flow rate=0; tunnels lead to valves EE, GG
Valve GG has flow rate=0; tunnels lead to valves FF, HH
Valve HH has flow rate=22; tunnel leads to valve GG
Valve II has flow rate=0; tunnels lead to valves AA, JJ
Valve JJ has flow rate=21; tunnel leads to valve II`

var input2 = `Valve DR has flow rate=22; tunnels lead to valves DC, YA
Valve IO has flow rate=14; tunnels lead to valves GE, CK, HY, XB
Valve XY has flow rate=0; tunnels lead to valves IP, AR
Valve UQ has flow rate=0; tunnels lead to valves XU, PD
Valve FO has flow rate=0; tunnels lead to valves DL, NC
Valve PU has flow rate=0; tunnels lead to valves ZJ, AN
Valve MK has flow rate=0; tunnels lead to valves ZS, SB
Valve HN has flow rate=0; tunnels lead to valves AA, DV
Valve XF has flow rate=0; tunnels lead to valves XB, AA
Valve OD has flow rate=13; tunnels lead to valves ZS, AF, SY, QQ, AR
Valve GE has flow rate=0; tunnels lead to valves KR, IO
Valve UF has flow rate=18; tunnels lead to valves QQ, AN, YE, GY
Valve WK has flow rate=19; tunnel leads to valve PQ
Valve PQ has flow rate=0; tunnels lead to valves WK, CW
Valve XU has flow rate=0; tunnels lead to valves DV, UQ
Valve SH has flow rate=0; tunnels lead to valves IP, AA
Valve SY has flow rate=0; tunnels lead to valves ZJ, OD
Valve OU has flow rate=0; tunnels lead to valves CK, DL
Valve IP has flow rate=8; tunnels lead to valves CY, ML, YI, XY, SH
Valve XZ has flow rate=0; tunnels lead to valves AM, PD
Valve ZU has flow rate=0; tunnels lead to valves CW, SB
Valve DC has flow rate=0; tunnels lead to valves CF, DR
Valve QY has flow rate=0; tunnels lead to valves CW, MQ
Valve XB has flow rate=0; tunnels lead to valves IO, XF
Valve AF has flow rate=0; tunnels lead to valves PD, OD
Valve GY has flow rate=0; tunnels lead to valves UF, ZC
Valve ZC has flow rate=0; tunnels lead to valves GY, CW
Valve ZJ has flow rate=25; tunnels lead to valves SY, PU
Valve NC has flow rate=6; tunnels lead to valves HY, ML, NJ, AT, FO
Valve DS has flow rate=0; tunnels lead to valves AT, DV
Valve DV has flow rate=7; tunnels lead to valves FD, KR, HN, DS, XU
Valve HY has flow rate=0; tunnels lead to valves NC, IO
Valve WF has flow rate=0; tunnels lead to valves NJ, AA
Valve CK has flow rate=0; tunnels lead to valves IO, OU
Valve YE has flow rate=0; tunnels lead to valves CY, UF
Valve LA has flow rate=0; tunnels lead to valves DL, ZM
Valve QQ has flow rate=0; tunnels lead to valves OD, UF
Valve AM has flow rate=0; tunnels lead to valves XZ, SB
Valve AN has flow rate=0; tunnels lead to valves UF, PU
Valve CL has flow rate=16; tunnels lead to valves YA, LD
Valve CF has flow rate=12; tunnel leads to valve DC
Valve FD has flow rate=0; tunnels lead to valves DV, DL
Valve QU has flow rate=0; tunnels lead to valves LD, PD
Valve AT has flow rate=0; tunnels lead to valves DS, NC
Valve SB has flow rate=24; tunnels lead to valves MK, AM, ZU
Valve YI has flow rate=0; tunnels lead to valves DL, IP
Valve ZM has flow rate=0; tunnels lead to valves AA, LA
Valve LD has flow rate=0; tunnels lead to valves CL, QU
Valve AR has flow rate=0; tunnels lead to valves OD, XY
Valve DL has flow rate=5; tunnels lead to valves FO, LA, YI, OU, FD
Valve MQ has flow rate=0; tunnels lead to valves QY, PD
Valve PD has flow rate=9; tunnels lead to valves MQ, QU, XZ, AF, UQ
Valve KR has flow rate=0; tunnels lead to valves GE, DV
Valve CY has flow rate=0; tunnels lead to valves YE, IP
Valve AA has flow rate=0; tunnels lead to valves SH, XF, ZM, HN, WF
Valve NJ has flow rate=0; tunnels lead to valves NC, WF
Valve YA has flow rate=0; tunnels lead to valves CL, DR
Valve ML has flow rate=0; tunnels lead to valves NC, IP
Valve CW has flow rate=15; tunnels lead to valves QY, PQ, ZC, ZU
Valve ZS has flow rate=0; tunnels lead to valves MK, OD`
