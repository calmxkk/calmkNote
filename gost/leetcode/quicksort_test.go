package leetcode

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestQuicksort(t *testing.T) {
	nums := []int{}
	nums1 := []int{}
	for i := 0; i < 20; i++ {
		temp := rand.Intn(1000)
		nums = append(nums, temp)
		nums1 = append(nums1, temp)
	}

	fmt.Println(nums1)

	Quicksort(nums)

	fmt.Println(nums)

	QuickSort1(nums1, 0, len(nums1)-1)
	fmt.Println(nums1)
}
