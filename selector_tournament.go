package genetic

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"sort"
)

// SelectorTournament runs "competitions" among random members of the population. The fittest in each random sampling
// is kept for the next generation. This is down until we have the number we want to keep. If the number of contenders
// is the size of the population, this ends up being an Elitism selector (but with more effort).
type SelectorTournament struct {
	KeepCount  int // How many specimens should be kept in each generation?
	Contenders int // How many specimens "compete" in each competition?
}

// Select runs small competitions to pick the fittest of the population.
func (s *SelectorTournament) Select(specimens []Specimen) (fittest []Specimen) {

	// If we are keeping the entire population, just pass it through.
	if s.KeepCount >= len(specimens) {
		return specimens
	}

	// Don't assume any order to the specimens.

	// Build a list of specimens we intend to keep for the next generation.
	// The list is just the indexes of the specimen slice.
	var keepers []int

	// Look at the population of contenders as just indexes into the specimen slice.
	var population []int
	for i := 0; i < len(specimens); i++ {
		population = append(population, i) // Naturally sorted ascending.
	}

	// Until we have enough keepers, the the games continue... "are you not entertained!"
	for len(keepers) < s.KeepCount {

		// First get a single contender. It is currently the best.
		var best int = randomLookupWithSkip(population, nil)

		// Keep track of the contenders so they aren't picked again.
		var contenders []int = []int{best}

		// Now, compete our current best with others.
		var otherContenderCount int = s.Contenders - 1 // Don't include the one we already have.
		for i := 0; i < otherContenderCount; i++ {

			// Who is the next contender?
			var contender int = randomLookupWithSkip(population, contenders)

			// Fight!
			if specimens[contender].SelectionScore > specimens[best].SelectionScore {
				best = contender
			}

			// Add the contender to our list and keep it sorted.
			contenders = append(contenders, contender)
			sort.Ints(contenders)
		}

		// The competition is over. The best is a specimen we want to keep.
		keepers = append(keepers, best)

		// Remove it from the population. It's not going to fight anymore.
		// Preserve the ascending order of the population.
		var foundIndex int = sort.SearchInts(population, best)
		if foundIndex < len(population) && population[foundIndex] == best {
			// We already have the best stored away. Just remove it from the slice.
			population = append(population[:foundIndex], population[foundIndex+1:]...)
		} else {
			log.Panic("ASSERT: Searched population for a non-existent value.")
		}
	}

	// Now turn the keepers into specimens.
	fittest = nil
	for _, keep := range keepers {
		fittest = append(fittest, specimens[keep])
	}
	return fittest
}

// randomLookupWithSkip picks a random lookup from the slice, skipping already grabbed lookups. The skip indexes are
// expected to be sorted.
func randomLookupWithSkip(lookups []int, skipIndexes []int) int {
	// The length of the pickable slice is less the number of elements already grabbed.
	var pickIndex int = rand.Intn(len(lookups) - len(skipIndexes))
	// Any of the skip indexes that are lower than this index pushes this index up by one.
	for _, skipIndex := range skipIndexes {
		if pickIndex >= skipIndex {
			// A skipped index is earlier than this pick, push it up one.
			pickIndex++
		} else {
			// The skipped indexes are sorted.
			// We just found a skip index higher than the pick index.
			// No more will push the pick index around.
			break
		}
	}
	// Get that lookup and return it.
	return lookups[pickIndex]
}

// LoadSelectorTournamentConfig loads the json filename as a new configuration.
func LoadSelectorTournamentConfig(filename string) (SelectorTournament, error) {
	var err error
	var bytes []byte
	var selector SelectorTournament

	log.Printf("Loading tournament selector Config: '%s'\n", filename)

	// Load and parse from json.
	if bytes, err = ioutil.ReadFile(filename); err != nil {
		return SelectorTournament{}, err
	}
	if err = json.Unmarshal(bytes, &selector); err != nil {
		return SelectorTournament{}, err
	}
	selector.validOrPanic()
	return selector, error(nil)
}

// validOrPanic panics if we're not ready for use.
func (s *SelectorTournament) validOrPanic() {
	if s.KeepCount < 1 {
		log.Panicf("KeepCount must be one or more: %d", s.KeepCount)
	}
	if s.Contenders < 2 {
		log.Panicf("Contenders must be two or more: %d", s.Contenders)
	}
}
