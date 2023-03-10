package data

// CountDivision - this function needs for date printing
func CountDivision(x uint64) uint64 {
	num := uint64(256)
	for x/num > 256*256-1 {
		num *= 256
	}
	if num > 256 {
		return (x / 256) % (256)
	}
	return x / uint64(num)
}

// CountMod - this function needs for date printing
func CountMod(x uint64) uint64 {
	num := uint64(256)
	for x/num > 256*256-1 {
		num *= 256
	}
	if num > 256 {
		return 256 * 256
	}
	return x % uint64(num)
}

// NumberLen - Length of the number
func NumberLen(x int) int {
	cnt := 0
	for x != 0 {
		cnt++
		x /= 10
	}
	return cnt
}

// Max - return maximum
func Max(x *int, y int) {
	if y > *x {
		*x = y
	}
}
