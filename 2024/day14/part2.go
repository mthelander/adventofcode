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

type Robot struct {
	Position Point
	Velocity Point
}

type Point struct {
	Row int
	Col int
}

func (r *Robot) MoveOnce(grid *map[Point]int) {
	old := r.Position
	r.Position = r.Position.Add(r.Velocity)
	(*grid)[old] -= 1
	(*grid)[r.Position] += 1
}

func (p Point) Add(o Point) Point {
	x, y := mod(p.Row+o.Row, Width), mod(p.Col+o.Col, Height)
	return Point{x, y}
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

func traverse(grid *map[Point]int, pt, dir Point, depth int) bool {
	if depth == 0 {
		return true
	}
	return (*grid)[pt] > 0 && traverse(grid, pt.Add(dir), dir, depth-1)
}

func findTree(grid *map[Point]int) bool {
	right := Point{0, 1}
	for pt, c := range *grid {
		if c > 0 {
			// Look for 20 vertical lines
			if traverse(grid, pt.Add(right), right, 20) {
				return true
			}
		}
	}
	return false
}

func buildGrid(robots []Robot) map[Point]int {
	grid := make(map[Point]int)
	for _, r := range robots {
		grid[r.Position] += 1
	}
	return grid
}

func toString(grid *map[Point]int) string {
	res := ""
	for i := range Width {
		for j := range Height {
			pt := Point{j, i}
			val := strconv.Itoa((*grid)[pt])
			if val == "0" {
				val = "."
			}
			res += val
		}
		res += "\n"
	}
	return res
}

func main() {
	robots := readLines(os.Args[1])
	grid := buildGrid(robots)
	num := 0

	for {
		num += 1
		for i := range len(robots) {
			robots[i].MoveOnce(&grid)
		}

		if findTree(&grid) {
			fmt.Printf("%v\n", toString(&grid))
			fmt.Printf("Num Seconds: %v\n", num)
			break
		}
	}
}
