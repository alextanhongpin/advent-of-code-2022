package main

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func main() {
	fmt.Println(part1(input1))
	fmt.Println(part1(input2)) // 1306
	fmt.Println(part2(input1))
	fmt.Println(part2(input2)) // 37604
}

type state struct {
	ores   map[string]int
	robots map[string]int
}

func (s state) produce() state {
	ss := s.clone()
	for r, c := range ss.robots {
		ss.ores[r] += c
	}
	return ss
}

func (s state) clone() state {
	return state{
		ores:   copyMap(s.ores),
		robots: copyMap(s.robots),
	}
}

func (s state) String() string {
	var ores, robots []string
	for k, v := range s.ores {
		ores = append(ores, fmt.Sprintf("%v:%v", k, v))
	}
	for k, v := range s.robots {
		robots = append(robots, fmt.Sprintf("%v:%v", k, v))
	}
	sort.Strings(ores)
	sort.Strings(robots)
	return fmt.Sprintf("%v:%v", ores, robots)
}

func newState() state {
	ores := make(map[string]int)
	robots := make(map[string]int)
	robots["ore"] = 1
	return state{
		ores:   ores,
		robots: robots,
	}
}

func part1(input string) int {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var total int
	blueprints := parse(input)
	wg.Add(len(blueprints))
	rateLimit := make(chan bool, 10)
	for i, recipes := range blueprints {
		rateLimit <- true
		go func(i int, recipes []recipe) {
			defer func() {
				<-rateLimit
			}()

			defer wg.Done()

			maxOr := maxOre(recipes)
			//fmt.Println(i, recipes)
			// Simulate 24 minutes.
			q := []state{newState()}

			c := make(map[string]bool)
			var bestScore int = math.MinInt
			for t := 0; t < 24; t++ {
				j := len(q)
				fmt.Println(i, t, j, bestScore, q[len(q)-1])
				for _, s := range q[:j] {
					if s.ores["geode"] > bestScore {
						bestScore = s.ores["geode"]
					}
					if s.ores["geode"] < bestScore-1 {
						continue
					}
					ns := next(s, recipes)
					alreadyMax := false
					for _, n := range ns {
						if c[n.String()] {
							continue
						}
						c[n.String()] = true
						if n.ores["ore"] >= maxOr {
							alreadyMax = true
						}
						q = append(q, n)
					}
					if !alreadyMax {
						if c[s.produce().String()] {
							continue
						}
						c[s.produce().String()] = true
						q = append(q, s.produce())
					}
				}
				q = q[j:]
			}
			for _, s := range q {
				if s.ores["geode"] > bestScore {
					bestScore = s.ores["geode"]
				}
			}

			qualityLevel := bestScore * (i + 1)
			fmt.Println("quality level for", i+1, bestScore, qualityLevel)
			mu.Lock()
			total += qualityLevel
			mu.Unlock()

		}(i, recipes)
	}
	wg.Wait()

	return total
}

func part2(input string) int {
	var wg sync.WaitGroup
	var mu sync.Mutex
	total := 1
	blueprints := parse(input)
	if len(blueprints) > 3 {
		blueprints = blueprints[:3]
	}
	wg.Add(len(blueprints))
	rateLimit := make(chan bool, 10)
	for i, recipes := range blueprints {
		rateLimit <- true
		go func(i int, recipes []recipe) {
			defer func() {
				<-rateLimit
			}()

			defer wg.Done()

			maxOr := maxOre(recipes)
			//fmt.Println(i, recipes)
			// Simulate 24 minutes.
			q := []state{newState()}

			c := make(map[string]bool)
			var bestScore int = math.MinInt
			for t := 0; t < 32; t++ {
				j := len(q)
				fmt.Println(i, t, j, bestScore, q[len(q)-1])
				for _, s := range q[:j] {
					if s.ores["geode"] > bestScore {
						bestScore = s.ores["geode"]
					}
					if s.ores["geode"] < bestScore-1 {
						continue
					}
					ns := next(s, recipes)
					alreadyMax := false
					for _, n := range ns {
						if c[n.String()] {
							continue
						}
						c[n.String()] = true
						if n.ores["ore"] >= maxOr {
							alreadyMax = true
						}
						q = append(q, n)
					}
					if !alreadyMax {
						if c[s.produce().String()] {
							continue
						}
						c[s.produce().String()] = true
						q = append(q, s.produce())
					}
				}
				q = q[j:]
			}
			for _, s := range q {
				if s.ores["geode"] > bestScore {
					bestScore = s.ores["geode"]
				}
			}

			fmt.Println("quality level for", i+1, bestScore)
			mu.Lock()
			total *= bestScore
			mu.Unlock()

		}(i, recipes)
	}
	wg.Wait()

	return total
}

