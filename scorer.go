package genetic

// Scorer is the scoring part of an experiment, implemented specifically by code that understands the material being studied.
type Scorer interface {

	// Score a particular member of the population. The cppn is what is being scored. The population is all members of
	// this generation, including the cppn being scored. cppnIndex is where the cppn is located in the population.
	//
	// The score is the score of the specimen. The bonus will be added to the score and represents any extra
	// quality valuing the CPPN (e.g. a novelty search.). The outcomes are used  with selectors that analyze
	// multiple outcomes (e.g. a hyper-volume indicator).
	Score(cppn NeatCppn, population []NeatCppn, cppnIndex int) (score float64, bonus float64, outcomes []float64)

	// The scorer may gather extra details we want to capture for each generation.
	GenerationStart(generationNum uint64)
	GenerationDetails() (json []byte)
}
