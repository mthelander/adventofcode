package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"slices"
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

func spf(m *Maze, start, goal Point) ([]Point, map[Point]int, map[Point]Point) {
	vertices, turns := buildVertices(m)
	prev := map[Point]Point{}
	data, dist, pindex := fillData(turns)
	(*dist)[start] = 0

	pq := &MinHeap{data, dist, pindex}
	heap.Init(pq)

	for pq.Len() > 0 {
		current := heap.Pop(pq).(Point)
		nbs := getNeighbors(current, vertices)
		for nb, cost := range *nbs {
			idx, inheap := (*pindex)[nb]
			if inheap {
				alt := (*dist)[current] + cost + (numSteps(start, nb, prev) * 1000)
				prevbest := (*dist)[nb]
				if alt <= prevbest {
					(*dist)[nb] = alt
					prev[nb] = current
					heap.Fix(pq, idx)
				}
			}
		}
	}

	path := nodesToPath(m, start, goal, reconstructPath(start, goal, prev))
	return path, *dist, prev
}

func nodesToPath(m *Maze, start, goal Point, nodes []Point) []Point {
	slices.Reverse(nodes)
	return walkPath(m, start, goal, nodes)
}

func numSteps(start, goal Point, prev map[Point]Point) int {
	num := 0
	for n := goal; inMap(n, &prev); n = prev[n] {
		num += 1
	}
	return num + 1
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

func reconstructPath(start, goal Point, prev map[Point]Point) []Point {
	path := make([]Point, 0)
	for n := goal; inMap(n, &prev); n = prev[n] {
		path = append(path, n)
	}
	return append(path, start)
}

func getNeighbors(pt Point, data *map[[2]Point]int) *map[Point]int {
	nbs := map[Point]int{}
	for k, c := range *data {
		if k[0].Equals(pt) {
			nbs[k[1]] = c
		}
	}
	return &nbs
}

func buildVertices(m *Maze) (*map[[2]Point]int, *[]Point) {
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
	return &paths, &nodes
}

func walkPath(m *Maze, start, goal Point, nodes []Point) []Point {
	if len(nodes) < 1 {
		return []Point{}
	}
	next := nodes[0]
	dir := start.Direction(next)
	result := make([]Point, 0)
	for curr := start; !curr.Equals(next); curr = curr.Add(dir) {
		result = append(result, curr)
	}

	return append(result, walkPath(m, next, goal, nodes[1:])...)
}

func calculateCost(path []Point) int {
	cost := 1
	lastdir := Right()

	for i := 0; i < len(path)-1; i++ {
		a, b := path[i], path[i+1]
		cost += 1
		dir := a.Direction(b)
		if !dir.Equals(lastdir) {
			cost += 1000
		}
		lastdir = dir
	}

	return cost
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func setPoints(points []Point, out *map[Point]bool) {
	for _, p := range points {
		(*out)[p] = true
	}
}

func walkMaze(m *Maze) int {
	goal := positionOf(End, m)
	start := positionOf(Start, m)

	path1, dist1, prev1 := spf(m, start, goal)
	path2, dist2, prev2 := spf(m, goal, start)

	cost1 := calculateCost(path1) // start to end
	cost2 := calculateCost(path2) // end to start
	offset := abs(cost2 - cost1)  // to account for different start positions

	nodes := map[Point]bool{}
	setPoints(path1, &nodes)
	setPoints(path2, &nodes)

	for k, v := range dist2 {
		// if the cost to start + cost to end == cost,
		// then we're on the shortest path
		x := dist1[k]
		if v > 0 && x > 0 && (v+x)-offset == cost1 {
			setPoints(nodesToPath(m, start, k, reconstructPath(start, k, prev1)), &nodes)
			setPoints(nodesToPath(m, goal, k, reconstructPath(goal, k, prev2)), &nodes)
		}
	}
	return len(nodes)
}

func main() {
	m := readLines(os.Args[1])
	fmt.Printf("Shortest path: %v\n", walkMaze(m))
}
