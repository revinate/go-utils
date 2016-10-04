package helper

func InArray(value string, arr []string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}
