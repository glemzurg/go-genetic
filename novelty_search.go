package genetic

import (
	"math/rand"
	"sort"
)

// NoveltySearchTally is a structure for keeping track of how many times a result has been seen.
// Novelty searches reward new outcomes higher than previously seen outcomes.
// Internally, the fingerprint is stored twice. The fingerprint is expected to be a md5 in text format as a 32 digit hexadecimal number.
// The count stored is a 32-bit integer, so memory needed is approximately: max fingeprints * (32-char string + 32-char string + 32-bit int).
type NoveltySearchTally struct {
	maxFingerprints   int            // We can't have this structure grow for ever. At some point if must start replacing its contents.
	fingerprints      []string       // The ordered list of all fingerprints seen.
	fingerprintCounts map[string]int // How many times have we seen each result? Keyed by md5.
}

// NewNoveltySearchTally creates a new novelty search tally ready for use. When max fingerprints are reached,
// each new unseen fingerprint will randomly remove an existing fingerprint so there is space for the new one.
func NewNoveltySearchTally(maxFingerprints int) NoveltySearchTally {
	if maxFingerprints == 0 {
		panic("NovelySearchTally cannot be made with maxFingerprints == 0")
	}
	return NoveltySearchTally{
		maxFingerprints:   maxFingerprints,
		fingerprints:      []string{},
		fingerprintCounts: map[string]int{}, // How many times have we seen each result? Keyed by md5.
	}
}

// Seen indicates how many times we've seen this particular fingerprint, including this one.
func (n *NoveltySearchTally) Seen(fingerprint string) int {
	// Has this fingerprint been seen already?
	var ok bool
	if _, ok = n.fingerprintCounts[fingerprint]; ok {

		// We already have seen this fingerprint. We've now seen it one more time.
		n.fingerprintCounts[fingerprint]++

	} else {

		// If we are already full, remove an existing fingerprint.
		if len(n.fingerprintCounts) >= n.maxFingerprints {
			n.removeRandomEntry()
		}

		// We've never seen this fingerprint. But we have now.
		n.fingerprintCounts[fingerprint] = 1

		// We need to remember this fingerprint in our sorted list of current fingerprints.
		var insertionIndex int = sort.SearchStrings(n.fingerprints, fingerprint)
		// https://code.google.com/p/go-wiki/wiki/SliceTricks
		n.fingerprints = append(n.fingerprints[:insertionIndex], append([]string{fingerprint}, n.fingerprints[insertionIndex:]...)...)
	}

	return n.fingerprintCounts[fingerprint]
}

// removeRandomEntry removes a single entry at random. This is used when we reach max fingerprints and want to add another.
// The purpose is to make a comprimise between removing a map enty without taking lots of memory or processing.
func (n *NoveltySearchTally) removeRandomEntry() {

	// Pick a random fingerprint.
	var fingerprintCount int = len(n.fingerprints)
	var indexToRemove int = rand.Intn(fingerprintCount) // Assume the seed has been set.
	var fingerprintToRemove string = n.fingerprints[indexToRemove]

	// Remove the lookup value.
	// https://code.google.com/p/go-wiki/wiki/SliceTricks
	n.fingerprints = append(n.fingerprints[:indexToRemove], n.fingerprints[indexToRemove+1:]...)

	// Remove the count.
	delete(n.fingerprintCounts, fingerprintToRemove)
}
