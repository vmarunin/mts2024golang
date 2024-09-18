package main

// https://leetcode.com/problems/the-two-sneaky-numbers-of-digitville

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hurray!")

	startTime := time.Now()

	out := getSneakyNumbers([]int{0, 1, 1, 0})
	// out := getSneakyNumbers([]int{0,3,2,1,3,2})
	// out := getSneakyNumbers([]int{7,1,5,4,3,4,6,0,9,5,8,2})

	duration := time.Since(startTime)
	fmt.Println(out)
	fmt.Println(duration)

}

func getSneakyNumbers(nums []int) []int {
	cnt := map[int]int{}
	for _, num := range nums {
		cnt[num]++
	}
	ret := []int{}
	for num, count := range cnt {
		if count > 1 {
			ret = append(ret, num)
		}
	}
	return ret
}
