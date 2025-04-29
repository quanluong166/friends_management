package utils

func FindCommon(arr1, arr2 []string) []string {
	lookup := make(map[string]bool)
	for _, val := range arr1 {
		lookup[val] = true
	}

	var common []string
	for _, val := range arr2 {
		if lookup[val] {
			common = append(common, val)
		}
	}

	return common
}

func Contains(arr []string, str string) (bool, string) {
	for _, v := range arr {
		if v == str {
			return true, v
		}
	}
	return false, ""
}

func RemoveSameElementsFromSecond(arr1, arr2 []string) []string {
	lookup := make(map[string]bool)
	for _, val := range arr1 {
		lookup[val] = true
	}

	var result []string
	for _, val := range arr2 {
		if !lookup[val] {
			result = append(result, val)
		}
	}

	return result
}

// Append all the input array into one array
func Combine(arr ...[]string) []string {
	totalSize := len(arr)
	result := make([]string, 0, totalSize)

	// Iterate over each array in the input
	for _, each := range arr {
		result = append(result, each...)
	}

	return result
}
