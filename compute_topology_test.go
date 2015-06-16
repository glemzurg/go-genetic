package genetic

import (
	. "gopkg.in/check.v1" // https://labix.org/gocheck
)

// Create a suite.
type ComputeTopologySuite struct{}

var _ = Suite(&ComputeTopologySuite{})

// Add the tests.

func (s *ComputeTopologySuite) Test_MakeComputeTopology_Minimal(c *C) {
	var inOut NeuralNetInOut
	var genes []neatGene
	var compute computeTopology
	var ok bool
	var expected computeTopology

	// A simple compute topology.
	inOut = NeuralNetInOut{
		Inputs:  []string{"i1"},
		Outputs: []string{"o1"},
	}
	genes = []neatGene{
		neatGene{GeneId: 1, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "o1", Weight: 0.1},
		neatGene{GeneId: 2, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "b", To: "o1", Weight: 0.2},
	}
	expected = computeTopology{
		orderedNodes: []string{NODE_BIAS, "i1", "o1"},
		nodes: map[string]topologicalNode{
			NODE_BIAS: topologicalNode{
				nodeId:     NODE_BIAS,
				inputCount: 0,
				sinks: map[string]float64{
					"o1": 0.2,
				},
			},
			"i1": topologicalNode{
				nodeId:     "i1",
				inputCount: 0,
				sinks: map[string]float64{
					"o1": 0.1,
				},
			},
			"o1": topologicalNode{
				nodeId:     "o1",
				inputCount: 2,
			},
		},
	}
	compute, ok = makeComputeTopology(inOut, genes)
	c.Assert(ok, Equals, true)
	c.Assert(compute, DeepEquals, expected)
}

func (s *ComputeTopologySuite) Test_MakeComputeTopology_SimpleHiddenNodes(c *C) {
	var inOut NeuralNetInOut
	var genes []neatGene
	var compute computeTopology
	var ok bool
	var expected computeTopology

	// A simple compute topology.
	inOut = NeuralNetInOut{
		Inputs:  []string{"i1"},
		Outputs: []string{"o1"},
	}
	genes = []neatGene{
		neatGene{GeneId: 1, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_INVERSE},
		neatGene{GeneId: 2, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_SINE},
		neatGene{GeneId: 3, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_RAMP},
		neatGene{GeneId: 4, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "b", To: "o1", Weight: 0.1},
		neatGene{GeneId: 5, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "1", Weight: 0.2},
		neatGene{GeneId: 6, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "1", To: "2", Weight: 0.3},
		neatGene{GeneId: 7, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "2", To: "3", Weight: 0.4},
		neatGene{GeneId: 8, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "3", To: "o1", Weight: 0.5},
	}
	expected = computeTopology{
		orderedNodes: []string{NODE_BIAS, "i1", "1", "2", "3", "o1"},
		nodes: map[string]topologicalNode{
			NODE_BIAS: topologicalNode{
				nodeId:     NODE_BIAS,
				inputCount: 0,
				sinks: map[string]float64{
					"o1": 0.1,
				},
			},
			"i1": topologicalNode{
				nodeId:     "i1",
				inputCount: 0,
				sinks: map[string]float64{
					"1": 0.2,
				},
			},
			"1": topologicalNode{
				nodeId:     "1",
				inputCount: 1,
				sinks: map[string]float64{
					"2": 0.3,
				},
				function: ACTIVATION_INVERSE,
			},
			"2": topologicalNode{
				nodeId:     "2",
				inputCount: 1,
				sinks: map[string]float64{
					"3": 0.4,
				},
				function: ACTIVATION_SINE,
			},
			"3": topologicalNode{
				nodeId:     "3",
				inputCount: 1,
				sinks: map[string]float64{
					"o1": 0.5,
				},
				function: ACTIVATION_RAMP,
			},
			"o1": topologicalNode{
				nodeId:     "o1",
				inputCount: 2,
			},
		},
	}
	compute, ok = makeComputeTopology(inOut, genes)
	c.Assert(ok, Equals, true)
	c.Assert(compute, DeepEquals, expected)
}

