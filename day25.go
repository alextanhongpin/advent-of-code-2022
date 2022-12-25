package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(part1(input1)) // 2=-1=0
	fmt.Println(part1(input2)) // 20-==01-2-=1-2---1-0
}

var digitBySymbol = map[rune]int{
	'=': -2,
	'-': -1,
	'0': 0,
	'1': 1,
	'2': 2,
}

var symbolByDigit = map[int]rune{
	-2: '=',
	-1: '-',
	0:  '0',
	1:  '1',
	2:  '2',
}

func toSnafu(dec int) string {
	var res []int
	for dec > 0 {
		rem := dec % 5
		res = append(res, rem)
		dec /= 5
	}
	var output []int
	carryForward := 0
	for _, r := range res {
		r += carryForward
		if r <= 2 {
			output = append(output, r)
			carryForward = 0
		} else {
			output = append(output, (r)-5)
			carryForward = 1
		}
	}
	if carryForward != 0 {
		output = append(output, carryForward)
	}
	for i, j := 0, len(output)-1; i < len(output)/2; i, j = i+1, j-1 {
		output[i], output[j] = output[j], output[i]
	}

	snafu := make([]rune, len(output))
	for i, v := range output {
		snafu[i] = symbolByDigit[v]
	}
	return string(snafu)
}

func fromSnafu(snafu string) (dec int) {
	for _, ch := range snafu {
		dec = dec*5 + digitBySymbol[ch]
	}
	return
}

func part1(input string) string {
	var sum int
	for _, row := range strings.Split(input, "\n") {
		sum += fromSnafu(row)
	}
	return toSnafu(sum)
}

var input1 = `1=-0-2
12111
2=0=
21
2=01
111
20012
112
1=-1=
1-12
12
1=
122`

var input2 = `21020=22-2=11002=2
2--221202=102-
22120-
1=0000=1-2=
10=221
110=1121002-10--=
10=00==--121==0=
202=-1221
101-01-
12==-
221101-=110=0
1=-=2-12-2202=
1=0-02-=0102-=-12
1==-0
1=212--00
202
1-11121--21-21=-2
22--2-=-2
2-00-0=---2-==-
1=2=-2-=0--1-22
1=000=0-112
121=-001=0
110-=-020--=0
1=0-1
2=0--==0
102=21=201=-21=21
11-22=-22=1-0=
1011----1----2-1
2===-00022
12-
20=0
10=011=00=0-
20-=0-0-1=2====0-2
2000-01011=001=
1=0-01
2=-1
12=10--=-1-
2=21=12=
22=21-0
1==-02202==--
1===
210020-122222202=
102-20211212==11---
1020=221-00=01
1=2=1=0---10
111-
1=02
10-10=-2-0=11
1002-
1-
2=12112-
20=202-=01=2-1
1=-
1=-=222221=1102012
20-1=02=--
2=2-=
1=1
1-0==1122-02==2=
1-20-0-=0=0
1=022-11-122-0010
1021===1
1=2020--
2=0==00-==0011
1--220-0=-=1-2=01---
2=1=021010=2=1=-
21=2-==22-==0=-010
200-=1
212-
21---=-022-22
1121211-110-00-10
22-=1
1=-=2-2100--=1=02
1=0
11102-202==122-=00
2=1-=
1=11--00
211-012--2
1=1021012100
12-101==-110102
1==--2-==-==1-=-
1===21-
2---1100112010202
1200-112=11=-0
1220--00-21--11
12-=0-00-=
100122==---=0
2221-==-210101==
202-1=1
100==12101=2=
102
210=20=11
12-=-=2--1-=---=20=
2222001
12=0
2--=20-10--10
2-1020000-1=1-
1-1
1101-111-
1==201=--==
2=
2=0=
12=
11=
1==
1-10000-2=
21=-21
10-=--20-2-00===
202122-
1=02--10011
1==-0-1=2111-12=2-2
20
1==-
2=-0-=222=-02=
1110=1-=1022-
122=2==0=2222--12
10111000-=2--=0--0
1-22
2=10-0=--21
1===-
1=022-
1=2===2-
22211=--10120120=0
101--=111101-=11
1212=10
20100
12=112=11==01
202111=--2===0`