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

func inBounds(pt, max, origin Point) bool {
	return pt.lte(max) && origin.lte(pt)
}

func tryAdd(data map[Point]bool, pt, max, origin Point) {
	if inBounds(pt, max, origin) {
		data[pt] = true
	}
}

func main() {
	data := readLines(os.Args[1])
	antennas := findAntennas(data)
	maxRow := maxBy(data, func(pt Point) int { return pt.X })
	maxCol := maxBy(data, func(pt Point) int { return pt.Y })
	antinodes := make(map[Point]bool)
	max, origin := Point{maxRow, maxCol}, Point{0, 0}

	for _, points := range antennas {
		for _, pt := range points {
			for _, other := range getOthers(pt, points) {
				tryAdd(antinodes, pt, max, origin)
				tryAdd(antinodes, other, max, origin)
				d1, d2 := pt.dist(other), other.dist(pt)
				a1, a2 := pt.add(d1), other.add(d2)

				for ; inBounds(a1, max, origin) || inBounds(a2, max, origin); a1, a2 = a1.add(d1), a2.add(d2) {
					tryAdd(antinodes, a1, max, origin)
					tryAdd(antinodes, a2, max, origin)
				}
			}
		}
	}
	fmt.Printf("%v", len(antinodes))
}
