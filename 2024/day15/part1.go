package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Point struct {
	Row int
	Col int
}

type Warehouse struct {
	Layout  map[Point]string
	Moves   []string
	NumRows int
}

const Robot = "@"
const Wall = "#"
const Box = "O"
const Empty = "."

func (w Warehouse) String() string {
	return fmt.Sprintf("Layout: %v, Moves: %v", w.Layout, strings.Join(w.Moves, ","))
}

func (w *Warehouse) Update(pt Point, val string) {
	w.Layout[pt] = val
}

func (w *Warehouse) Get(pt Point) string {
	return w.Layout[pt]
}

func (p Point) Add(o Point) Point {
	return Point{p.Row + o.Row, p.Col + o.Col}
}

func (p Point) String() string {
	return fmt.Sprintf("{%v, %v}", p.Row, p.Col)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readLines(path string) Warehouse {
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	warehouse := Warehouse{
		Layout: make(map[Point]string),
		Moves:  make([]string, 0),
	}

	for row := 0; scanner.Scan(); row++ {
		line := scanner.Text()

		if strings.Contains(line, "#") {
			for col, char := range line {
				pt := Point{row, col}
				warehouse.Layout[pt] = string(char)
			}
			warehouse.NumRows++
		} else if len(line) > 0 {
			for _, char := range line {
				warehouse.Moves = append(warehouse.Moves, string(char))
			}
		}
	}
	check(scanner.Err())
	return warehouse
}

func pointFromChar(mv string) Point {
	switch mv {
	case "<":
		return Point{0, -1}
	case ">":
		return Point{0, 1}
	case "^":
		return Point{-1, 0}
	default:
		return Point{1, 0}
	}
}

func positionOf(c string, w Warehouse) Point {
	for pt, v := range w.Layout {
		if v == c {
			return pt
		}
	}
	return Point{-1, -1}
}

func moveAll(w *Warehouse) {
	for _, mv := range w.Moves {
		moveRobot(w, pointFromChar(mv))
	}
}

func moveRobot(w *Warehouse, dir Point) {
	start := positionOf(Robot, *w)
	movePoint(w, start, dir)
}

func shift(w *Warehouse, pt, next Point, val string) {
	w.Update(next, val)
	w.Update(pt, Empty)
}

func movePoint(w *Warehouse, pt, dir Point) bool {
	val := w.Get(pt)
	switch next := pt.Add(dir); w.Get(next) {
	case Empty:
		shift(w, pt, next, val)
		return true
	case Box:
		if movePoint(w, next, dir) {
			shift(w, pt, next, val)
			return true
		}
	}
	return false
}

func prettyPrint(w Warehouse) {
	n := w.NumRows
	for row := range n {
		for col := range n {
			pt := Point{row, col}
			fmt.Printf("%v", string(w.Layout[pt]))
		}
		fmt.Println("")
	}
}

func gps(w *Warehouse) int {
	sum := 0
	for pt, val := range w.Layout {
		if val == Box {
			sum += (100 * pt.Row) + pt.Col
		}
	}
	return sum
}

func main() {
	w := readLines(os.Args[1])
	moveAll(&w)
	prettyPrint(w)
	g := gps(&w)
	fmt.Printf("gps: %v\n", g)
}
