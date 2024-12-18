package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"time"
)

type Num struct {
	Val   int64
	Right *Num
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

func readLines(path string) *Num {
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
	return buildLL(result)
}

func buildLL(list []int64) *Num {
	if len(list) == 0 {
		return &Num{}
	}
	return &Num{list[0], buildLL(list[1:])}
}

func calculateSides(val int64) (int64, int64) {
	n := numDigits(val)
	exp := int64(math.Pow(10, float64(n/2)))
	left := val / exp
	right := val % exp
	return left, right
}

func blink(head *Num) {
	curr := head
	for curr.Right != nil {
		if curr.Val == 0 {
			curr.Val = 1
			curr = curr.Right
		} else if evenDigits(curr.Val) {
			next := curr.Right
			left, right := calculateSides(curr.Val)
			curr.Val = left
			curr.Right = &Num{right, next}
			curr = next
		} else {
			curr.Val *= 2024
			curr = curr.Right
		}
	}
}

func printList(head *Num) {
	val := *head
	for val.Right != nil {
		fmt.Printf("val: %v\n", val.Val)
		val = *val.Right
	}
}

func length(head *Num) int64 {
	val := *head
	count := int64(0)
	for ; val.Right != nil; count++ {
		val = *val.Right
	}
	return count
}

func main() {
	head := readLines(os.Args[1])

	numBlinks, err := strconv.ParseInt(os.Args[2], 10, 32)
	check(err)
	start := time.Now()

	for i := range numBlinks {
		elapsed := time.Since(start)
		start = time.Now()
		fmt.Printf("i: %v (%v)\n", i, elapsed)
		blink(head)
	}

	fmt.Printf("%v\n", length(head))
}
