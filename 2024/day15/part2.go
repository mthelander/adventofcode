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
	NumCols int
}

const Robot = "@"
const Wall = "#"
const OneBox = "O"
const LBox = "["
const RBox = "]"
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
			col := 0
			for _, char := range line {
				var c1, c2 string
				switch val := string(char); val {
				case Wall:
					c1, c2 = Wall, Wall
				case OneBox:
					c1, c2 = LBox, RBox
				case Empty:
					c1, c2 = Empty, Empty
				case Robot:
					c1, c2 = Robot, Empty
				}
				pt1 := Point{row, col}
				pt2 := Point{row, col + 1}
				warehouse.Layout[pt1] = c1
				warehouse.Layout[pt2] = c2
				col += 2
				warehouse.NumCols = col
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

func Left() Point {
	return Point{0, -1}
}
func Right() Point {
	return Point{0, 1}
}
func Up() Point {
	return Point{-1, 0}
}
func Down() Point {
	return Point{1, 0}
}

func pointFromChar(mv string) Point {
	switch mv {
	case "<":
		return Left()
	case ">":
		return Right()
	case "^":
		return Up()
	default:
		return Down()
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

func canMove(w *Warehouse, pt, dir Point) bool {
	switch next := pt.Add(dir); w.Get(next) {
	case Empty:
		return true
	case LBox, RBox:
		if isLR(dir) {
			return canMove(w, next, dir)
		} else {
			otherside := next.Add(Left())
			if w.Get(next) == LBox {
				otherside = next.Add(Right())
			}
			return canMove(w, next, dir) && canMove(w, otherside, dir)
		}
	}
	return false
}

func shift(w *Warehouse, pt, next Point, val string) {
	w.Update(next, val)
	w.Update(pt, Empty)
}

func movingTwoBoxes(w *Warehouse, pt, dir Point) (Point, bool) {
	if !isLR(dir) {
		switch val := w.Get(pt); val {
		case LBox:
			return pt.Add(Right()), true
		case RBox:
			return pt.Add(Left()), true
		}
	}
	return pt, false
}

func movePoint(w *Warehouse, pt, dir Point) bool {
	// TODO: this is overly complicated, but it works :shrug:
	val := w.Get(pt)
	switch next := pt.Add(dir); w.Get(next) {
	case Empty:
		shift(w, pt, next, val)
		return true
	case LBox, RBox:
		if isLR(dir) {
			if movePoint(w, next, dir) {
				shift(w, pt, next, val)
				return true
			}
		} else {
			otherside, ok := movingTwoBoxes(w, next, dir)
			if ok && canMove(w, next, dir) && canMove(w, otherside, dir) {
				if movePoint(w, otherside, dir) && movePoint(w, next, dir) {
					shift(w, pt, next, val)
					return true
				}
			}
		}
	}
	return false
}

func isLR(dir Point) bool {
	return dir.Row == 0 && abs(dir.Col) == 1
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func prettyPrint(w Warehouse) {
	for row := range w.NumRows {
		for col := range w.NumCols {
			pt := Point{row, col}
			fmt.Printf("%v", string(w.Layout[pt]))
		}
		fmt.Println("")
	}
}

func gps(w *Warehouse) int {
	sum := 0
	for pt, val := range w.Layout {
		if val == LBox {
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
