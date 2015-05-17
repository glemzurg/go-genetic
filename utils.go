package genetic

import (
	"sort"
)

// inStrings hides the slighly confusing code for searching a slice of strings for the presence of a string.
// The slice of strings must be sorted.
func inStrings(haystack []string, needle string) (wasFound bool) {
	// The strings being searched must be sorted.
	var insertionIndex int = sort.SearchStrings(haystack, needle)
	if insertionIndex < len(haystack) && haystack[insertionIndex] == needle {
		// We found the needle.
		return true
	}
	return false
}
