package main

import "fmt"

func main() {
	var number, capacity int
	_, err := fmt.Scanln(&number, &capacity)
	if err != nil {
		fmt.Println(err)
		return
	}
	weight := make([]int, number)
	value := make([]int, number)
	for i := 0; i < number; i++ {
		_, _ = fmt.Scanln(&weight[i])
	}

	for i := 0; i < number; i++ {
		_, _ = fmt.Scanln(&value[i])
	}
}

func pa(weight, value []int, capacity int) int {
	res := make([][]int, len(weight))

	for i := 0; i < len(weight); i++ {
		res[i] = make([]int, capacity+1)
		res[i][0] = 0
	}

	for i := 0; i < capacity+1; i++ {
		if i >= weight[0] {
			res[0][i] = value[0]
		}
	}

	for i := 1; i < len(weight); i++ {
		for j := 1; i <= capacity; j++ {
			res[i][j] = max(res[i-1][j], res[i-1][j-weight[i]]+value[i])
		}
	}

	return res[len(weight)-1][capacity]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
