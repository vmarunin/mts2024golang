package main

import (
	"fmt"
	"time"
)

// Вычислим остаток от деления n-го числа Фибоначчи на 10^9+7.
// вычислять будем рекурсивно с мемоизацией

const MOD = 1_000_000_007

func main() {
	ts := time.Now()
	out := fib(4_000_000)
	duration := time.Since(ts)
	fmt.Println(out)
	fmt.Println(duration)
}

func fib(n int) int {
	cache := map[int]int{}

	var recurrentFunc func(n int) int
	recurrentFunc = func(n int) int {
		if n < 2 {
			return n
		}
		if v, ok := cache[n]; ok {
			return v
		}
		v := (recurrentFunc(n-1) + recurrentFunc(n-2)) % MOD
		cache[n] = v
		return v
	}
	return recurrentFunc(n)
}
