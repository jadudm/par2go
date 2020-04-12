package verify

// Check is a cheat.
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// CheckForFile does stuff.
func CheckForFile(file string) (bool, error) {
	return true, nil
}
