package windows

import "math"

func maxSatisfied(customers []int, grumpy []int, minutes int) int {
	left := 0
	right := 0
	windowCount := 0
	ans := 0

	// 移动窗口，计算在 grumpy 状态下的最大可增加满意度
	for right < len(customers) {
		if grumpy[right] == 1 {
			windowCount += customers[right]
		}
		if right-left+1 > minutes {
			if grumpy[left] == 1 {
				windowCount -= customers[left]
			}
			left += 1
		}
		right += 1
		ans = int(math.Max(float64(ans), float64(windowCount)))
	}

	// 加上所有不 grumpy 状态下的满意度
	for i := 0; i < len(customers); i++ {
		if grumpy[i] == 0 {
			ans += customers[i]
		}
	}

	return ans
}
