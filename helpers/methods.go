package helpers

// Contains returns true if a value exists in a slice, else returns false
func Contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}
