package utils

// Stba String to ByteArray
func Stba(input string) []byte {
	return []byte(input)
}

func Contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}
