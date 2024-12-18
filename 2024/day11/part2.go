package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
)

type Key struct {
	Round int64
	Left  int64
	Right int64
}

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

func lookup(round, left, right int64, cache *map[Key]int64) int64 {
	k := Key{Left: left, Right: right, Round: round}
	val, ok := (*cache)[k]
	if ok {
		return val
	}
	res := blink(round, []int64{left, right}, cache)
	(*cache)[k] = res
	return res
}

func blink(round int64, line []int64, cache *map[Key]int64) int64 {
	count := int64(len(line) - 1)

	for i := 0; i < len(line); i++ {
		for n := range round {
			val := line[i]
			if val == 0 {
				line[i] = 1
			} else if evenDigits(val) {
				left, right := calculateSides(val)
				line[i] = left
				res := lookup(round-(n+1), left, right, cache)
				count += res
				break
			} else {
				line[i] *= 2024
			}
		}
	}

	return count
}

func main() {
	data := readLines(os.Args[1])

	numBlinks, err := strconv.ParseInt(os.Args[2], 10, 32)
	check(err)

	cache := make(map[Key]int64)

	c := blink(numBlinks, data, &cache) + 1

	fmt.Printf("%v\n", c)
}