func next(s state, recipes []recipe) []state {
	states := []state{}
	for _, r := range recipes {
		ores, ok := r.build(s.ores)
		if ok {
			//fmt.Println("before", s)
			sc := s.clone()
			sc.ores = ores
			sc = sc.produce()
			sc.robots[r.output]++
			//fmt.Println("after", sc)
			states = append(states, sc)
		}
	}

	return states
}

type ingredient struct {
	robot string
	cost  int
}

type recipe struct {
	input  []ingredient
	output string
}

func maxOre(recipes []recipe) int {
	cost := math.MinInt
	for _, r := range recipes {
		for _, i := range r.input {
			if i.robot == "ore" && i.cost > cost {
				cost = i.cost
			}
		}
	}

	return cost
}

func (r recipe) build(m map[string]int) (map[string]int, bool) {
	cm := copyMap(m)
	for _, in := range r.input {
		cm[in.robot] -= in.cost
		if cm[in.robot] < 0 {
			return m, false
		}
	}

	return cm, true
}

func parse(input string) [][]recipe {
	var res [][]recipe
	rows := strings.Split(input, "\n")
	for _, row := range rows {
		parts := strings.Split(row, ":")

		recipes := strings.Split(strings.TrimSpace(parts[1]), ".")
		var out []recipe

		for _, r := range recipes {
			if r == "" {
				continue
			}
			p := strings.Fields(strings.TrimSpace(r))
			switch len(p) {
			case 6:
				rec := recipe{
					output: p[1],
					input:  []ingredient{{p[5], toInt(p[4])}},
				}
				out = append(out, rec)
			case 9:
				rec := recipe{
					output: p[1],
					input:  []ingredient{{p[5], toInt(p[4])}, {p[8], toInt(p[7])}},
				}
				out = append(out, rec)
			default:
				panic(r)
			}
		}
		res = append(res, out)
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

func copyMap[K comparable, V any](m map[K]V) map[K]V {
	res := make(map[K]V)
	for k, v := range m {
		res[k] = v
	}
	return res
}

var input1 = `Blueprint 1: Each ore robot costs 4 ore. Each clay robot costs 2 ore. Each obsidian robot costs 3 ore and 14 clay. Each geode robot costs 2 ore and 7 obsidian.
Blueprint 2: Each ore robot costs 2 ore. Each clay robot costs 3 ore. Each obsidian robot costs 3 ore and 8 clay. Each geode robot costs 3 ore and 12 obsidian.`
var input2 = `Blueprint 1: Each ore robot costs 4 ore. Each clay robot costs 4 ore. Each obsidian robot costs 4 ore and 12 clay. Each geode robot costs 3 ore and 8 obsidian.
Blueprint 2: Each ore robot costs 2 ore. Each clay robot costs 2 ore. Each obsidian robot costs 2 ore and 15 clay. Each geode robot costs 2 ore and 7 obsidian.
Blueprint 3: Each ore robot costs 4 ore. Each clay robot costs 3 ore. Each obsidian robot costs 4 ore and 18 clay. Each geode robot costs 4 ore and 11 obsidian.
Blueprint 4: Each ore robot costs 2 ore. Each clay robot costs 2 ore. Each obsidian robot costs 2 ore and 10 clay. Each geode robot costs 2 ore and 11 obsidian.
Blueprint 5: Each ore robot costs 3 ore. Each clay robot costs 3 ore. Each obsidian robot costs 2 ore and 9 clay. Each geode robot costs 2 ore and 9 obsidian.
Blueprint 6: Each ore robot costs 3 ore. Each clay robot costs 3 ore. Each obsidian robot costs 2 ore and 12 clay. Each geode robot costs 2 ore and 10 obsidian.
Blueprint 7: Each ore robot costs 4 ore. Each clay robot costs 4 ore. Each obsidian robot costs 4 ore and 10 clay. Each geode robot costs 2 ore and 7 obsidian.
Blueprint 8: Each ore robot costs 4 ore. Each clay robot costs 4 ore. Each obsidian robot costs 2 ore and 10 clay. Each geode robot costs 3 ore and 14 obsidian.
Blueprint 9: Each ore robot costs 3 ore. Each clay robot costs 4 ore. Each obsidian robot costs 3 ore and 17 clay. Each geode robot costs 3 ore and 8 obsidian.
Blueprint 10: Each ore robot costs 3 ore. Each clay robot costs 3 ore. Each obsidian robot costs 3 ore and 11 clay. Each geode robot costs 2 ore and 8 obsidian.
Blueprint 11: Each ore robot costs 4 ore. Each clay robot costs 3 ore. Each obsidian robot costs 3 ore and 20 clay. Each geode robot costs 2 ore and 19 obsidian.
Blueprint 12: Each ore robot costs 3 ore. Each clay robot costs 3 ore. Each obsidian robot costs 3 ore and 20 clay. Each geode robot costs 2 ore and 12 obsidian.
Blueprint 13: Each ore robot costs 3 ore. Each clay robot costs 4 ore. Each obsidian robot costs 4 ore and 6 clay. Each geode robot costs 2 ore and 20 obsidian.
Blueprint 14: Each ore robot costs 4 ore. Each clay robot costs 4 ore. Each obsidian robot costs 3 ore and 5 clay. Each geode robot costs 3 ore and 18 obsidian.
Blueprint 15: Each ore robot costs 3 ore. Each clay robot costs 3 ore. Each obsidian robot costs 4 ore and 19 clay. Each geode robot costs 4 ore and 7 obsidian.
Blueprint 16: Each ore robot costs 3 ore. Each clay robot costs 4 ore. Each obsidian robot costs 4 ore and 19 clay. Each geode robot costs 4 ore and 11 obsidian.
Blueprint 17: Each ore robot costs 2 ore. Each clay robot costs 4 ore. Each obsidian robot costs 3 ore and 20 clay. Each geode robot costs 2 ore and 16 obsidian.
Blueprint 18: Each ore robot costs 3 ore. Each clay robot costs 4 ore. Each obsidian robot costs 4 ore and 18 clay. Each geode robot costs 3 ore and 8 obsidian.
Blueprint 19: Each ore robot costs 3 ore. Each clay robot costs 3 ore. Each obsidian robot costs 2 ore and 14 clay. Each geode robot costs 3 ore and 17 obsidian.
Blueprint 20: Each ore robot costs 2 ore. Each clay robot costs 3 ore. Each obsidian robot costs 3 ore and 11 clay. Each geode robot costs 3 ore and 14 obsidian.
Blueprint 21: Each ore robot costs 3 ore. Each clay robot costs 3 ore. Each obsidian robot costs 3 ore and 6 clay. Each geode robot costs 2 ore and 16 obsidian.
Blueprint 22: Each ore robot costs 2 ore. Each clay robot costs 4 ore. Each obsidian robot costs 4 ore and 20 clay. Each geode robot costs 3 ore and 14 obsidian.
Blueprint 23: Each ore robot costs 4 ore. Each clay robot costs 4 ore. Each obsidian robot costs 3 ore and 10 clay. Each geode robot costs 2 ore and 14 obsidian.
Blueprint 24: Each ore robot costs 4 ore. Each clay robot costs 4 ore. Each obsidian robot costs 2 ore and 7 clay. Each geode robot costs 4 ore and 13 obsidian.
Blueprint 25: Each ore robot costs 3 ore. Each clay robot costs 4 ore. Each obsidian robot costs 4 ore and 18 clay. Each geode robot costs 4 ore and 12 obsidian.
Blueprint 26: Each ore robot costs 4 ore. Each clay robot costs 4 ore. Each obsidian robot costs 4 ore and 11 clay. Each geode robot costs 4 ore and 12 obsidian.
Blueprint 27: Each ore robot costs 4 ore. Each clay robot costs 4 ore. Each obsidian robot costs 4 ore and 9 clay. Each geode robot costs 4 ore and 16 obsidian.
Blueprint 28: Each ore robot costs 4 ore. Each clay robot costs 3 ore. Each obsidian robot costs 3 ore and 7 clay. Each geode robot costs 2 ore and 7 obsidian.
Blueprint 29: Each ore robot costs 4 ore. Each clay robot costs 4 ore. Each obsidian robot costs 2 ore and 14 clay. Each geode robot costs 4 ore and 19 obsidian.
Blueprint 30: Each ore robot costs 4 ore. Each clay robot costs 3 ore. Each obsidian robot costs 4 ore and 20 clay. Each geode robot costs 2 ore and 15 obsidian.`
