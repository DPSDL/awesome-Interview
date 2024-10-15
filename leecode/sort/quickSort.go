package sort

import "math/rand"

func sortArray(nums []int) []int {
	quickSort(nums, 0, len(nums)-1)

	return nums
}

func quickSort(nums []int, left, right int) {
	if left >= right {
		return
	}

	pivot := pivotFunc(nums, left, right)

	quickSort(nums, left, pivot-1)
	quickSort(nums, pivot+1, right)

	return
}

func pivotFunc(nums []int, left, right int) int {
	pivotIndex := rand.Intn(right-left+1) + left

	nums[pivotIndex], nums[right] = nums[right], nums[pivotIndex]
	pivot := nums[right]
	j := left

	for i := left; i < right; i++ {
		if nums[i] < pivot {
			nums[j], nums[i] = nums[i], nums[j]
			j++
		}
	}

	nums[right], nums[j] = nums[j], nums[right]
	return j

}
