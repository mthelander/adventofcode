package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type MinHeap struct {
	list   *[]Point
	dist   *map[Point]int
	pindex *map[Point]int
}

func (h MinHeap) Len() int {
	return len(*h.list)
}

func (h MinHeap) Less(i, j int) bool {
	a, b := (*h.list)[i], (*h.list)[j]
	dst := *h.dist
	return dst[a] < dst[b]
}

func (h *MinHeap) Swap(i, j int) {
	p1, p2 := (*h.list)[i], (*h.list)[j]
	(*h.list)[i], (*h.list)[j] = p2, p1
	(*h.pindex)[p1], (*h.pindex)[p2] = j, i
}

func (h *MinHeap) Push(x any) {
	(*h.pindex)[x.(Point)] = h.Len()
	(*h.list) = append(*h.list, x.(Point))
}

func (h *MinHeap) Pop() any {
	old := *h.list
	n := len(old)
	x := old[n-1]
	(*h.list) = old[:n-1]
	delete(*h.pindex, x)
	return x
}

type Point struct {
	Row int
	Col int
}

type Maze struct {
	Grid    map[Point]string
	NumRows int
	NumCols int
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

func (p Point) Direction(o Point) Point {
	if p.Row < o.Row {
		return Down()
	}
	if p.Row > o.Row {
		return Up()
	}
	if p.Col > o.Col {
		return Left()
	}
	return Right()
}

func (p Point) String() string {
	return fmt.Sprintf("{%v, %v}", p.Row, p.Col)
}

func (m *Maze) IsClear(pt Point) bool {
	return Get(m, pt) != Wall
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
		Right(),
		Left(),
		Up(),
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

func Get(m *Maze, pt Point) string {
	val, ok := m.Grid[pt]
	if ok {
		return val
	}
	return Wall
}

func isPath(m *Maze, p1, p2 Point) (int, bool) {
	for _, d := range directions() {
		cost := 0
		for n := p1.Add(d); m.IsClear(n); n = n.Add(d) {
			cost += 1
			if n.Equals(p2) {
				return cost, true
			}
		}
	}
	return 0, false
}

func buildTurns(m *Maze) []Point {
	turns := make([]Point, 0)
	for i := range m.NumRows + 1 {
		for j := range m.NumCols + 1 {
			pt := Point{i, j}
			if val := Get(m, pt); val != Wall {
				rl := m.IsClear(pt.Add(Right())) || m.IsClear(pt.Add(Left()))
				ud := m.IsClear(pt.Add(Up())) || m.IsClear(pt.Add(Down()))

				if (rl && ud) || val == Start || val == End {
					turns = append(turns, pt)
				}
			}
		}
	}
	return turns
}

func inMap(c Point, m *map[Point]Point) bool {
	_, ok := (*m)[c]
	return ok
}

func fillData(nodes *[]Point) (*[]Point, *map[Point]int, *map[Point]int) {
	data := make([]Point, 0)
	pindex := map[Point]int{}
	dist := map[Point]int{}

	for i, k := range *nodes {
		dist[k] = 1_000_000 // infinity
		pindex[k] = i
		data = append(data, k)
	}

	return &data, &dist, &pindex
}

func spf(m *Maze, start, goal Point) int {
	vertices, turns := buildVertices(m)
	prev := map[Point]Point{}
	data, dist, pindex := fillData(&turns)
	(*dist)[start] = 0

	pq := &MinHeap{data, dist, pindex}
	heap.Init(pq)

	for pq.Len() > 0 {
		current := heap.Pop(pq).(Point)
		for nb, cost := range getNeighbors(current, vertices) {
			idx, inheap := (*pindex)[nb]
			if inheap {
				subpath := reconstructPath(start, nb, prev)
				alt := (*dist)[current] + cost + len(subpath)*1000
				if alt < (*dist)[nb] {
					(*dist)[nb] = alt
					prev[nb] = current
					heap.Fix(pq, idx)
				}
			}
		}
	}

	//printWithPath(m, start, goal, prev)

	return adjustScore(m, start, goal, prev, (*dist)[goal])
}

func adjustScore(m *Maze, start, goal Point, prev map[Point]Point, score int) int {
	path := reconstructPath(start, goal, prev)
	dir := path[len(path)-1].Direction(path[len(path)-2])
	if dir.Equals(Right()) {
		// No turn was necessary at the start, so the cost is 1 turn over
		return score - 1000
	}
	return score
}

func scorePath(m *Maze, start, goal Point, path []Point) int {
	cost := 0

	for i := 0; i < len(path)-1; i++ {
		a, b := path[i], path[i+1]
		c, _ := isPath(m, a, b)
		cost += c
	}

	return cost
}

func printWithPath(m *Maze, start, goal Point, prev map[Point]Point) {
	sp := reconstructPath(start, goal, prev)
	c := 0
	for i := range len(sp) - 1 {
		a, b := sp[i], sp[i+1]
		dir := b.Direction(a)
		ch := dirToChar(dir)
		for x := b; !x.Equals(a); x = x.Add(dir) {
			c += 1
			m.Grid[x] = ch
		}
	}
	m.Grid[start] = Start
	prettyPrint(m)
	fmt.Printf("path size: %v, %v\n", len(sp), c)
	fmt.Printf("num steps: %v\n", scorePath(m, start, goal, sp))
	fmt.Printf("path: %v\n", sp)
}

func dirToChar(dir Point) string {
	switch dir {
	case Up():
		return "^"
	case Down():
		return "v"
	case Left():
		return "<"
	case Right():
		return ">"
	default:
		return ""
	}
}

func reconstructPath(start, goal Point, prev map[Point]Point) []Point {
	path := make([]Point, 0)
	for n := goal; inMap(n, &prev); n = prev[n] {
		path = append(path, n)
	}
	return append(path, start)
}

func getNeighbors(pt Point, data map[[2]Point]int) map[Point]int {
	nbs := map[Point]int{}
	for k, c := range data {
		if k[0].Equals(pt) {
			nbs[k[1]] = c
		}
	}
	return nbs
}

func buildVertices(m *Maze) (map[[2]Point]int, []Point) {
	nodes := buildTurns(m)
	paths := map[[2]Point]int{}

	for _, t1 := range nodes {
		for _, t2 := range nodes {
			if !t1.Equals(t2) {
				if cost, ok := isPath(m, t1, t2); ok {
					paths[[2]Point{t1, t2}] = cost
				}
			}
		}
	}
	return paths, nodes
}

func walkMaze(m *Maze) int {
	goal := positionOf(End, m)
	start := positionOf(Start, m)

	return spf(m, start, goal)
}

func prettyPrint(m *Maze) {
	st := ""
	for i := range m.NumRows + 1 {
		for j := range m.NumCols + 1 {
			st += m.Grid[Point{i, j}]
		}
		st += "\n"
	}
	fmt.Println(st)
}

func main() {
	m := readLines(os.Args[1])
	fmt.Printf("Shortest path: %v\n", walkMaze(m))
}
