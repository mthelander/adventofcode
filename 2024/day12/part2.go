package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
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

type Pointmap = map[Point]string

type Region struct {
	Id     int
	Area   int
	Sides  int
	Char   string
	Points []Point
}

type Pointdir struct {
	pt  Point
	dir Point
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readLines(path string) Pointmap {
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	data := make(Pointmap)
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

func directions() []Point {
	return []Point{
		Point{-1, 0}, // up
		Point{0, 1},  // right
		Point{1, 0},  // down
		Point{0, -1}, // left
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

func getPointsInRegion(data Pointmap, curr Point, points *[]Point) {
	*points = append(*points, curr)
	for _, dir := range directions() {
		next := curr.add(dir)
		nextVal, ok := data[next]
		if ok && nextVal == data[curr] && !contains(next, points) {
			getPointsInRegion(data, next, points)
		}
	}
}

func turnLeft(dir Point) Point {
	return turnRight(turnRight(turnRight(dir)))
}

func turnRight(dir Point) Point {
	dirs := directions()
	i := slices.Index(dirs, dir)
	return dirs[(i+1)%len(dirs)]

}

func calculateSides(data Pointmap, points []Point) int {
	seen := make(map[Pointdir]bool, 0)
	pt := points[0]
	result := 0
	target := valAt(data, pt)
	for _, d := range directions() {
		result += countSides(data, d, pt, target, &seen)
	}
	return result
}

func calculateArea(data Pointmap, points []Point) int {
	return len(points)
}

func classifyRegions(data Pointmap) []Region {
	regionMap := make(map[Point]int, len(data))
	regionNum := 0
	regions := make([]Region, 0)

	for k, v := range data {
		_, ok := regionMap[k]
		if !ok {
			pointsInRegion := make([]Point, 0)
			getPointsInRegion(data, k, &pointsInRegion)
			regions = append(regions, Region{
				Id:     regionNum,
				Char:   v,
				Points: pointsInRegion,
				Area:   calculateArea(data, pointsInRegion),
				Sides:  calculateSides(data, pointsInRegion),
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

func countSides(data Pointmap, dir, pt Point, target string, seen *map[Pointdir]bool) int {
	pd := Pointdir{pt, dir}
	_, ok := (*seen)[pd]
	if ok {
		return 0
	}
	(*seen)[pd] = true

	val := valAt(data, pt)
	if val != target {
		markAll(data, target, pt, dir, turnRight(dir), seen)
		markAll(data, target, pt, dir, turnLeft(dir), seen)
		return 1
	}

	result := 0
	for _, d := range directions() {
		next := pt.add(d)
		result += countSides(data, d, next, target, seen)
	}
	return result
}

func valAt(data Pointmap, pt Point) string {
	val, ok := data[pt]
	if !ok {
		return ""
	}
	return val
}

func markAll(d Pointmap, t string, pt, fw, lr Point, seen *map[Pointdir]bool) {
	bw := turnRight(turnRight(fw))
	for x := pt.add(lr); valAt(d, x) != t && valAt(d, x.add(bw)) == t; x = x.add(lr) {
		pd := Pointdir{x, fw}
		(*seen)[pd] = true
	}
}

func main() {
	data := readLines(os.Args[1])

	regions := classifyRegions(data)
	totalPrice := 0
	for _, reg := range regions {
		totalPrice += reg.Area * reg.Sides
	}
	fmt.Printf("%v\n", totalPrice)
}
