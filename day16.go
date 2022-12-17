package main

import (
	"container/heap"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	fmt.Println(part1(input1))
	fmt.Println(part1(input2))
	// Doesn't work
	fmt.Println(part2(input1))
	fmt.Println(part2(input2)) // 2306
}

type state struct {
	valve   string
	minutes int
	rate    int
	visited map[string]bool
}

func checkMove(h state, valves map[string]valve, timeTakenToVisitValveByOrigin map[string]map[string]int) []state {
	var res []state

	for v, t := range timeTakenToVisitValveByOrigin[h.valve] {
		rate := valves[v].rate
		// Uncomment to debug.
		//fmt.Printf("left %d, visiting %v takes %d (%t), score=%d\n", h.minutes, v, t, h.visited[v], (h.minutes-(t+1))*rate)
		// Don't revisit the same valve.
		if h.visited[v] {
			continue
		}

		// Not enough time to visit (1 minute) and open (another 1 minute)
		if h.minutes-t < 1 {
			continue
		}
		visited := make(map[string]bool)
		for k, v := range h.visited {
			visited[k] = v
		}
		visited[v] = true

		// The t+1 refers to the additional time taken to open the valve after visiting.
		res = append(res, state{
			valve:   v,
			rate:    (h.minutes-(t+1))*rate + h.rate,
			minutes: h.minutes - (t + 1),
			visited: visited,
		})
	}

	return res
}

func part1(input string) int {
	valves := parse(input)
	timeTakenToVisitValveByOrigin := computeTimeTakenToVisitValveByOrigin(valves)

	var valvesCount int
	for _, d := range valves {
		if d.rate > 0 {
			valvesCount++
		}
	}

	q := []state{state{
		valve:   "AA",
		minutes: 30,
		rate:    0,
		visited: make(map[string]bool),
	}}

	var max int
	for len(q) > 0 {
		var h state
		h, q = q[0], q[1:]
		// Terminates after 30 minutes has elapsed, or when all the valves
		// (non-zero flow rate) are opened.
		if h.minutes == 0 || len(h.visited) == valvesCount {
			if h.rate > max {
				max = h.rate
			}
		}

		moves := checkMove(h, valves, timeTakenToVisitValveByOrigin)
		q = append(q, moves...)
	}

	return max
}

func part2(input string) int {
	valves := parse(input)
	timeTakenToVisitValveByOrigin := computeTimeTakenToVisitValveByOrigin(valves)

	var valvesCount int
	for _, d := range valves {
		if d.rate > 0 {
			valvesCount++
		}
	}

	pq := &PriorityQueue{&Item{
		// y for yours, e for elephant.
		yvalve:   "AA",
		evalve:   "AA",
		yminutes: 26,
		eminutes: 26,
		yrate:    0,
		erate:    0,
		visited:  make(map[string]bool),
	}}

	checkHasRemainingTime := func(h *Item) bool {
		for v, t := range timeTakenToVisitValveByOrigin[h.yvalve] {
			if h.visited[v] {
				continue
			}
			if h.yminutes-t < 1 {
				continue
			}
			return true
		}

		for vv, tt := range timeTakenToVisitValveByOrigin[h.evalve] {
			if h.visited[vv] {
				continue
			}
			if h.eminutes-tt < 1 {
				continue
			}
			return true
		}

		return false
	}

	var counter int
	var max int
	for pq.Len() > 0 {
		//for len(q) > 0 {
		counter++
		//if counter%10_000 == 0 {
		//fmt.Println(len(q))
		fmt.Println(pq.Len())
		//}

		h := heap.Pop(pq).(*Item)
		//var h state
		//h, q = q[0], q[1:]

		if !checkHasRemainingTime(h) || h.yminutes == 0 || h.eminutes == 0 || len(h.visited) == valvesCount {
			rate := h.yrate + h.erate
			if rate > max {
				max = rate
			}
			fmt.Println("max", max, pq.Len())
			//return max
			continue
			//return max
		}

		ymoves := checkMove(state{
			valve:   h.yvalve,
			rate:    h.yrate,
			minutes: h.yminutes,
			visited: copyMap(h.visited),
		}, valves, timeTakenToVisitValveByOrigin)

		emoves := checkMove(state{
			valve:   h.evalve,
			rate:    h.erate,
			minutes: h.eminutes,
			visited: copyMap(h.visited),
		}, valves, timeTakenToVisitValveByOrigin)

		if len(ymoves) == 0 {
			for _, emove := range emoves {
				heap.Push(pq, &Item{
					priority: emove.rate,
					yvalve:   h.yvalve,
					yrate:    h.yrate,
					yminutes: h.yminutes,
					evalve:   emove.valve,
					erate:    emove.rate,
					eminutes: emove.minutes,
					visited:  copyMap(emove.visited),
				})
			}
		} else {
			for _, ymove := range ymoves {
				if len(emoves) == 0 {
					heap.Push(pq, &Item{
						priority: ymove.rate,
						yvalve:   ymove.valve,
						yrate:    ymove.rate,
						yminutes: ymove.minutes,
						evalve:   h.evalve,
						erate:    h.erate,
						eminutes: h.eminutes,
						visited:  copyMap(ymove.visited),
					})
				} else {
					for _, emove := range emoves {
						if emove.visited[ymove.valve] || ymove.visited[emove.valve] {
							continue
						}
						visited := copyMap(ymove.visited)
						for k, v := range emove.visited {
							visited[k] = v
						}

						heap.Push(pq, &Item{
							priority: ymove.rate + emove.rate,
							yvalve:   ymove.valve,
							yrate:    ymove.rate,
							yminutes: ymove.minutes,
							evalve:   emove.valve,
							erate:    emove.rate,
							eminutes: emove.minutes,
							visited:  visited,
						})
					}
				}
			}
		}
	}

	return max
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
