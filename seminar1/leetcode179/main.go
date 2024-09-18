package main

// https://leetcode.com/problems/largest-number/description

import (
	"fmt"
	"sort"
	"time"
)

func main() {
	fmt.Println("Hurray!")

	startTime := time.Now()

	out := largestNumber([]int{10, 2})
	// out := largestNumber([]int{3, 30, 34, 5, 9})

	duration := time.Since(startTime)
	fmt.Println(out)
	fmt.Println(duration)

}

type MySortStr []string

func (a MySortStr) Len() int      { return len(a) }
func (a MySortStr) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a MySortStr) Less(i, j int) bool {
	s1 := a[i] + a[j]
	s2 := a[j] + a[i]
	return s1 > s2
}

func largestNumber(nums []int) string {
	s := make([]string, len(nums))
	bSize := 0
	for i, v := range nums {
		s[i] = fmt.Sprint(v)
		bSize += len(s[i])
	}
	sort.Sort(MySortStr(s))
	b := make([]byte, bSize)
	bPtr := 0
	for _, v := range s {
		for i := 0; i < len(v); i++ {
			b[bPtr] = v[i]
			bPtr++
		}
	}
	if b[0] == '0' {
		return "0"
	}
	return string(b)
}
