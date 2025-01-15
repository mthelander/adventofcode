package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Space map[complex64]string

const Wall = "#"
const Path = "O"
const Empty = "."

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readLines(path string) []complex64 {
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	result := make([]complex64, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		d := strings.Split(scanner.Text(), ",")
		result = append(result, parseComplex(d))
	}

	check(scanner.Err())
	return result
}

func parseComplex(d []string) complex64 {
	a, b := atoi(d[0]), atoi(d[1])
	return cpx(a, b)
}

func atoi(val string) int {
	i, err := strconv.Atoi(val)
	check(err)
	return i
}

func cpx(i, j int) complex64 {
	return complex(float32(i), float32(j))
}

func buildGrid(size int) Space {
	grid := Space{}
	for row := range size {
		for col := range size {
			grid[cpx(col, row)] = Empty
		}
	}
	return grid
}

func prettyPrint(grid Space) {
	n := sizeOf(grid)
	st := ""
	for i := range n {
		for j := range n {
			c := cpx(j, i)
			st += grid[c]
		}
		st += "\n"
	}
	fmt.Print(st)
}

func setAll(grid *Space, points []complex64, val string) {
	for _, b := range points {
		(*grid)[b] = val
	}
}

func sizeOf(grid Space) int {
	return int(math.Sqrt(float64(len(grid))))
}

func getPath(grid Space) []complex64 {
	n := sizeOf(grid)
	start, goal := cpx(0, 0), cpx(n-1, n-1)

	prev := map[complex64]complex64{}
	visited := map[complex64]bool{}
	queue := []complex64{start}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if !visited[curr] {
			visited[curr] = true
			for _, n := range neighbors(grid, curr) {
				if !visited[n] && grid[n] != "#" {
					prev[n] = curr
					queue = append(queue, n)
				}
			}
		}
	}

	return pathFromPrev(prev, goal)
}

func neighbors(grid Space, pt complex64) []complex64 {
	nbs := make([]complex64, 0)
	dirs := []complex64{1, -1, 1i, -1i}

	for _, d := range dirs {
		nb := pt + d
		if _, ok := grid[nb]; ok {
			nbs = append(nbs, nb)
		}
	}
	return nbs
}

func pathFromPrev(prev map[complex64]complex64, curr complex64) []complex64 {
	result := []complex64{}
	for {
		if next, ok := prev[curr]; ok {
			result = append(result, next)
			curr = next
		} else {
			return result
		}
	}
}

func main() {
	bytes := readLines(os.Args[1])
	size := atoi(os.Args[2])

	grid := buildGrid(size)
	setAll(&grid, bytes, Wall)

	slices.Reverse(bytes)

	for i, b := range bytes {
		grid[b] = Empty

		path := getPath(grid)
		if len(path) > 0 {
			fmt.Printf("Result: %v: %v\n", i, b)
			break
		}
	}
}
