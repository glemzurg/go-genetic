package genetic

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
)

// NeatCppn is NEAT CPPN: NeuroEvolution of Augmenting Topologies CPPN, a nerual net that builds its own structure through
// mating and mutation. NEAT CPPNs tend to develop minimal internal connections to do the work they need. It is also not necessary
// to attempt to structure their insides.
type NeatCppn struct {
	InOut    CppnInOut
	Genome   NeatGenome
	topology computeTopology
}

// NewNeatCppn creates a new well-formed CPPN for the given inputs/outputs. All outputs must be able to produce a value when
// the CPPN is run so one random connction to an inptu will be made for each output. The new next innovation number (used to
// identify genes across the experiment)
func NewNeatCppn(inOut CppnInOut) (cppn NeatCppn) {

	// Start a new CPPN.
	cppn = NeatCppn{
		InOut: inOut,
	}

	// Connect every output to one of the input values.
	for _, out := range cppn.InOut.Outputs {

		// Pick a random input.
		var ok bool
		var inputIndex int = rand.Intn(len(cppn.InOut.Inputs))
		var in string = cppn.InOut.Inputs[inputIndex]

		// Pick a random weight.
		var weight float64
		weight = rand.Float64() // Actually will never by 1.0 but will be less than 1.0.

		// Make the connection. Should always work.
		if ok = cppn.addConnection(in, out, weight); !ok {
			panic(fmt.Sprintf("Failed to create connection from '%s' to '%s' when creating a new NeatCppn", in, out))
		}
	}

	// We have a minimal well-formed CPPN.
	return cppn
}

// addConnection creates a new connection in the genome. Indicates if the connect was added. It will not be added if the
// connection woudl be invalid because it duplicates an existing connection or creates a circular dependency.
func (c *NeatCppn) addConnection(from string, to string, weight float64) (wasAdded bool) {

	// Verify we can add this gene.

	// We cannot add a connection from a gene to itself. The ultimate circular dependency.
	if from == to {
		return false
	}

	// An enabled version of the gene must not already exist.
	for _, gene := range c.Genome.Genes {

		// We cannot add the same connection again.
		if gene.IsEnabled && gene.From == from && gene.To == to {
			return false
		}

		// We cannot add the same connection but in reverse either.
		// This is the cheapest circular dependency to find.
		if gene.IsEnabled && gene.From == to && gene.To == from {
			return false
		}
	}

	// Connections cannot be made to bias.
	if to == NODE_BIAS {
		panic(fmt.Sprintf("Cannot use bias as sink: '%s'", to))
	}

	// Connections cannot be made to inputs.
	if inStrings(c.InOut.Inputs, to) {
		panic(fmt.Sprintf("Cannot use input as sink: '%s'", to))
	}

	// Connections cannot be made from outputs.
	if inStrings(c.InOut.Outputs, from) {
		panic(fmt.Sprintf("Cannot use output as source: '%s'", from))
	}

	// Start with a gene with no id (we don't want to use an id until we know this gene is valid).
	var newGene NeatGene = NeatGene{
		IsEnabled: true,
		Type:      _GENE_TYPE_CONNECTION,
		From:      from,
		To:        to,
		Weight:    weight,
	}

	// We don't need to test for circular dependencies of the from is an input or bias, and to is an output.
	var isFromInput bool
	if from == NODE_BIAS || inStrings(c.InOut.Inputs, from) {
		isFromInput = true
	}
	var isToOutput bool
	if inStrings(c.InOut.Outputs, to) {
		isToOutput = true
	}

	// Test the new gene (before getting a real geneId), verify it does not create any circular dependencies.
	// As long as we are not a simple wiring from input to output.
	if !(isFromInput && isToOutput) {
		var testGenes []NeatGene = make([]NeatGene, len(c.Genome.Genes))
		copy(testGenes, c.Genome.Genes)
		testGenes = append(testGenes, newGene) // The fake node can be tested without giving it a node id.
		var ok bool
		if _, ok = makeComputeTopology(c.InOut, testGenes); !ok {
			// This gene made a circular dependency.
			return false
		}
	}

	// Its all good. Get a new gene id.
	newGene.GeneId = newGeneId()
	c.Genome.Genes = append(c.Genome.Genes, newGene)
	return true
}

