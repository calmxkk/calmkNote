package main

import (
	"fmt"
	"math/rand"
)

func main() {
	test := []int{1, 3, 4, 9, 7, 8, 100, 2}
	for i := 0; i < 10; i++ {
		test = append(test, rand.Intn(100))
	}
	fmt.Println(test)
	quickSort1(test, 0, len(test)-1)
	fmt.Println(test)
}

func quickSort1(nums []int, left, right int) {
	if left > right {
		return
	}

	moveleft, moveright := left, right
	base := nums[left]

	for moveleft < moveright {
		for moveleft < moveright && nums[moveright] > base {
			moveright--
		}

		for moveleft < moveright && nums[moveleft] <= base {
			moveleft++
		}

		if moveleft < moveright {
			nums[moveright], nums[moveleft] = nums[moveleft], nums[moveright]
		}
	}
	nums[moveleft], nums[left] = base, nums[moveleft]

	quickSort1(nums, left, moveleft-1)
	quickSort1(nums, moveleft+1, right)
}

func quicksort(nums []int) {
	left, right := 0, len(nums)-1
	if left >= right {
		return
	}

	left_index, right_index := left, right

	base := nums[left]

	for left_index < right_index {
		for left_index < right_index && nums[right_index] > base {
			right_index--
		}
		nums[left_index] = nums[right_index]

		for left_index < right_index && nums[left_index] <= base {
			left_index++
		}
		nums[right_index] = nums[left_index]
	}
	nums[left_index] = base

	quicksort(nums[left:left_index])
	quicksort(nums[left_index+1 : right+1])
}
