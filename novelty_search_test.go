package genetic

import (
	crypto_rand "crypto/rand"
	"encoding/base64"
	. "gopkg.in/check.v1" // https://labix.org/gocheck
	"math/rand"
	"sort"
	"time"
)

// Create a suite.
type NoveltySearchSuite struct{}

var _ = Suite(&NoveltySearchSuite{})

// Add the tests.

func (s *NoveltySearchSuite) Test_NoveltySearchTally(c *C) {

	// Get the randomness rolling.
	rand.Seed(time.Now().UnixNano())

	// Report how many times we've seen a fingerprint.
	var seen int

	// To ensure the internals stay insycn double check the keys of the lookup structure.
	var fingerprintKeys []string

	// Creat a new tally with a small maximum remembered keys.
	var maxFingerprintsToRemember int = 3
	var tally NoveltySearchTally = NewNoveltySearchTally(maxFingerprintsToRemember)

	// Say we've seen something.
	seen = tally.Seen("FINGERPRINT_A")
	c.Assert(seen, Equals, 1)
	c.Assert(tally.fingerprintCounts, DeepEquals, map[string]int{"FINGERPRINT_A": 1})
	fingerprintKeys = sortedMapKeys(tally.fingerprintCounts)
	c.Assert(fingerprintKeys, DeepEquals, tally.fingerprints) // The internal data structure must be intact.

	// Say we've seen something again.
	seen = tally.Seen("FINGERPRINT_A")
	c.Assert(seen, Equals, 2)
	c.Assert(tally.fingerprintCounts, DeepEquals, map[string]int{"FINGERPRINT_A": 2}) // A second viewing.
	fingerprintKeys = sortedMapKeys(tally.fingerprintCounts)
	c.Assert(fingerprintKeys, DeepEquals, tally.fingerprints) // The internal data structure must be intact.

	// Say we've seen something else.
	seen = tally.Seen("FINGERPRINT_B")
	c.Assert(seen, Equals, 1)
	c.Assert(tally.fingerprintCounts, DeepEquals, map[string]int{"FINGERPRINT_A": 2, "FINGERPRINT_B": 1})
	fingerprintKeys = sortedMapKeys(tally.fingerprintCounts)
	c.Assert(fingerprintKeys, DeepEquals, tally.fingerprints) // The internal data structure must be intact.

	// Say we've seen something else.
	seen = tally.Seen("FINGERPRINT_C")
	c.Assert(seen, Equals, 1)
	c.Assert(tally.fingerprintCounts, DeepEquals, map[string]int{"FINGERPRINT_A": 2, "FINGERPRINT_B": 1, "FINGERPRINT_C": 1})
	fingerprintKeys = sortedMapKeys(tally.fingerprintCounts)
	c.Assert(fingerprintKeys, DeepEquals, tally.fingerprints) // The internal data structure must be intact.

	// We're now at max finger prints to remember.

	// Visiting an existing fingerprint will not cause anything to be removed.
	seen = tally.Seen("FINGERPRINT_C")
	c.Assert(seen, Equals, 2)
	c.Assert(tally.fingerprintCounts, DeepEquals, map[string]int{"FINGERPRINT_A": 2, "FINGERPRINT_B": 1, "FINGERPRINT_C": 2})
	fingerprintKeys = sortedMapKeys(tally.fingerprintCounts)
	c.Assert(fingerprintKeys, DeepEquals, tally.fingerprints) // The internal data structure must be intact.

	// Seeing a new fingerprint will now remove a random existing fingerprint from the tally. Could be any of them.
	seen = tally.Seen("FINGERPRINT_D")
	c.Assert(seen, Equals, 1)
	c.Assert(len(tally.fingerprints), Equals, 3)
	fingerprintKeys = sortedMapKeys(tally.fingerprintCounts)
	c.Assert(fingerprintKeys, DeepEquals, tally.fingerprints) // The internal data structure must be intact.

	// Re-seeing a new fingerprint will not change the other fingerprints.
	seen = tally.Seen("FINGERPRINT_D")
	c.Assert(seen, Equals, 2)
	c.Assert(len(tally.fingerprints), Equals, 3)
	fingerprintKeys = sortedMapKeys(tally.fingerprintCounts)
	c.Assert(fingerprintKeys, DeepEquals, tally.fingerprints) // The internal data structure must be intact.

	// Stress the code a bit by spamming a bunch of randome insertions and deletions.
	for i := 0; i < 30; i++ {
		var fingerprint string = randomFingerprint()
		_ = tally.Seen(fingerprint)
		c.Assert(len(tally.fingerprints), Equals, 3)
		fingerprintKeys = sortedMapKeys(tally.fingerprintCounts)
		c.Assert(fingerprintKeys, DeepEquals, tally.fingerprints) // The internal data structure must be intact.
	}

	// Invalid parameters.
	c.Assert(func() { NewNoveltySearchTally(0) }, PanicMatches, `NovelySearchTally cannot be made with maxFingerprints == 0`)

}

// sortedMapKeys gets keys of a map in an order suitable for testing.
func sortedMapKeys(theMap map[string]int) []string {
	var keys []string = []string{}
	for key, _ := range theMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// randomFingerprint creates a random fingerprint for testing purposes.
func randomFingerprint() string {
	var err error
	var bytes []byte = make([]byte, 32)
	if _, err = crypto_rand.Read(bytes); err != nil {
		panic(err)
	}
	var fingerprint string = base64.URLEncoding.EncodeToString(bytes)
	return fingerprint
}
