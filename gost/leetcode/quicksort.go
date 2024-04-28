package leetcode

func QuickSort1(nums []int, left, right int) {
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

	QuickSort1(nums, left, moveleft-1)
	QuickSort1(nums, moveleft+1, right)
}

func Quicksort(nums []int) {
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

	Quicksort(nums[left:left_index])
	Quicksort(nums[left_index+1 : right+1])
}
