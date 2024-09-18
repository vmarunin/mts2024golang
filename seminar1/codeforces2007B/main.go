package main

// https://codeforces.com/problemset/problem/2007/B

import (
	"bufio"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	fh, err := os.Open("./test.txt")
	if err != nil {
		fh = os.Stdin
	}

	reader := bufio.NewReaderSize(fh, 1024*1024)
	writer := bufio.NewWriterSize(os.Stdout, 1024*1024)

	str, _ := reader.ReadString('\n')
	fields := strings.Fields(str)
	tCnt, _ := strconv.Atoi(fields[0])

	for t := 0; t < tCnt; t++ {
		str, _ = reader.ReadString('\n')
		fields = strings.Fields(str)
		n, _ := strconv.Atoi(fields[0])
		m, _ := strconv.Atoi(fields[1])

		data := make([]int, 0, n)

		str, _ := reader.ReadString('\n')
		fields := strings.Fields(str)
		for i := 0; i < n; i++ {
			v, _ := strconv.ParseInt(fields[i], 10, 64)
			data = append(data, int(v))
		}

		ops := [][3]int{}
		for i := 0; i < m; i++ {
			str, _ := reader.ReadString('\n')
			fields := strings.Fields(str)
			v0 := +1
			if fields[0] == "-" {
				v0 = -1
			}
			v1, _ := strconv.ParseInt(fields[1], 10, 64)
			v2, _ := strconv.ParseInt(fields[2], 10, 64)
			ops = append(ops, [3]int{int(v0), int(v1), int(v2)})
		}

		sol := task(data, ops)

		for i, v := range sol {
			if i > 0 {
				writer.WriteByte(' ')
			}
			writer.WriteString(strconv.Itoa(v))
		}
		writer.WriteByte('\n')
	}

	writer.Flush()
}

func task(data []int, ops [][3]int) []int {
	v := slices.Max(data)

	result := make([]int, 0, len(ops))

	for _, op := range ops {
		if op[1] <= v && op[2] >= v {
			v += op[0]
		}
		result = append(result, v)
	}

	return result
}
