package _024Lee

func twoSum(nums []int, target int) []int {
	var (
		res    []int
		numMap = make(map[int]int, len(nums))
	)

	for i, num := range nums {
		tem := target - num
		if value, ok := numMap[tem]; ok {
			res = []int{i, value}
			return res
		}
		numMap[num] = i
	}
	return res
}
