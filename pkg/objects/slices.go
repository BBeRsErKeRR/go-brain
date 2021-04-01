package objects

func Intersection[T comparable](slices ...[]T) []T {
	counts := map[T]int{}
	result := []T{}

	for _, slice := range slices {
		for _, val := range slice {
			counts[val]++
		}
	}

	for val, count := range counts {
		if count == len(slices) {
			result = append(result, val)
		}
	}

	return result
}

func EqSlices[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func FilterSlice[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}
