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

	offset := 10_000_000_000_000

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
			Prize:   Pair{toi(c[1]) + offset, toi(c[2]) + offset},
		})
	}
	check(scanner.Err())
	return results
}

func solve(s Scenario) int {
	a, b := calc(s.ButtonA.X, s.ButtonB.X, s.Prize.X, s.ButtonA.Y, s.ButtonB.Y, s.Prize.Y)
	return a*3 + b
}

func verify(x, y, z int) bool {
	return (x + y) == z
}

func calc(x1, y1, z1, x2, y2, z2 int) (int, int) {
	// do some algebra
	r1 := ((y1 * z2) - (y2 * z1)) / ((y1 * x2) - (y2 * x1))
	r2 := (z2 - (x2 * r1)) / y2

	// check to account for int conversion
	if verify(r1*x1, r2*y1, z1) && verify(r1*x2, r2*y2, z2) {
		return r1, r2
	}

	return 0, 0
}

func main() {
	scenarios := readLines(os.Args[1])
	tokens := 0

	for _, sc := range scenarios {
		tokens += solve(sc)
	}
	fmt.Printf("Total Tokens: %v\n", tokens)
}