// addNode adds a node to the NEAT CPPN. The node is always a hidden node appearing on an existing connection,
// splitting it into two connections. One connection goes from the original source node to the hidden node.
// The other goes from the hidden node to the original destination node.
func (c *NeatCppn) addNode(connectionGeneIndex int, function string) {

	// Only enabled genes can be split..
	if !c.Genome.Genes[connectionGeneIndex].IsEnabled {
		panic("Disabled genes cannot be split with a new node.")
	}

	// Only connection genees can be split.
	if c.Genome.Genes[connectionGeneIndex].Type != _GENE_TYPE_CONNECTION {
		panic(fmt.Sprintf("Only genes of type '%s' can have nodes added, not type: '%s'", _GENE_TYPE_CONNECTION, c.Genome.Genes[connectionGeneIndex].Type))
	}

	// Disable the original node.
	c.Genome.Genes[connectionGeneIndex].IsEnabled = false

	// What is the from, to, weight of the original node?
	var from string = c.Genome.Genes[connectionGeneIndex].From
	var to string = c.Genome.Genes[connectionGeneIndex].To
	var weight float64 = c.Genome.Genes[connectionGeneIndex].Weight

	// Create a new new node named after its gene id.
	var nodeGeneId uint64 = newGeneId()
	var nodeId string = strconv.FormatUint(nodeGeneId, _BASE_10)

	// Add the node.
	c.Genome.Genes = append(c.Genome.Genes, NeatGene{GeneId: nodeGeneId, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: function})

	// Add new connections that take the place of the disabled connectino but have the node in the middle.
	c.Genome.Genes = append(c.Genome.Genes, NeatGene{GeneId: newGeneId(), IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: from, To: nodeId, Weight: weight})
	c.Genome.Genes = append(c.Genome.Genes, NeatGene{GeneId: newGeneId(), IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: nodeId, To: to, Weight: weight})
}

// MutateAddNode adds a new node to the CPPN with a randomly picked activation function and spliting a randomly selected
// existing connection, putting the node in the middle of it.
func (c *NeatCppn) MutateAddNode(availableFunctions []string) {

	// If we can't pick an activation function, we can't create a new node.
	if len(availableFunctions) == 0 {
		panic("Available functions must be defined to mutetate add node.")
	}

	// Randomly pick one activiation function.
	var functionIndex int = rand.Intn(len(availableFunctions))
	var function string = availableFunctions[functionIndex]

	// Randomly pick an enabled connection.
	// First count and remember enabled connections.
	var enabledConnectionIndexes []int
	for i, gene := range c.Genome.Genes {
		if gene.IsEnabled == true && gene.Type == _GENE_TYPE_CONNECTION {
			enabledConnectionIndexes = append(enabledConnectionIndexes, i)
		}
	}
	// Pick one of those indexes.
	var pickedIndex int = rand.Intn(len(enabledConnectionIndexes))
	var geneIndex int = enabledConnectionIndexes[pickedIndex]

	// Add the node.
	c.addNode(geneIndex, function)
}

// MutateAddConnection adds a new valid connection to the CPPN randomly wiring two nodes together.
// It's possible that it randomly attempts to make a connection that is invalid (creating a circular depenency).
// It will try up to max attempts to keep making connections, and indicate if one was made.
func (c *NeatCppn) MutateAddConnection(maxAttempts int) (wasAdded bool) {

	// If we don't know how long we can go, report an issue.
	if maxAttempts < 1 {
		panic(fmt.Sprintf("Must have a 1 or more max attempts to mutate add connection, not: %d", maxAttempts))
	}

	// What are all the hidden nodes?
	var hiddenNodes []string
	for _, gene := range c.Genome.Genes {
		if gene.IsEnabled == true && gene.Type == _GENE_TYPE_NODE {
			var nodeId string = strconv.FormatUint(gene.GeneId, _BASE_10)
			hiddenNodes = append(hiddenNodes, nodeId)
		}
	}

	// What are all the nodes we can make a connection from?
	var fromNodes []string = []string{NODE_BIAS}     // Start with the bias node itself.
	fromNodes = append(fromNodes, c.InOut.Inputs...) // Add the inputs.
	fromNodes = append(fromNodes, hiddenNodes...)    // Add the hidden nodes.

	// What are all the nodes we can make a connection to?
	var toNodes []string
	toNodes = append(toNodes, c.InOut.Outputs...) // Add the outputs.
	toNodes = append(toNodes, hiddenNodes...)     // Add the hidden nodes.

	// Attempt to create a connection until successful or as long as we can.
	for i := 0; i < maxAttempts; i++ {

		// Pick a random from node.
		var fromIndex int = rand.Intn(len(fromNodes))
		var from string = fromNodes[fromIndex]

		// Pick a random to node.
		var toIndex int = rand.Intn(len(toNodes))
		var to string = toNodes[toIndex]

		// Pick a random weight between 0.0 and 1.0.
		var weight float64
		weight = rand.Float64() // Actually will never by 1.0 but will be less than 1.0.

		// Make the connection. If it works we've done what we need to in this function.
		if wasAdded = c.addConnection(from, to, weight); wasAdded {
			return wasAdded // Success!
		}
	}

	// We were not succesful.
	return false
}

