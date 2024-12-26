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

type Region struct {
	Id        int
	Area      int
	Perimeter int
	Char      string
	Points    []Point
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

func allDirections() []Point {
	return []Point{
		Point{-1, 0}, // up
		Point{1, 0},  // down
		Point{0, -1}, // left
		Point{0, 1},  // right
	}
}

func contains(pt Point, points *[]Point) bool {
	for _, p := range *points {
		if p.equal(pt) {
			return true
		}
	}
	return false
}

func getPointsInRegion(data map[Point]string, curr Point, points *[]Point) {
	*points = append(*points, curr)
	for _, dir := range allDirections() {
		next := curr.add(dir)
		nextVal, ok := data[next]
		if ok && nextVal == data[curr] && !contains(next, points) {
			getPointsInRegion(data, next, points)
		}
	}
}

func calculatePerimeter(data map[Point]string, points []Point) int {
	perimeter := 0
	for _, pt := range points {
		for _, dir := range allDirections() {
			next := pt.add(dir)
			nextVal, ok := data[next]
			if !(ok && nextVal == data[pt]) {
				perimeter += 1
			}
		}
	}
	return perimeter
}

func calculateArea(data map[Point]string, points []Point) int {
	return len(points)
}

func classifyRegions(data map[Point]string) []Region {
	regionMap := make(map[Point]int, len(data))
	regionNum := 0
	regions := make([]Region, 0)

	for k, v := range data {
		_, ok := regionMap[k]
		if !ok {
			pointsInRegion := make([]Point, 0)
			getPointsInRegion(data, k, &pointsInRegion)
			regions = append(regions, Region{
				Id:        regionNum,
				Char:      v,
				Points:    pointsInRegion,
				Area:      calculateArea(data, pointsInRegion),
				Perimeter: calculatePerimeter(data, pointsInRegion),
			})
			for _, pt := range pointsInRegion {
				// these have already been categorized; skip then when we get there
				regionMap[pt] = regionNum
			}
			regionNum += 1
		}
	}

	return regions
}

func main() {
	data := readLines(os.Args[1])
	regions := classifyRegions(data)
	totalPrice := 0
	for _, reg := range regions {
		totalPrice += reg.Area * reg.Perimeter
		fmt.Printf("%v: A/%v, P/%v\n", reg.Char, reg.Area, reg.Perimeter)
	}
	fmt.Printf("%v\n", totalPrice)
}
