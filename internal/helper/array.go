package helper

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