// MutateChangeConnectionWeight changes the weight of a randomly selected enabled connection to a random value between 0.0 and 1.0.
func (c *NeatCppn) MutateChangeConnectionWeight() {

	// Randomly pick an enabled connection.
	// First count and remember enabled connections.
	var enabledConnectionIndexes []int
	for i, gene := range c.Genome.Genes {
		if gene.IsEnabled == true && gene.Type == _GENE_TYPE_CONNECTION {
			enabledConnectionIndexes = append(enabledConnectionIndexes, i)
		}
	}
	// Pick one of those indexes.
	var pickedIndex int = rand.Intn(len(enabledConnectionIndexes))
	var geneIndex int = enabledConnectionIndexes[pickedIndex]

	// Change just the weight of this gene.
	c.Genome.Genes[geneIndex].Weight = rand.Float64() // 0.0 to 1.0. Actually will never by 1.0 but will be less than 1.0.
}

// Mate mates two CPPNs to create a new offspring. The structure of the child's genome is the genome of the fitter parent
// (so the hidden nodes in the child will be the hidden nodes of the fitter parent). For every enabled connection gene
// shared between the two parents, the connection weight will be randomly picked from one or the other. If no genes
// are modified (a possibility), the child will be identical to the fitter parent.
func Mate(fitterParent NeatCppn, otherParent NeatCppn) (child NeatCppn) {
	// Start the child from the parent.
	child = NeatCppn{
		InOut:  fitterParent.InOut, // in/out is fixed for an experiment so not a problem if it gets cross referenced in anyway.
		Genome: NeatGenome{},
	}

	// Get the genomes we are working with.
	var fitterGenes []NeatGene = fitterParent.Genome.Genes
	var otherGenes []NeatGene = otherParent.Genome.Genes

	// The genomes should be always sorted ascending by gene id but do a sanity check.
	if !sort.IsSorted(ByGeneId(fitterGenes)) {
		panic(fmt.Sprintf("genome not sorted correctly by gene id: %+v", fitterGenes))
	}
	if !sort.IsSorted(ByGeneId(otherGenes)) {
		panic(fmt.Sprintf("genome not sorted correctly by gene id: %+v", otherGenes))
	}

	// The gene structure and hidden nodes come from the fitter parent. Create them in the child now.
	for _, gene := range fitterGenes {

		// This goes to the child. The only question is if it uses the weight from the other parent.
		// Is this an enabled connection (a gene with a weight)?
		if gene.IsEnabled == true && gene.Type == _GENE_TYPE_CONNECTION {

			// Does the less-fit parent also have this gene?
			var foundIndex int = sort.Search(len(otherGenes), func(i int) bool { return otherGenes[i].GeneId >= gene.GeneId })
			if foundIndex < len(otherGenes) && otherGenes[foundIndex].GeneId == gene.GeneId {

				// Both parents have this gene. There is a 50% change we'll keep the fitter weight on this gene, and a 50%
				// chance we'll use the weight from the less-fit parent's gene. Pick either 0 or 1.
				var coinFlip int = rand.Intn(2) // Pick either 0 or 1.
				if coinFlip == 0 {
					// Use the less-fit gene's weight.
					gene.Weight = otherGenes[foundIndex].Weight
				}
			}
		}

		// Add this gene to the child.
		child.Genome.Genes = append(child.Genome.Genes, gene)
	}

	return child
}

// Clone creates a clone of the CPPN, identical but no shared data.
func (c *NeatCppn) Clone() (clone NeatCppn) {
	clone = NeatCppn{
		InOut:  c.InOut,          // in/out is fixed for an experiment so not a problem if it gets cross referenced in anyway.
		Genome: c.Genome.Clone(), // No shared gene data. Copied instead.
	}
	return clone
}

