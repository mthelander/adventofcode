package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Point struct {
	X int
	Y int
}

func (p Point) invert() Point {
	return Point{-p.X, -p.Y}
}

func (p Point) equal(other Point) bool {
	return p.X == other.X && p.Y == other.Y
}

func (p Point) add(other Point) Point {
	return Point{p.X + other.X, p.Y + other.Y}
}

func (p Point) dist(other Point) Point {
	return Point{p.X - other.X, p.Y - other.Y}
}

func (p Point) lte(other Point) bool {
	return p.X <= other.X && p.Y <= other.Y
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readLines(path string) map[Point]int {
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	data := make(map[Point]int)
	scanner := bufio.NewScanner(file)

	for x := 0; scanner.Scan(); x++ {
		for y, char := range scanner.Text() {
			pt := Point{x, y}
			data[pt], err = strconv.Atoi(string(char))
			check(err)
		}
	}
	check(scanner.Err())
	return data
}

func allDirections() []Point {
	return []Point{
		Point{-1, 0}, // up
		Point{1, 0},  // down
		Point{0, -1}, // left
		Point{0, 1},  // right
	}
}

func rateTrailhead(data map[Point]int, pt Point) int {
	val := data[pt]
	if val == 9 {
		return 1
	}
	sum := 0
	for _, d := range allDirections() {
		peek := pt.add(d)
		if data[peek] == val+1 {
			sum += rateTrailhead(data, peek)
		}
	}
	return sum
}

func main() {
	data := readLines(os.Args[1])
	sum := 0
	for k, v := range data {
		if v == 0 {
			sum += rateTrailhead(data, k)
		}
	}
	fmt.Printf("%v\n", sum)
}
