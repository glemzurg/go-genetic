package genetic

// The geneiIds are the "innovation numbers" used in NEAT neural nets. If two genes have the same gene id, they are the same gene
// and should mean the two genes were created at the same time.
// Just keep the value in a static value. Only one call should ever be attempting to read this value.
var _maxGeneId uint64

// setMaxGeneId is used to prepare the gene id counter for a new experiment. When an experiment starts,
// it may set the max gene id to ensure that new evolutions don't get confused with neural net genes that already exist.
func setMaxGeneId(geneId uint64) {
	_maxGeneId = geneId
}

// newGeneId gets a geneid unique in the experiment.
func newGeneId() (geneId uint64) {
	// Increment and return the new gene id.
	_maxGeneId++
	return _maxGeneId
}
