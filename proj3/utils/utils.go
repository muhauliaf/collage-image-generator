package utils

// Sum calculates the sum of an array if int
func Sum(arr []int) int {
	total := 0
	for _, v := range arr {
		total += v
	}
	return total
}

// ArangeFloat creates an incrementing array of float64 from 0 to n-1
func ArangeFloat(n int) []float64 {
	sequence := make([]float64, n)
	for i := 0; i < n; i++ {
		sequence[i] = float64(i)
	}
	return sequence
}

// Interpolate projects values in x based on mapping from values in xp to values in fp
// reference: https://www.geeksforgeeks.org/interpolation-in-python/
func Interpolate(x []float64, xp []float64, fp []float64) []float64 {
	result := make([]float64, len(x))
	for n, xi := range x {
		i := 0
		for i < len(xp)-1 && xp[i+1] <= xi {
			i++
		}

		if i < len(xp)-1 {
			x0 := xp[i]
			x1 := xp[i+1]
			y0 := fp[i]
			y1 := fp[i+1]

			y := y0 + (y1-y0)*(xi-x0)/(x1-x0)
			result[n] = y
		} else {
			result[n] = 0.0
		}
	}
	return result
}
