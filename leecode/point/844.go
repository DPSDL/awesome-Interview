package point

// 844. 比较含退格的字符串
func backspaceCompare(s string, t string) bool {
	finalStr := func(str string) string {
		var stack []rune

		for _, char := range str {
			if char != '#' {
				stack = append(stack, char)
			} else if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		}
		return string(stack)
	}

	return finalStr(s) == finalStr(t)

}
