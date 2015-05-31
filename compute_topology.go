package genetic

import (
	"fmt"
	"strconv"
)

// The structure of a CPPN, ready for computation
type computeTopology struct {
	orderedNodes []string                   // The order nodes should be calculated in.
	nodes        map[string]topologicalNode // The nodes in a form that is easy to compute.
}

// topologicalNode is a node for computation that keeps track of its inputs and outputs.
type topologicalNode struct {
	nodeId     string             // What node is this?
	inputCount uint               // How many inputs is this node waiting on before it acts?
	sinks      map[string]float64 // What nodes does this node send its output to, with what weight?
	function   string             // If a hidden node, what is the function to run?
}

// addSink adds a connection.
func (n *topologicalNode) addSink(node string, weight float64) {
	// Allcate the sinks if not allocated yet.
	if n.sinks == nil {
		n.sinks = map[string]float64{}
	}
	// Sanity check we are not creating the same connection twice.
	var ok bool
	if _, ok = n.sinks[node]; ok {
		panic(fmt.Sprintf("Connection made twice from node '%s' to node: '%s'", n.nodeId, node))
	}
	// Indicate what weight the connection has.
	n.sinks[node] = weight
}

// makeComputeTopology returns the computational form of a CPPN, a data format fit for computing the output from a CPPN.
// This method is also used to determine if there is a circular dependency when randomly adding new connections.
// All input errors will panic except circular dependencies. All inputs should be sanitized by the time they reach
// this code, but the circular dependencies can be made by random mutation of connections. For circular dependencies,
// ok will be false.
func makeComputeTopology(inOut CppnInOut, genes []NeatGene) (compute computeTopology, noCircular bool) {
	var ok bool

	// Create a map of all the nodes and inputs to them.
	var nodeMap map[string]*topologicalNode = map[string]*topologicalNode{}

	// The bias.
	nodeMap[NODE_BIAS] = &topologicalNode{nodeId: NODE_BIAS}

	// The inputs.
	for _, in := range inOut.Inputs {
		nodeMap[in] = &topologicalNode{nodeId: in}
	}

	// The outputs.
	for _, out := range inOut.Outputs {
		nodeMap[out] = &topologicalNode{nodeId: out}
	}

	// The hidden nodes.
	for _, gene := range genes {
		if gene.IsEnabled == true && gene.Type == _GENE_TYPE_NODE {
			var nodeId string = strconv.FormatUint(gene.GeneId, _BASE_10)
			nodeMap[nodeId] = &topologicalNode{nodeId: nodeId, function: gene.Function}
		}
	}

	// All nodes have been added.

	// Go through the gene one at a time and construct the connection data.
	for _, gene := range genes {
		if gene.IsEnabled == true && gene.Type == _GENE_TYPE_CONNECTION {

			// Sanity check the from/to exist.
			if _, ok = nodeMap[gene.From]; !ok {
				panic(fmt.Sprintf("Unknown from node: '%s'", gene.From))
			}
			if _, ok = nodeMap[gene.To]; !ok {
				panic(fmt.Sprintf("Unknown to node: '%s'", gene.To))
			}

			// Add the connection to the source.
			nodeMap[gene.From].addSink(gene.To, gene.Weight)

			// Add the reference to the sink.
			nodeMap[gene.To].inputCount++
		}
	}

	// Verify that each output has sources feeding it a value. To compute, each output must be defined.
	for _, out := range inOut.Outputs {
		if nodeMap[out].inputCount == 0 {
			panic(fmt.Sprintf("CPPN output '%s' has no values feeding it.", out))
		}
	}

	// Now that we have everything laid out in a clear structure, determine the order of processing
	// using a depth first topological sort.

	// We need to keep track of nodes we have completed and added so we don't re-process them.
	var sunkNodes map[string]uint = map[string]uint{}

	// The starting nodes are the inputs and the bias.
	var orderedNodeIds []string = []string{NODE_BIAS}        // Start with the bias node itself.
	orderedNodeIds = append(orderedNodeIds, inOut.Inputs...) // Add the inputs.

	// Keep looping until we have examined all the nodes.
	// The length of the nodes will keep getting larger until we are done.
	for i := 0; i < len(orderedNodeIds); i++ {
		var nodeId string = orderedNodeIds[i]

		// For each node this node feeds, add it to the ordered list if it hasn't already been added.
		for sink, _ := range nodeMap[nodeId].sinks {

			// Indicate this node has a sink to it.
			sunkNodes[sink]++

			// Is this the last sink for this node?
			if sunkNodes[sink] == nodeMap[sink].inputCount {
				// This node should now be added to the ordered nodes, it is ready to calculate.
				orderedNodeIds = append(orderedNodeIds, sink)
			}
		}
	}

	// Check for circular dependencies. They way they will manifest is that the ordered nodeids wil be shorter
	// than the number of actual nodes. This is because, while processing ordered nodes, we will never reach the input count
	// for a node in the circular dependency, so it never gets added. Since it doesn't get added, the ordered node list will
	// be too short.
	if len(orderedNodeIds) != len(nodeMap) {
		return computeTopology{}, false // Circular dependency.
	}

	// Convert the pointers to topological nodes to the form we want in the compute topology.
	var nodes map[string]topologicalNode = map[string]topologicalNode{}
	for nodeId, nodePtr := range nodeMap {
		nodes[nodeId] = *nodePtr
	}

	// Return teh well-formed computeTopology
	return computeTopology{orderedNodes: orderedNodeIds, nodes: nodes}, true
}