func (s *ComputeTopologySuite) Test_MakeComputeTopology_ConnectionMadeTwice(c *C) {
	var inOut NeuralNetInOut
	var genes []neatGene

	// A simple compute topology.
	inOut = NeuralNetInOut{
		Inputs:  []string{"i1"},
		Outputs: []string{"o1"},
	}
	genes = []neatGene{
		neatGene{GeneId: 1, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_INVERSE},
		neatGene{GeneId: 2, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_SINE},
		neatGene{GeneId: 3, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_RAMP},
		neatGene{GeneId: 4, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "b", To: "o1", Weight: 0.1},
		neatGene{GeneId: 5, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "1", Weight: 0.2},
		neatGene{GeneId: 6, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "1", To: "2", Weight: 0.3},
		neatGene{GeneId: 7, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "2", To: "3", Weight: 0.4},
		neatGene{GeneId: 8, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "3", To: "o1", Weight: 0.5},
		neatGene{GeneId: 9, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "1", To: "2", Weight: 0.6}, // Duplicate connection.
	}

	c.Assert(func() { makeComputeTopology(inOut, genes) }, Panics, `Connection made twice from node '1' to node: '2'`)
}

func (s *ComputeTopologySuite) Test_MakeComputeTopology_UnknownSourceNode(c *C) {
	var inOut NeuralNetInOut
	var genes []neatGene

	// A simple compute topology.
	inOut = NeuralNetInOut{
		Inputs:  []string{"i1"},
		Outputs: []string{"o1"},
	}
	genes = []neatGene{
		neatGene{GeneId: 1, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_INVERSE},
		neatGene{GeneId: 2, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_SINE},
		neatGene{GeneId: 3, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_RAMP},
		neatGene{GeneId: 4, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "b", To: "o1", Weight: 0.1},
		neatGene{GeneId: 5, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "1", Weight: 0.2},
		neatGene{GeneId: 6, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "1", To: "2", Weight: 0.3},
		neatGene{GeneId: 7, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "2", To: "3", Weight: 0.4},
		neatGene{GeneId: 8, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "3", To: "o1", Weight: 0.5},
		neatGene{GeneId: 9, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "unknown", To: "2", Weight: 0.6}, // Unknown source node.
	}

	c.Assert(func() { makeComputeTopology(inOut, genes) }, Panics, `Unknown from node: 'unknown'`)
}

func (s *ComputeTopologySuite) Test_MakeComputeTopology_UnknownSinkNode(c *C) {
	var inOut NeuralNetInOut
	var genes []neatGene

	// A simple compute topology.
	inOut = NeuralNetInOut{
		Inputs:  []string{"i1"},
		Outputs: []string{"o1"},
	}
	genes = []neatGene{
		neatGene{GeneId: 1, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_INVERSE},
		neatGene{GeneId: 2, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_SINE},
		neatGene{GeneId: 3, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_RAMP},
		neatGene{GeneId: 4, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "b", To: "o1", Weight: 0.1},
		neatGene{GeneId: 5, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "1", Weight: 0.2},
		neatGene{GeneId: 6, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "1", To: "2", Weight: 0.3},
		neatGene{GeneId: 7, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "2", To: "3", Weight: 0.4},
		neatGene{GeneId: 8, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "3", To: "o1", Weight: 0.5},
		neatGene{GeneId: 9, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "2", To: "unknown", Weight: 0.6}, // Unknown source node.
	}

	c.Assert(func() { makeComputeTopology(inOut, genes) }, Panics, `Unknown to node: 'unknown'`)
}

func (s *ComputeTopologySuite) Test_MakeComputeTopology_UnfedOutput(c *C) {
	var inOut NeuralNetInOut
	var genes []neatGene

	// A simple compute topology.
	inOut = NeuralNetInOut{
		Inputs:  []string{"i1"},
		Outputs: []string{"o1", "o2"}, // The new output has no connections to it.
	}
	genes = []neatGene{
		neatGene{GeneId: 1, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_INVERSE},
		neatGene{GeneId: 2, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_SINE},
		neatGene{GeneId: 3, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_RAMP},
		neatGene{GeneId: 4, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "b", To: "o1", Weight: 0.1},
		neatGene{GeneId: 5, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "1", Weight: 0.2},
		neatGene{GeneId: 6, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "1", To: "2", Weight: 0.3},
		neatGene{GeneId: 7, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "2", To: "3", Weight: 0.4},
		neatGene{GeneId: 8, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "3", To: "o1", Weight: 0.5},
	}

	c.Assert(func() { makeComputeTopology(inOut, genes) }, Panics, `ANN output 'o2' has no values feeding it.`)
}

