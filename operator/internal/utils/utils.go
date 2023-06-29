package utils

func Contains(arr []string, val string) bool {
	for _, item := range arr {
		if item == val {
			return true
		}
	}
	return false
}

func ContainsOne(target []string, arr []string) bool {
	for _, item := range arr {
		if Contains(target, item) {
			return true
		}
	}
	return false
}
