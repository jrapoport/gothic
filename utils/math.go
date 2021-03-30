package utils

// Max returns the max of x & y
func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// Clamp clamps int x to min and max
func Clamp(x, min, max int) int {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}
