package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct {
	X int
	Y int
}

func (p Point) equal(other Point) bool {
	return p.X == other.X && p.Y == other.Y
}

func (p Point) negate() Point {
	return Point{-p.X, -p.Y}
}

func (p Point) add(other Point) Point {
	return Point{p.X + other.X, p.Y + other.Y}
}

func (p Point) dist(other Point) Point {
	return Point{p.X - other.X, p.Y - other.Y}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readLines(path string) map[Point]string {
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	data := make(map[Point]string)
	scanner := bufio.NewScanner(file)

	for x := 0; scanner.Scan(); x++ {
		for y, char := range scanner.Text() {
			pt := Point{x, y}
			data[pt] = string(char)
		}
	}
	check(scanner.Err())
	return data
}

func findAntennas(data map[Point]string) map[string][]Point {
	antennas := make(map[string][]Point)

	for key, val := range data {
		if val != "." {
			_, ok := antennas[val]
			if !ok {
				antennas[val] = []Point{key}
			} else {
				antennas[val] = append(antennas[val], key)
			}
		}
	}
	return antennas
}

func getOthers(pt Point, antennas []Point) []Point {
	result := []Point{}
	for _, other := range antennas {
		if !other.equal(pt) {
			result = append(result, other)
		}
	}
	return result
}

func maxBy(pointsMap map[Point]string, mapper func(Point) int) int {
	var max int
	for point, _ := range pointsMap {
		val := mapper(point)
		if val > max {
			max = val
		}
	}
	return max
}

func addIfInBounds(data map[Point]bool, maxRow, maxCol int, pt Point) {
	if pt.X <= maxRow && pt.Y <= maxCol && pt.X >= 0 && pt.Y >= 0 {
		data[pt] = true
	}
}

func main() {
	data := readLines(os.Args[1])
	antennas := findAntennas(data)
	maxRow := maxBy(data, func(pt Point) int { return pt.X })
	maxCol := maxBy(data, func(pt Point) int { return pt.Y })
	antinodes := make(map[Point]bool)

	for _, points := range antennas {
		for _, pt := range points {
			for _, other := range getOthers(pt, points) {
				antinode1 := pt.add(pt.dist(other))
				antinode2 := other.add(other.dist(pt))
				addIfInBounds(antinodes, maxRow, maxCol, antinode1)
				addIfInBounds(antinodes, maxRow, maxCol, antinode2)
			}
		}
	}
	fmt.Printf("%v", len(antinodes))
}
