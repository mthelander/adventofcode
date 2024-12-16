package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"slices"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func numDigits(i int64) int64 {
	return int64(math.Log10(float64(i))) + 1
}

func evenDigits(i int64) bool {
	return (numDigits(i) % 2) == 0
}

func readLines(path string) []int64 {
	file, err := os.Open(path)
	check(err)
	defer file.Close()
	rg := regexp.MustCompile("[ ]+")
	result := make([]int64, 0)
	scanner := bufio.NewScanner(file)

	if scanner.Scan() {
		data := rg.Split(scanner.Text(), -1)
		for _, str := range data {
			i64, err := strconv.ParseInt(str, 10, 64)
			check(err)
			result = append(result, i64)
		}
	}
	check(scanner.Err())
	return result
}

func calculateSides(val int64) (int64, int64) {
	n := numDigits(val)
	exp := int64(math.Pow(10, float64(n/2)))
	left := val / exp
	right := val % exp
	return left, right
}

func blink(line []int64) []int64 {
	for i := 0; i < len(line); i++ {
		val := line[i]
		if val == 0 {
			line[i] = 1
		} else if evenDigits(val) {
			left, right := calculateSides(val)
			line[i] = left
			line = slices.Insert(line, i+1, right)
			i++
		} else {
			line[i] *= 2024
		}
	}
	return line
}

func main() {
	data := readLines(os.Args[1])

	numBlinks, err := strconv.ParseInt(os.Args[2], 10, 32)
	check(err)

	for range numBlinks {
		data = blink(data)
	}

	fmt.Printf("%v\n", len(data))
}
