package windows

import "math"

func maxScore1(cardPoints []int, k int) int {
	//将长度为k的的cardPoints 拼接在 后面
	//然后就是滑动窗口找最大值

	cardPoints = append(cardPoints, cardPoints[:k]...)

	left := 0
	right := 0

	windowsSum := 0
	maxSum := math.MinInt32

	for right < len(cardPoints) {
		windowsSum += cardPoints[right]
		right++

		if right-left > k {
			maxSum = max(maxSum, windowsSum)
			windowsSum -= cardPoints[left]
			left++
		}

	}

	return maxSum

}
