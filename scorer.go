package genetic

// Scorer is the scoring part of an experiment, implemented specifically by code that understands the material being studied.
type Scorer interface {

	// Score a particular member of the population. The neural net is what is being scored. The population is all members of
	// this generation, including the neural net being scored. neuralNetIndex is where the neural net is located in the population.
	//
	// The score is the score of the specimen. The bonus will be added to the score and represents any extra
	// quality valuing the neural net (e.g. a novelty search.). The outcomes are used  with selectors that analyze
	// multiple outcomes (e.g. a hyper-volume indicator).
	Score(neuralNet NeatNeuralNet, population []NeatNeuralNet, neuralNetIndex int) (score float64, bonus float64, outcomes []float64)

	// The scorer may gather extra details we want to capture for each generation.
	GenerationStart(generationNum uint64)
	GenerationDetails() (json []byte)
}