func (s *ComputeTopologySuite) Test_MakeComputeTopology_CircularDependencyA(c *C) {
	var inOut NeuralNetInOut
	var genes []neatGene
	var compute computeTopology
	var ok bool

	// A simple compute topology.
	inOut = NeuralNetInOut{
		Inputs:  []string{"i1"},
		Outputs: []string{"o1"},
	}
	genes = []neatGene{
		neatGene{GeneId: 1, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_INVERSE},
		neatGene{GeneId: 2, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_SINE},
		neatGene{GeneId: 3, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_RAMP},
		neatGene{GeneId: 4, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "b", To: "o1", Weight: 0.1},
		neatGene{GeneId: 5, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "1", Weight: 0.2},
		neatGene{GeneId: 6, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "1", To: "2", Weight: 0.3},
		neatGene{GeneId: 7, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "2", To: "3", Weight: 0.4},
		neatGene{GeneId: 8, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "3", To: "o1", Weight: 0.5},
		neatGene{GeneId: 9, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "1", To: "1", Weight: 0.6}, // Create circular dependency.
	}

	compute, ok = makeComputeTopology(inOut, genes)
	c.Assert(ok, Equals, false) // Circular dependency
	c.Assert(compute, DeepEquals, computeTopology{})
}

func (s *ComputeTopologySuite) Test_MakeComputeTopology_CircularDependencyB(c *C) {
	var inOut NeuralNetInOut
	var genes []neatGene
	var compute computeTopology
	var ok bool

	// A simple compute topology.
	inOut = NeuralNetInOut{
		Inputs:  []string{"i1"},
		Outputs: []string{"o1"},
	}
	genes = []neatGene{
		neatGene{GeneId: 1, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_INVERSE},
		neatGene{GeneId: 2, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_SINE},
		neatGene{GeneId: 3, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_RAMP},
		neatGene{GeneId: 4, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "b", To: "o1", Weight: 0.1},
		neatGene{GeneId: 5, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "1", Weight: 0.2},
		neatGene{GeneId: 6, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "1", To: "2", Weight: 0.3},
		neatGene{GeneId: 7, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "2", To: "3", Weight: 0.4},
		neatGene{GeneId: 8, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "3", To: "o1", Weight: 0.5},
		neatGene{GeneId: 9, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "2", To: "1", Weight: 0.6}, // Create circular dependency.
	}

	compute, ok = makeComputeTopology(inOut, genes)
	c.Assert(ok, Equals, false) // Circular dependency
	c.Assert(compute, DeepEquals, computeTopology{})
}

func (s *ComputeTopologySuite) Test_MakeComputeTopology_CircularDependencyC(c *C) {
	var inOut NeuralNetInOut
	var genes []neatGene
	var compute computeTopology
	var ok bool

	// A simple compute topology.
	inOut = NeuralNetInOut{
		Inputs:  []string{"i1"},
		Outputs: []string{"o1"},
	}
	genes = []neatGene{
		neatGene{GeneId: 1, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_INVERSE},
		neatGene{GeneId: 2, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_SINE},
		neatGene{GeneId: 3, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_RAMP},
		neatGene{GeneId: 4, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "b", To: "o1", Weight: 0.1},
		neatGene{GeneId: 5, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "1", Weight: 0.2},
		neatGene{GeneId: 6, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "1", To: "2", Weight: 0.3},
		neatGene{GeneId: 7, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "2", To: "3", Weight: 0.4},
		neatGene{GeneId: 8, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "3", To: "o1", Weight: 0.5},
		neatGene{GeneId: 9, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "3", To: "1", Weight: 0.6}, // Create circular dependency.
	}

	compute, ok = makeComputeTopology(inOut, genes)
	c.Assert(ok, Equals, false) // Circular dependency
	c.Assert(compute, DeepEquals, computeTopology{})
}
