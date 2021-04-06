package common

func Mod(a, b int) int {
	var res int = a % b
	if (res < 0 && b > 0) || (res > 0 && b < 0) {
		return res + b
	}
	return res
}
