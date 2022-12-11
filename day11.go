package main

import (
	"fmt"
	"log"
	"sort"
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

func part1(input string) int {
	monkeys := parse(input)
	worryLevel := 3
	rounds := 20

	for r := 0; r < rounds; r++ {
		for i, m := range monkeys {
			for _, item := range m.items {
				j, n := m.inspect(item, worryLevel)
				monkeys[j].items = append(monkeys[j].items, n)
			}
			monkeys[i].inspected += len(m.items)
			monkeys[i].items = []int{}
		}
	}

	counts := make([]int, len(monkeys))
	for i, m := range monkeys {
		counts[i] = m.inspected
	}
	sort.Ints(counts)
	return counts[len(counts)-1] * counts[len(counts)-2]
}

func part2(input string) int {
	monkeys := parse(input)
	worryLevel := 1
	rounds := 10_000

	gcd := 1
	for _, m := range monkeys {
		gcd = gcd * m.test
	}

	for r := 0; r < rounds; r++ {
		for i, m := range monkeys {
			for _, item := range m.items {
				j, n := m.inspect(item, worryLevel)
				monkeys[j].items = append(monkeys[j].items, n%gcd)
			}
			monkeys[i].inspected += len(m.items)
			monkeys[i].items = []int{}
		}
	}

	counts := make([]int, len(monkeys))
	for i, m := range monkeys {
		counts[i] = m.inspected
	}
	sort.Ints(counts)
	return counts[len(counts)-1] * counts[len(counts)-2]
}

func parse(input string) []monkey {
	var monkeys []monkey
	var m monkey
	rows := strings.Split(input, "\n")

	for _, row := range rows {
		row = strings.TrimSpace(row)
		if row == "" {
			monkeys = append(monkeys, m)
			m = monkey{}
			continue
		}
		if strings.HasPrefix(row, "Monkey") {
			continue
		} else if strings.HasPrefix(row, "Starting") {
			row = strings.TrimPrefix(row, "Starting items: ")
			items := strings.Split(row, ", ")
			m.items = make([]int, len(items))
			for i, item := range items {
				m.items[i] = toInt(item)
			}
		} else if strings.HasPrefix(row, "Operation:") {
			parts := strings.Fields(row)
			if opVal := parts[len(parts)-1]; opVal == "old" {
				m.opOld = true
			} else {
				m.opVal = toInt(parts[len(parts)-1])
			}
			m.op = parts[len(parts)-2]
		} else if strings.HasPrefix(row, "Test:") {
			parts := strings.Fields(row)
			m.test = toInt(parts[len(parts)-1])
		} else if strings.HasPrefix(row, "If true:") {
			parts := strings.Fields(row)
			m.yes = toInt(parts[len(parts)-1])
		} else if strings.HasPrefix(row, "If false:") {
			parts := strings.Fields(row)
			m.no = toInt(parts[len(parts)-1])
		} else {
			log.Fatal("invalid command:", row)
		}
	}
	monkeys = append(monkeys, m)

	return monkeys
}

func toInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

type monkey struct {
	items     []int
	op        string
	opVal     int
	opOld     bool
	test      int
	yes       int
	no        int
	inspected int
}

func (m monkey) inspect(item int, worryLevel int) (nextMonkey int, n int) {
	switch m.op {
	case "+":
		if m.opOld {
			n = item + item
		} else {
			n = item + m.opVal
		}
	case "*":
		if m.opOld {
			n = item * item
		} else {
			n = item * m.opVal
		}
	}
	n = n / worryLevel
	if n%m.test == 0 {
		return m.yes, n
	}
	return m.no, n
}

var input1 = `Monkey 0:
  Starting items: 79, 98
  Operation: new = old * 19
  Test: divisible by 23
    If true: throw to monkey 2
    If false: throw to monkey 3

Monkey 1:
  Starting items: 54, 65, 75, 74
  Operation: new = old + 6
  Test: divisible by 19
    If true: throw to monkey 2
    If false: throw to monkey 0

Monkey 2:
  Starting items: 79, 60, 97
  Operation: new = old * old
  Test: divisible by 13
    If true: throw to monkey 1
    If false: throw to monkey 3

Monkey 3:
  Starting items: 74
  Operation: new = old + 3
  Test: divisible by 17
    If true: throw to monkey 0
    If false: throw to monkey 1`

var input2 = `Monkey 0:
  Starting items: 83, 97, 95, 67
  Operation: new = old * 19
  Test: divisible by 17
    If true: throw to monkey 2
    If false: throw to monkey 7

Monkey 1:
  Starting items: 71, 70, 79, 88, 56, 70
  Operation: new = old + 2
  Test: divisible by 19
    If true: throw to monkey 7
    If false: throw to monkey 0

Monkey 2:
  Starting items: 98, 51, 51, 63, 80, 85, 84, 95
  Operation: new = old + 7
  Test: divisible by 7
    If true: throw to monkey 4
    If false: throw to monkey 3

Monkey 3:
  Starting items: 77, 90, 82, 80, 79
  Operation: new = old + 1
  Test: divisible by 11
    If true: throw to monkey 6
    If false: throw to monkey 4

Monkey 4:
  Starting items: 68
  Operation: new = old * 5
  Test: divisible by 13
    If true: throw to monkey 6
    If false: throw to monkey 5

Monkey 5:
  Starting items: 60, 94
  Operation: new = old + 5
  Test: divisible by 3
    If true: throw to monkey 1
    If false: throw to monkey 0

Monkey 6:
  Starting items: 81, 51, 85
  Operation: new = old * old
  Test: divisible by 5
    If true: throw to monkey 5
    If false: throw to monkey 1

Monkey 7:
  Starting items: 98, 81, 63, 65, 84, 71, 84
  Operation: new = old + 3
  Test: divisible by 2
    If true: throw to monkey 2
    If false: throw to monkey 3`