// RandomizedClone creates a clone of the CPPN and randomizes the clone's connection weights without altering the
// the structure. Once an initial CPPN structure is created for an experiment, create new members of the population by
// by cloning the template CPPN. The new weights will be between 0.0 and 1.0.
func (c *NeatCppn) RandomizedClone() (clone NeatCppn) {
	clone = NeatCppn{
		InOut:  c.InOut, // in/out is fixed for an experiment so not a problem if it gets cross referenced in anyway.
		Genome: NeatGenome{},
	}
	// The genomes need to be copied/modified one at a time and referentially distinct between the CPPNs.
	for _, origGene := range c.Genome.Genes {
		var cloneGene NeatGene = origGene
		if cloneGene.IsEnabled == true && cloneGene.Type == _GENE_TYPE_CONNECTION {
			cloneGene.Weight = rand.Float64() // 0.0 to 1.0. Actually will never by 1.0 but will be less than 1.0.
		}
		clone.Genome.Genes = append(clone.Genome.Genes, cloneGene)
	}
	return clone
}

// PrepareComputeTopology prepares a CPPN to be computed, building internal datastructures for the task.
func (c *NeatCppn) PrepareComputeTopology() {
	var ok bool
	if c.topology, ok = makeComputeTopology(c.InOut, c.Genome.Genes); !ok {
		// Somehow we ended up with a circular dependency in our genome.
		// Should never happen.
		panic(fmt.Sprintf("CPPN has a circular dependency in Genome: %+v", c.Genome.Genes))
	}
}

// Compute takes all the inputs and passes them through the CPPN to get the outputs.
func (c *NeatCppn) Compute(inputs map[string]float64) (outputs map[string]float64) {
	var ok bool

	// Have we created a topology yet?
	if c.topology.orderedNodes == nil {
		c.PrepareComputeTopology()
	}

	// Keep track of the current node values.
	var nodeValues map[string]float64 = map[string]float64{}

	// Add a sanity double-check to ensure we are not making any mistakes.
	var sinkTally map[string]uint = map[string]uint{}

	// Start by putting in all the input values.
	for _, in := range c.InOut.Inputs {
		// Did we pass in this input?
		var value float64
		if value, ok = inputs[in]; !ok {
			panic(fmt.Sprintf("Missing input: '%s'", in))
		}
		nodeValues[in] = value
		sinkTally[in] = 0 // Inputs will never have other nodes use them as a sink.
	}

	// All inputs passed in should now have node values.
	// Sanity check we didn't pass in any invalid inputs.
	for in, _ := range inputs {
		if _, ok = nodeValues[in]; !ok {
			panic(fmt.Sprintf("Unknown input: '%s'", in))
		}
	}

	// Add the bias. It always has a value of 1.0.
	nodeValues[NODE_BIAS] = 1.0
	sinkTally[NODE_BIAS] = 0 // The bias will never have other nodes use it as a sink.

	// Now just start processing the nodes one at a time.
	for _, nodeId := range c.topology.orderedNodes {

		// Get the topological details of this node.
		var node topologicalNode = c.topology.nodes[nodeId]

		// Verify we have the correct number of sinks to this node.
		if sinkTally[nodeId] != node.inputCount {
			panic(fmt.Sprintf("Incorrect sink count: expected %d but found %d", node.inputCount, sinkTally[nodeId]))
		}

		// This node has the value of all source nodes using it as a sink.
		// Or initial value for inputs and bias.
		var value float64 = nodeValues[nodeId]

		// If this node has a function, run the function on the value to get the value it will pass on.
		if node.function != "" {
			value = activate(node.function, value)
		}

		// Send this node to each sink it has, applying the weight of the connection.
		for sink, weight := range node.sinks {

			// What value are we sending to the sink?
			var weightedValue float64 = value * weight

			// Add this value to the value building on that node.
			if _, ok = nodeValues[sink]; !ok {
				nodeValues[sink] = 0.0
			}
			nodeValues[sink] += weightedValue

			// Increment out tally to this sink.
			if _, ok = sinkTally[sink]; !ok {
				sinkTally[sink] = 0
			}
			sinkTally[sink]++
		}
	}

	// Extract the output values and retun them.
	outputs = map[string]float64{}
	for _, out := range c.InOut.Outputs {
		// Did we calcualte in this output? All should be calculated.
		var value float64
		if value, ok = nodeValues[out]; !ok {
			panic(fmt.Sprintf("Failed to calculate output: '%s'", out))
		}
		outputs[out] = value
	}
	return outputs
}
