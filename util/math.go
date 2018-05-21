package util

func Abs(x int16) int16 {
	if x < 0 {
		return -x
	}
	return x
}