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
	quickSort(test, 0, len(test)-1)
	fmt.Println(test)
}

func quickSort(nums []int, left, right int) {
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

	quickSort(nums, left, moveleft-1)
	quickSort(nums, moveleft+1, right)
}
