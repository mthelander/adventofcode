package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

const Width = 101
const Height = 103
const Seconds = 100

type Robot struct {
	Position Point
	Velocity Point
}

type Point struct {
	Row int
	Col int
}

func (r *Robot) MoveOnce() {
	r.Position = r.Position.Add(r.Velocity)
}

func (p Point) Add(o Point) Point {
	x, y := mod(p.Row+o.Row, Width), mod(p.Col+o.Col, Height)
	return Point{x, y}
}

func (p Point) Gte(a Point) bool {
	return p.Row >= a.Row && p.Col >= a.Col
}

func (p Point) Lt(a Point) bool {
	return p.Row < a.Row && p.Col < a.Col
}

func (p Point) Between(a, b Point) bool {
	return p.Gte(a) && p.Lt(b)
}

func (p Point) Equal(a Point) bool {
	return p.Row == a.Row && p.Col == a.Col
}

func mod(a, b int) int {
	m := a % b
	if a < 0 && b < 0 {
		m -= b
	}
	if a < 0 && b > 0 {
		m += b
	}
	return m
}

func (p Point) String() string {
	return fmt.Sprintf("{%v, %v}", p.Row, p.Col)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func toi(st string) int {
	x, err := strconv.Atoi(st)
	check(err)
	return x
}

func buildPoint(s1, s2 string) Point {
	return Point{toi(s1), toi(s2)}
}

func readLines(path string) []Robot {
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	re := regexp.MustCompile(`p=(\d+),(\d+) v=([0-9-]+),([0-9-]+)`)

	robots := make([]Robot, 0)

	for scanner.Scan() {
		data := re.FindStringSubmatch(scanner.Text())

		robots = append(robots, Robot{
			Position: buildPoint(data[1], data[2]),
			Velocity: buildPoint(data[3], data[4]),
		})
	}
	check(scanner.Err())
	return robots
}

func quadrant(robots []Robot, min, max Point) int {
	count := 0
	for _, r := range robots {
		if r.Position.Between(min, max) {
			count += 1
		}
	}
	return count
}

func countQuadrants(r []Robot) int {
	midrow, midcol := Height/2, Width/2
	quadrants := [][]Point{
		{{0, 0}, {midcol, midrow}},                  // topleft
		{{midcol + 1, 0}, {Width, midrow}},          // topright
		{{0, midrow + 1}, {midcol, Height}},         // bottomleft
		{{midcol + 1, midrow + 1}, {Width, Height}}, // bottomright
	}
	product := 1
	for _, pts := range quadrants {
		product *= quadrant(r, pts[0], pts[1])
	}
	return product
}

func countAt(robots []Robot, pt Point) int {
	count := 0
	for _, r := range robots {
		if r.Position.Equal(pt) {
			count += 1
		}
	}
	return count
}

func main() {
	robots := readLines(os.Args[1])

	for i := range len(robots) {
		for range Seconds {
			robots[i].MoveOnce()
		}
	}
	fmt.Printf("Count: %v\n", countQuadrants(robots))
}
