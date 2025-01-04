package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

type Point struct {
	Row int
	Col int
}

type Path struct {
	Steps []Point
	Cost  int
}

type Maze struct {
	Grid    map[Point]string
	NumRows int
	NumCols int
}

type Node struct {
	Loc  Point
	Cost int
}

const Wall = "#"
const Start = "S"
const End = "E"
const Empty = "."

func (p Point) Add(o Point) Point {
	return Point{p.Row + o.Row, p.Col + o.Col}
}

func (p Point) Equals(o Point) bool {
	return p.Row == o.Row && p.Col == o.Col
}

func (p Point) String() string {
	return fmt.Sprintf("{%v, %v}", p.Row, p.Col)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readLines(path string) *Maze {
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	maze := Maze{map[Point]string{}, 0, 0}

	for row := 0; scanner.Scan(); row++ {
		line := scanner.Text()

		for col, char := range line {
			pt := Point{row, col}
			maze.Grid[pt] = string(char)
			maze.NumCols = col
		}
		maze.NumRows = row
	}
	check(scanner.Err())
	return &maze
}

func positionOf(c string, m *Maze) Point {
	for pt, v := range m.Grid {
		if v == c {
			return pt
		}
	}
	return Point{-1, -1}
}

func directions() []Point {
	return []Point{
		Left(),
		Up(),
		Right(),
		Down(),
	}
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

func rotate(pt Point, dirs *[]Point) Point {
	i := slices.Index(*dirs, pt)
	return (*dirs)[(i+1)%len(*dirs)]
}

func walkFast(m *Maze, dirs *[]Point, costs *[]int, pt, lastdir Point, lowestScore, currentScore int) (int, bool) {
	val, exists := m.Grid[pt]
	if !exists || val == Wall {
		return 0, false
	}

	if val == End {
		fmt.Printf("%v\n", currentScore)
		return currentScore, true
	}

	m.Grid[pt] = Wall

	d := lastdir

	for i := range 4 {
		next := pt.Add(d)
		if Get(m, next) != Wall {
			score := (*costs)[i] + currentScore
			if score < lowestScore {
				if sc, valid := walkFast(m, dirs, costs, next, d, lowestScore, score); valid {
					if sc < lowestScore {
						lowestScore = sc
					}
				}
			}
		}
		d = rotate(d, dirs)
	}

	m.Grid[pt] = val

	return lowestScore, true
}

func Get(m *Maze, pt Point) string {
	val, ok := m.Grid[pt]
	if ok {
		return val
	}
	return Wall
}

func walkMaze(m *Maze) int {
	// plus one since we advance forward in addition to (maybe) turning
	costs := []int{1, 1001, 1001, 1001}
	dirs := directions()
	score, _ := walkFast(m, &dirs, &costs, positionOf(Start, m), Right(), 110_520, 0)
	return score
}

func main() {
	m := readLines(os.Args[1])
	fmt.Printf("Shortest path: %v\n", walkMaze(m))
}
