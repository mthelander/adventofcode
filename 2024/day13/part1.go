package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Pair struct {
	X int
	Y int
}

type Scenario struct {
	ButtonA Pair
	ButtonB Pair
	Prize   Pair
}

func toi(st string) int {
	x, err := strconv.Atoi(st)
	check(err)
	return x
}

func readLines(path string) []Scenario {
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	btna := regexp.MustCompile(`Button A: X\+(\d+), Y\+(\d+)`)
	btnb := regexp.MustCompile(`Button B: X\+(\d+), Y\+(\d+)`)
	prze := regexp.MustCompile(`Prize: X=(\d+), Y=(\d+)`)
	blnk := regexp.MustCompile(`\s*`)

	results := make([]Scenario, 0)

	for x := 0; scanner.Scan(); x++ {
		a := btna.FindStringSubmatch(scanner.Text())
		scanner.Scan()
		b := btnb.FindStringSubmatch(scanner.Text())
		scanner.Scan()
		c := prze.FindStringSubmatch(scanner.Text())
		scanner.Scan()
		blnk.FindStringSubmatch(scanner.Text())

		results = append(results, Scenario{
			ButtonA: Pair{toi(a[1]), toi(a[2])},
			ButtonB: Pair{toi(b[1]), toi(b[2])},
			Prize:   Pair{toi(c[1]), toi(c[2])},
		})
	}
	check(scanner.Err())
	return results
}

func solvesIt(x, y, z, i, j int) bool {
	return (x*i)+(y*j) == z
}

func solve(s Scenario) int {
	for x := range 100 {
		for y := range 100 {
			if solvesIt(s.ButtonA.X, s.ButtonB.X, s.Prize.X, x, y) && solvesIt(s.ButtonA.Y, s.ButtonB.Y, s.Prize.Y, x, y) {
				return x*3 + y
			}
		}
	}
	return 0
}

func main() {
	scenarios := readLines(os.Args[1])
	tokens := 0

	for _, sc := range scenarios {
		tokens += solve(sc)
	}
	fmt.Printf("Total Tokens: %v\n", tokens)
}
