package genetic

import (
	. "gopkg.in/check.v1" // https://labix.org/gocheck
	"math/rand"
	"time"
)

// Create a suite.
type NeatCppnSuite struct{}

var _ = Suite(&NeatCppnSuite{})

// Add the tests.

func (s *NeatCppnSuite) Test_NewNeatCppn(c *C) {
	setMaxGeneId(0)

	c.Skip("This test has been verified but is unpredictable so should be manually reviewed.")

	// Get the randomness rolling.
	rand.Seed(time.Now().UnixNano())

	// Create an interface of inputs and outputs for the CPPN.
	var inOut CppnInOut = CppnInOut{
		Inputs:  []string{"i1", "i2", "i3", "i4"},
		Outputs: []string{"o1", "o2"},
	}

	// Make a new CPPN.
	var cppn NeatCppn = NewNeatCppn(inOut)

	// The contents are random. Just inspect it with a test.
	c.Assert(cppn, Equals, "unpredictable")
}

func (s *NeatCppnSuite) Test_NeatCppn_AddConnection(c *C) {
	var ok bool

	// Make a new CPPN (avoiding randomness).
	var cppn NeatCppn = NeatCppn{
		InOut: CppnInOut{
			Inputs:  []string{"i1", "i2", "i3", "i4"},
			Outputs: []string{"o1", "o2"},
		},
	}

	// Set the max gene id.
	setMaxGeneId(0)
	c.Assert(_maxGeneId, Equals, uint64(0))

	// The genome is empty at the moment.
	c.Assert(cppn.Genome, DeepEquals, NeatGenome{})

	// Add a connection.
	ok = cppn.addConnection("i1", "o1", 0.5)
	c.Assert(ok, Equals, true)
	c.Check(cppn.Genome, DeepEquals, NeatGenome{Genes: []NeatGene{
		NeatGene{GeneId: 1, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "o1", Weight: 0.5},
	}})
	c.Check(_maxGeneId, Equals, uint64(1))

	// Add a different connection.
	ok = cppn.addConnection("i2", "o2", 0.3)
	c.Assert(ok, Equals, true)
	c.Check(cppn.Genome, DeepEquals, NeatGenome{Genes: []NeatGene{
		NeatGene{GeneId: 1, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "o1", Weight: 0.5},
		NeatGene{GeneId: 2, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i2", To: "o2", Weight: 0.3},
	}})
	c.Check(_maxGeneId, Equals, uint64(2))

	// Attempt to add the same connection again.
	ok = cppn.addConnection("i2", "o2", 0.5)
	c.Assert(ok, Equals, false) // Wasn't added.
	c.Check(cppn.Genome, DeepEquals, NeatGenome{Genes: []NeatGene{
		NeatGene{GeneId: 1, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "o1", Weight: 0.5},
		NeatGene{GeneId: 2, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i2", To: "o2", Weight: 0.3},
	}})
	c.Check(_maxGeneId, Equals, uint64(2)) // GeneId not incremented.

	// Attempt to add a connection with a input as the sink.
	c.Assert(func() { cppn.addConnection("i1", "i2", 0.88888) }, Panics, `Cannot use input as sink: 'i2'`)
	c.Check(_maxGeneId, Equals, uint64(2)) // GeneId not incremented.

	// Attempt to add a connection with a input as the sink.
	c.Assert(func() { cppn.addConnection("i1", "b", 0.88888) }, Panics, `Cannot use bias as sink: 'b'`)
	c.Check(_maxGeneId, Equals, uint64(2)) // GeneId not incremented.

	// Attempt to add a connection with an output as the source.
	c.Assert(func() { cppn.addConnection("o1", "o2", 0.88888) }, Panics, `Cannot use output as source: 'o1'`)
	c.Check(_maxGeneId, Equals, uint64(2)) // GeneId not incremented.
}

func (s *NeatCppnSuite) Test_NeatCppn_AddConnection_HiddenNodes(c *C) {
	var ok bool

	// Build a cppn we can test challenging connections on.
	var cppn NeatCppn = NeatCppn{
		InOut: CppnInOut{
			Inputs:  []string{"i1"},
			Outputs: []string{"o1"},
		},
		Genome: NeatGenome{Genes: []NeatGene{
			NeatGene{GeneId: 1, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_INVERSE},
			NeatGene{GeneId: 2, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_SINE},
			NeatGene{GeneId: 3, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_RAMP},
			NeatGene{GeneId: 4, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "b", To: "o1", Weight: 0.1},
			NeatGene{GeneId: 5, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "1", Weight: 0.2},
			NeatGene{GeneId: 6, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "1", To: "2", Weight: 0.3},
			NeatGene{GeneId: 7, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "2", To: "3", Weight: 0.4},
			NeatGene{GeneId: 8, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "3", To: "o1", Weight: 0.5},
			NeatGene{GeneId: 9, IsEnabled: false, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "o1", Weight: 0.6}, // Disabled!
		},
		},
	}

	// Set the max gene id.
	setMaxGeneId(9)
	c.Assert(_maxGeneId, Equals, uint64(9))

	// Attempt to add a circular dependency
	ok = cppn.addConnection("1", "1", 0.88888)
	c.Assert(ok, Equals, false)                 // Wasn't added.
	c.Assert(len(cppn.Genome.Genes), Equals, 9) // Unchanged.
	c.Assert(_maxGeneId, Equals, uint64(9))     // GeneId not incremented.

	// Attempt to add a circular dependency
	ok = cppn.addConnection("2", "1", 0.88888)
	c.Assert(ok, Equals, false)                 // Wasn't added.
	c.Assert(len(cppn.Genome.Genes), Equals, 9) // Unchanged.
	c.Assert(_maxGeneId, Equals, uint64(9))     // GeneId not incremented.

	// Attempt to add a circular dependency
	ok = cppn.addConnection("3", "1", 0.88888)
	c.Assert(ok, Equals, false)                 // Wasn't added.
	c.Assert(len(cppn.Genome.Genes), Equals, 9) // Unchanged.
	c.Assert(_maxGeneId, Equals, uint64(9))     // GeneId not incremented.

	// Unknown nodes.
	c.Assert(func() { cppn.addConnection("unknown", "3", 0.88888) }, Panics, `Unknown from node: 'unknown'`)
	c.Assert(func() { cppn.addConnection("3", "unknown", 0.88888) }, Panics, `Unknown to node: 'unknown'`)
	c.Assert(_maxGeneId, Equals, uint64(9)) // GeneId not incremented.

	// Make some valid connection to hidden nodes.

	// Connect an input to a hidden node.
	ok = cppn.addConnection("i1", "2", 0.12)
	c.Assert(ok, Equals, true) // Wasn't added.
	c.Check(cppn.Genome.Genes[9], DeepEquals, NeatGene{GeneId: 10, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "2", Weight: 0.12})
	c.Check(_maxGeneId, Equals, uint64(10)) // Gene id incremented.

	// Connect a hidden node to a hidden node.
	ok = cppn.addConnection("1", "3", 0.13)
	c.Assert(ok, Equals, true) // Wasn't added.
	c.Check(cppn.Genome.Genes[10], DeepEquals, NeatGene{GeneId: 11, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "1", To: "3", Weight: 0.13})
	c.Check(_maxGeneId, Equals, uint64(11)) // Gene id incremented.

	// Connect a hidden node to an output.
	ok = cppn.addConnection("2", "o1", 0.14)
	c.Assert(ok, Equals, true) // Wasn't added.
	c.Check(cppn.Genome.Genes[11], DeepEquals, NeatGene{GeneId: 12, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "2", To: "o1", Weight: 0.14})
	c.Check(_maxGeneId, Equals, uint64(12)) // Gene id incremented.

	// Connect the bias to a hidden node.
	ok = cppn.addConnection("b", "1", 0.15)
	c.Assert(ok, Equals, true) // Wasn't added.
	c.Check(cppn.Genome.Genes[12], DeepEquals, NeatGene{GeneId: 13, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "b", To: "1", Weight: 0.15})
	c.Check(_maxGeneId, Equals, uint64(13)) // Gene id incremented.
}

func (s *NeatCppnSuite) Test_NeatCppn_AddNode(c *C) {
	var cppn NeatCppn

	// Make a new CPPN (avoiding randomness).
	cppn = NeatCppn{
		InOut: CppnInOut{
			Inputs:  []string{"i1", "i2", "i3", "i4"},
			Outputs: []string{"o1", "o2"},
		},
		Genome: NeatGenome{Genes: []NeatGene{
			NeatGene{GeneId: 1, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "o1", Weight: 1.1},  // Pick weight that can never be randomized to.
			NeatGene{GeneId: 2, IsEnabled: false, Type: _GENE_TYPE_CONNECTION, From: "i2", To: "o2", Weight: 1.2}, // Pick weight that can never be randomized to.
			NeatGene{GeneId: 3, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i3", To: "o1", Weight: 1.3},  // Pick weight that can never be randomized to.
			NeatGene{GeneId: 4, IsEnabled: true, Type: "SOMETHING_ELSE", From: "i4", To: "o2", Weight: 1.4},       // Pick weight that can never be randomized to.
		}},
	}

	// Add a node to the first gene.
	setMaxGeneId(345)
	c.Assert(_maxGeneId, Equals, uint64(345))
	cppn.addNode(0, ACTIVATION_SIGMOID)
	c.Check(cppn.Genome.Genes, DeepEquals, []NeatGene{
		NeatGene{GeneId: 1, IsEnabled: false, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "o1", Weight: 1.1}, // Disable the original connection.
		NeatGene{GeneId: 2, IsEnabled: false, Type: _GENE_TYPE_CONNECTION, From: "i2", To: "o2", Weight: 1.2},
		NeatGene{GeneId: 3, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i3", To: "o1", Weight: 1.3},
		NeatGene{GeneId: 4, IsEnabled: true, Type: "SOMETHING_ELSE", From: "i4", To: "o2", Weight: 1.4},
		NeatGene{GeneId: 346, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_SIGMOID},             // Node added to genome.
		NeatGene{GeneId: 347, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "346", Weight: 1.1}, // First half of original connection.
		NeatGene{GeneId: 348, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "346", To: "o1", Weight: 1.1}, // First half of original connection.
	})

	// Make a new CPPN (avoiding randomness).
	cppn = NeatCppn{
		InOut: CppnInOut{
			Inputs:  []string{"i1", "i2", "i3", "i4"},
			Outputs: []string{"o1", "o2"},
		},
		Genome: NeatGenome{Genes: []NeatGene{
			NeatGene{GeneId: 1, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "o1", Weight: 1.1},  // Pick weight that can never be randomized to.
			NeatGene{GeneId: 2, IsEnabled: false, Type: _GENE_TYPE_CONNECTION, From: "i2", To: "o2", Weight: 1.2}, // Pick weight that can never be randomized to.
			NeatGene{GeneId: 3, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i3", To: "o1", Weight: 1.3},  // Pick weight that can never be randomized to.
			NeatGene{GeneId: 4, IsEnabled: true, Type: "SOMETHING_ELSE", From: "i4", To: "o2", Weight: 1.4},       // Pick weight that can never be randomized to.
		}},
	}

	// Add a node to the second gene.
	setMaxGeneId(345)
	c.Assert(_maxGeneId, Equals, uint64(345))
	c.Check(func() { cppn.addNode(1, ACTIVATION_SIGMOID) }, Panics, `Disabled genes cannot be split with a new node.`)

	// Make a new CPPN (avoiding randomness).
	cppn = NeatCppn{
		InOut: CppnInOut{
			Inputs:  []string{"i1", "i2", "i3", "i4"},
			Outputs: []string{"o1", "o2"},
		},
		Genome: NeatGenome{Genes: []NeatGene{
			NeatGene{GeneId: 1, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "o1", Weight: 1.1},  // Pick weight that can never be randomized to.
			NeatGene{GeneId: 2, IsEnabled: false, Type: _GENE_TYPE_CONNECTION, From: "i2", To: "o2", Weight: 1.2}, // Pick weight that can never be randomized to.
			NeatGene{GeneId: 3, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i3", To: "o1", Weight: 1.3},  // Pick weight that can never be randomized to.
			NeatGene{GeneId: 4, IsEnabled: true, Type: "SOMETHING_ELSE", From: "i4", To: "o2", Weight: 1.4},       // Pick weight that can never be randomized to.
		}},
	}

	// Add a node to the third gene.
	setMaxGeneId(345)
	c.Assert(_maxGeneId, Equals, uint64(345))
	cppn.addNode(2, ACTIVATION_SIGMOID)
	c.Check(cppn.Genome.Genes, DeepEquals, []NeatGene{
		NeatGene{GeneId: 1, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "o1", Weight: 1.1},
		NeatGene{GeneId: 2, IsEnabled: false, Type: _GENE_TYPE_CONNECTION, From: "i2", To: "o2", Weight: 1.2},
		NeatGene{GeneId: 3, IsEnabled: false, Type: _GENE_TYPE_CONNECTION, From: "i3", To: "o1", Weight: 1.3}, // Disable the original connection.
		NeatGene{GeneId: 4, IsEnabled: true, Type: "SOMETHING_ELSE", From: "i4", To: "o2", Weight: 1.4},
		NeatGene{GeneId: 346, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_SIGMOID},             // Node added to genome.
		NeatGene{GeneId: 347, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i3", To: "346", Weight: 1.3}, // First half of original connection.
		NeatGene{GeneId: 348, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "346", To: "o1", Weight: 1.3}, // First half of original connection.
	})

	// Make a new CPPN (avoiding randomness).
	cppn = NeatCppn{
		InOut: CppnInOut{
			Inputs:  []string{"i1", "i2", "i3", "i4"},
			Outputs: []string{"o1", "o2"},
		},
		Genome: NeatGenome{Genes: []NeatGene{
			NeatGene{GeneId: 1, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "o1", Weight: 1.1},  // Pick weight that can never be randomized to.
			NeatGene{GeneId: 2, IsEnabled: false, Type: _GENE_TYPE_CONNECTION, From: "i2", To: "o2", Weight: 1.2}, // Pick weight that can never be randomized to.
			NeatGene{GeneId: 3, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i3", To: "o1", Weight: 1.3},  // Pick weight that can never be randomized to.
			NeatGene{GeneId: 4, IsEnabled: true, Type: "SOMETHING_ELSE", From: "i4", To: "o2", Weight: 1.4},       // Pick weight that can never be randomized to.
		}},
	}

	// Add a node to the forth gene.
	setMaxGeneId(345)
	c.Assert(_maxGeneId, Equals, uint64(345))
	c.Check(func() { cppn.addNode(3, ACTIVATION_SIGMOID) }, Panics, `Only genes of type 'connection' can have nodes added, not type: 'SOMETHING_ELSE'`)
}

func (s *NeatCppnSuite) Test_NeatCppn_MutateAddNode(c *C) {
	setMaxGeneId(0)

	c.Skip("This test has been verified but is unpredictable so should be manually reviewed.")

	// Get the randomness rolling.
	rand.Seed(time.Now().UnixNano())

	// Make a new CPPN (avoiding randomness).
	var cppn NeatCppn = NeatCppn{
		InOut: CppnInOut{
			Inputs:  []string{"i1", "i2", "i3", "i4"},
			Outputs: []string{"o1", "o2"},
		},
		Genome: NeatGenome{Genes: []NeatGene{
			NeatGene{GeneId: 1, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "o1", Weight: 1.1},  // Pick weight that can never be randomized to.
			NeatGene{GeneId: 2, IsEnabled: false, Type: _GENE_TYPE_CONNECTION, From: "i2", To: "o2", Weight: 1.2}, // Pick weight that can never be randomized to.
			NeatGene{GeneId: 3, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i3", To: "o1", Weight: 1.3},  // Pick weight that can never be randomized to.
			NeatGene{GeneId: 4, IsEnabled: true, Type: "SOMETHING_ELSE", From: "i4", To: "o2", Weight: 1.4},       // Pick weight that can never be randomized to.
		}},
	}

	// Mutate the CPPN.
	cppn.MutateAddNode([]string{ACTIVATION_BIPOLAR_SIGMOID, ACTIVATION_INVERSE, ACTIVATION_SINE})

	// The contents are random. Just inspect it with a test.
	c.Assert(cppn, Equals, "unpredictable")
}

func (s *NeatCppnSuite) Test_NeatCppn_MutateAddNode_NoActivationFunctions(c *C) {

	// Make a new CPPN (avoiding randomness).
	var cppn NeatCppn = NeatCppn{
		InOut: CppnInOut{
			Inputs:  []string{"i1", "i2", "i3", "i4"},
			Outputs: []string{"o1", "o2"},
		},
		Genome: NeatGenome{Genes: []NeatGene{
			NeatGene{GeneId: 1, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "o1", Weight: 1.1},  // Pick weight that can never be randomized to.
			NeatGene{GeneId: 2, IsEnabled: false, Type: _GENE_TYPE_CONNECTION, From: "i2", To: "o2", Weight: 1.2}, // Pick weight that can never be randomized to.
			NeatGene{GeneId: 3, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i3", To: "o1", Weight: 1.3},  // Pick weight that can never be randomized to.
			NeatGene{GeneId: 4, IsEnabled: true, Type: "SOMETHING_ELSE", From: "i4", To: "o2", Weight: 1.4},       // Pick weight that can never be randomized to.
		}},
	}

	// Invalid parameters.
	c.Check(func() { cppn.MutateAddNode(nil) }, Panics, `Available functions must be defined to mutetate add node.`)
	c.Check(func() { cppn.MutateAddNode([]string{}) }, Panics, `Available functions must be defined to mutetate add node.`)
}

func (s *NeatCppnSuite) Test_NeatCppn_MutateAddConnection(c *C) {

	c.Skip("This test has been verified but is unpredictable so should be manually reviewed.")

	// Get the randomness rolling.
	rand.Seed(time.Now().UnixNano())

	// Make a new CPPN (avoiding randomness).
	var cppn NeatCppn = NeatCppn{
		InOut: CppnInOut{
			Inputs:  []string{"i1", "i2", "i3", "i4"},
			Outputs: []string{"o1", "o2"},
		},
		Genome: NeatGenome{Genes: []NeatGene{
			NeatGene{GeneId: 1, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "o1", Weight: 1.1},  // Pick weight that can never be randomized to.
			NeatGene{GeneId: 2, IsEnabled: false, Type: _GENE_TYPE_CONNECTION, From: "i2", To: "o2", Weight: 1.2}, // Pick weight that can never be randomized to.
			NeatGene{GeneId: 3, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i3", To: "o1", Weight: 1.3},  // Pick weight that can never be randomized to.
			NeatGene{GeneId: 4, IsEnabled: true, Type: "SOMETHING_ELSE", From: "i4", To: "o2", Weight: 1.4},       // Pick weight that can never be randomized to.
			NeatGene{GeneId: 5, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_INVERSE},
			NeatGene{GeneId: 6, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_INVERSE},
		}},
	}
	setMaxGeneId(6)

	// Mutate the CPPN.
	var wasAdded bool = cppn.MutateAddConnection(1)

	// The contents are random. Just inspect it with a test.
	c.Assert(wasAdded, Equals, true)
	c.Assert(cppn, Equals, "unpredictable")
}

func (s *NeatCppnSuite) Test_NeatCppn_MutateChangeConnectionWeight(c *C) {

	c.Skip("This test has been verified but is unpredictable so should be manually reviewed.")

	// Get the randomness rolling.
	rand.Seed(time.Now().UnixNano())

	// Make a new CPPN (avoiding randomness).
	var cppn NeatCppn = NeatCppn{
		InOut: CppnInOut{
			Inputs:  []string{"i1", "i2", "i3", "i4"},
			Outputs: []string{"o1", "o2"},
		},
		Genome: NeatGenome{Genes: []NeatGene{
			NeatGene{GeneId: 1, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "o1", Weight: 1.1},  // Pick weight that can never be randomized to.
			NeatGene{GeneId: 2, IsEnabled: false, Type: _GENE_TYPE_CONNECTION, From: "i2", To: "o2", Weight: 1.2}, // Pick weight that can never be randomized to.
			NeatGene{GeneId: 3, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i3", To: "o1", Weight: 1.3},  // Pick weight that can never be randomized to.
			NeatGene{GeneId: 4, IsEnabled: true, Type: "SOMETHING_ELSE", From: "i4", To: "o2", Weight: 1.4},       // Pick weight that can never be randomized to.
			NeatGene{GeneId: 5, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_INVERSE},
			NeatGene{GeneId: 6, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_INVERSE},
		}},
	}

	// Mutate the CPPN.
	cppn.MutateChangeConnectionWeight()

	// The contents are random. Just inspect it with a test.
	c.Assert(cppn, Equals, "unpredictable")
}

func (s *NeatCppnSuite) Test_NeatCppn_MutateAddConnection_NoMaxAttempts(c *C) {

	// Make a new CPPN (avoiding randomness).
	var cppn NeatCppn = NeatCppn{
		InOut: CppnInOut{
			Inputs:  []string{"i1", "i2", "i3", "i4"},
			Outputs: []string{"o1", "o2"},
		},
		Genome: NeatGenome{Genes: []NeatGene{
			NeatGene{GeneId: 1, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "o1", Weight: 1.1},  // Pick weight that can never be randomized to.
			NeatGene{GeneId: 2, IsEnabled: false, Type: _GENE_TYPE_CONNECTION, From: "i2", To: "o2", Weight: 1.2}, // Pick weight that can never be randomized to.
			NeatGene{GeneId: 3, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i3", To: "o1", Weight: 1.3},  // Pick weight that can never be randomized to.
			NeatGene{GeneId: 4, IsEnabled: true, Type: "SOMETHING_ELSE", From: "i4", To: "o2", Weight: 1.4},       // Pick weight that can never be randomized to.
		}},
	}

	// Invalid parameters.
	c.Check(func() { cppn.MutateAddConnection(0) }, Panics, `Must have a 1 or more max attempts to mutate add connection, not: 0`)
	c.Check(func() { cppn.MutateAddConnection(-1) }, Panics, `Must have a 1 or more max attempts to mutate add connection, not: -1`)
}

func (s *NeatCppnSuite) Test_Mate(c *C) {

	c.Skip("This test has been verified but is unpredictable so should be manually reviewed.")

	// Get the randomness rolling.
	rand.Seed(time.Now().UnixNano())

	// Make a CPPN (avoiding randomness).
	var fitterCppn NeatCppn = NeatCppn{
		InOut: CppnInOut{
			Inputs:  []string{"i1", "i2", "i3", "i4"},
			Outputs: []string{"o1", "o2"},
		},
		Genome: NeatGenome{Genes: []NeatGene{
			NeatGene{GeneId: 1, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "o1", Weight: 1.1},  // Pick weight that can never be randomized to.
			NeatGene{GeneId: 2, IsEnabled: false, Type: _GENE_TYPE_CONNECTION, From: "i2", To: "o2", Weight: 1.2}, // Pick weight that can never be randomized to.
			NeatGene{GeneId: 3, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i3", To: "o1", Weight: 1.3},  // Pick weight that can never be randomized to.
			NeatGene{GeneId: 4, IsEnabled: true, Type: "SOMETHING_ELSE", From: "i4", To: "o2", Weight: 1.4},       // Pick weight that can never be randomized to.
			NeatGene{GeneId: 5, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_INVERSE},
			NeatGene{GeneId: 6, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_INVERSE},
			NeatGene{GeneId: 8, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i4", To: "o1", Weight: 1.5}, // Gene in just this CPPN.
		}},
	}

	// Make anotherCPPN (avoiding randomness).
	var otherCppn NeatCppn = NeatCppn{
		InOut: CppnInOut{
			Inputs:  []string{"i1", "i2", "i3", "i4"},
			Outputs: []string{"o1", "o2"},
		},
		Genome: NeatGenome{Genes: []NeatGene{
			NeatGene{GeneId: 1, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "o1", Weight: 2.1},  // Pick weight that can never be randomized to.
			NeatGene{GeneId: 2, IsEnabled: false, Type: _GENE_TYPE_CONNECTION, From: "i2", To: "o2", Weight: 2.2}, // Pick weight that can never be randomized to.
			NeatGene{GeneId: 3, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i3", To: "o1", Weight: 2.3},  // Pick weight that can never be randomized to.
			NeatGene{GeneId: 4, IsEnabled: true, Type: "SOMETHING_ELSE", From: "i4", To: "o2", Weight: 2.4},       // Pick weight that can never be randomized to.
			NeatGene{GeneId: 5, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_INVERSE},
			NeatGene{GeneId: 6, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_INVERSE},
			NeatGene{GeneId: 7, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_INVERSE}, // Gene in just this CPPN.
		}},
	}

	// Mutate the CPPN.
	var child NeatCppn = Mate(fitterCppn, otherCppn)

	// The contents are random. Just inspect it with a test.
	c.Assert(child, Equals, "unpredictable")
}

func (s *NeatCppnSuite) Test_Mate_UnorderedGenes(c *C) {

	// Make a CPPN with properly ordered genes.
	var orderedCppn NeatCppn = NeatCppn{
		InOut: CppnInOut{
			Inputs:  []string{"i1", "i2", "i3", "i4"},
			Outputs: []string{"o1", "o2"},
		},
		Genome: NeatGenome{Genes: []NeatGene{
			NeatGene{GeneId: 1, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "o1", Weight: 1.1},  // Pick weight that can never be randomized to.
			NeatGene{GeneId: 2, IsEnabled: false, Type: _GENE_TYPE_CONNECTION, From: "i2", To: "o2", Weight: 1.2}, // Pick weight that can never be randomized to.
		}},
	}

	// Make a CPPN with improperly ordered genes.
	var unorderedCppn NeatCppn = NeatCppn{
		InOut: CppnInOut{
			Inputs:  []string{"i1", "i2", "i3", "i4"},
			Outputs: []string{"o1", "o2"},
		},
		Genome: NeatGenome{Genes: []NeatGene{
			NeatGene{GeneId: 2, IsEnabled: false, Type: _GENE_TYPE_CONNECTION, From: "i2", To: "o2", Weight: 1.2}, // Pick weight that can never be randomized to.
			NeatGene{GeneId: 1, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "o1", Weight: 1.1},  // Pick weight that can never be randomized to.
		}},
	}

	// Invalid parameters.
	c.Check(func() { Mate(orderedCppn, unorderedCppn) }, Panics, `genome not sorted correctly by gene id: [{GeneId:2 IsEnabled:false Type:connection From:i2 To:o2 Weight:1.2 Function:} {GeneId:1 IsEnabled:true Type:connection From:i1 To:o1 Weight:1.1 Function:}]`)
	c.Check(func() { Mate(unorderedCppn, orderedCppn) }, Panics, `genome not sorted correctly by gene id: [{GeneId:2 IsEnabled:false Type:connection From:i2 To:o2 Weight:1.2 Function:} {GeneId:1 IsEnabled:true Type:connection From:i1 To:o1 Weight:1.1 Function:}]`)
}

func (s *NeatCppnSuite) Test_NeatCppn_RandomizedClone(c *C) {

	// Get the randomness rolling.
	rand.Seed(time.Now().UnixNano())

	// Make a new CPPN (avoiding randomness).
	var cppn NeatCppn = NeatCppn{
		InOut: CppnInOut{
			Inputs:  []string{"i1", "i2", "i3", "i4"},
			Outputs: []string{"o1", "o2"},
		},
		Genome: NeatGenome{Genes: []NeatGene{
			NeatGene{GeneId: 1, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "o1", Weight: 1.1},  // Pick weight that can never be randomized to.
			NeatGene{GeneId: 2, IsEnabled: false, Type: _GENE_TYPE_CONNECTION, From: "i2", To: "o2", Weight: 1.2}, // Pick weight that can never be randomized to.
			NeatGene{GeneId: 3, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i3", To: "o1", Weight: 1.3},  // Pick weight that can never be randomized to.
			NeatGene{GeneId: 4, IsEnabled: true, Type: "SOMETHING_ELSE", From: "i4", To: "o2", Weight: 1.4},       // Pick weight that can never be randomized to.
		}},
	}

	// Create a randomized clone.
	var clone NeatCppn = cppn.RandomizedClone()

	// The clone and original share the same inputs and outputs.
	c.Assert(clone.InOut, DeepEquals, cppn.InOut)
	// The clone and original share the same number of genes.
	c.Assert(len(clone.Genome.Genes), Equals, len(cppn.Genome.Genes))

	// The original cppn is unchanged.
	c.Check(cppn.Genome.Genes[0].Weight, Equals, 1.1)
	c.Check(cppn.Genome.Genes[1].Weight, Equals, 1.2)
	c.Check(cppn.Genome.Genes[2].Weight, Equals, 1.3)
	c.Check(cppn.Genome.Genes[3].Weight, Equals, 1.4)

	// For the clone, some genes have been given new weights.
	// The first gene should have a new weight.
	c.Check(clone.Genome.Genes[0].Weight, Not(Equals), 1.1)
	// The second gene should be unchanged, it was disabled.
	c.Check(clone.Genome.Genes[1].Weight, Equals, 1.2)
	// The third gene should have a new weight.
	c.Check(clone.Genome.Genes[2].Weight, Not(Equals), 1.3)
	// The forth gene should be unchanged, it was not a connection.
	c.Check(clone.Genome.Genes[3].Weight, Equals, 1.4)
}

func (s *NeatCppnSuite) Test_NeatCppn_Compute(c *C) {

	// Make a new CPPN (avoiding randomness).
	var cppn NeatCppn = NeatCppn{
		InOut: CppnInOut{
			Inputs:  []string{"i1", "i2"},
			Outputs: []string{"o1", "o2"},
		},
		Genome: NeatGenome{Genes: []NeatGene{
			NeatGene{GeneId: 1, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "o1", Weight: 0.1},
			NeatGene{GeneId: 2, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_INVERSE},
			NeatGene{GeneId: 3, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i2", To: "2", Weight: 0.25},
			NeatGene{GeneId: 4, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "2", To: "o1", Weight: 0.4},
			NeatGene{GeneId: 5, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "b", To: "2", Weight: 0.5},
			NeatGene{GeneId: 6, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "b", To: "o2", Weight: 0.5},
		}},
	}

	// Compute and get the expected outputs.
	//
	// i1 = 10.0
	// i2 = 100.0
	// b  = 1.0
	// h2 = -1.0 * (100.0 * 0.25 + 1.0 * 0.5) = -25.5
	// o1 = (10.0 * 0.1 - 25.5 * 0.4) = 1 - 10.2 = -9.2
	// o2 = (1.0 * 0.5) = 0.5
	var expectedOutputs map[string]float64 = map[string]float64{
		"o1": (10.0*0.1 + ((100.0*0.25+1.0*0.5)*-1.0)*0.4),
		"o2": (1.0 * 0.5),
	}
	var outputs map[string]float64 = cppn.Compute(map[string]float64{"i1": 10.0, "i2": 100.0})
	c.Assert(len(outputs), Equals, len(expectedOutputs))
	c.Assert(int64(outputs["o1"]*10000.0), Equals, int64(expectedOutputs["o1"]*10000.0)) // Round with typecast.
	c.Assert(int64(outputs["o2"]*10000.0), Equals, int64(expectedOutputs["o2"]*10000.0)) // Round with typecast.

	// Attempt to compute missing an input.
	c.Check(func() { cppn.Compute(map[string]float64{"i1": 10.0}) }, Panics, `Missing input: 'i2'`)

	// Attempt to compute with unknown input.
	c.Check(func() { cppn.Compute(map[string]float64{"i1": 10.0, "i2": 100.0, "i3": 100.0}) }, Panics, `Unknown input: 'i3'`)
}

func (s *NeatCppnSuite) Test_NeatCppn_PrepareComputeTopology_CircularDependency(c *C) {

	// Make a new CPPN (avoiding randomness).
	var cppn NeatCppn = NeatCppn{
		InOut: CppnInOut{
			Inputs:  []string{"i1", "i2"},
			Outputs: []string{"o1", "o2"},
		},
		Genome: NeatGenome{Genes: []NeatGene{
			NeatGene{GeneId: 1, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i1", To: "o1", Weight: 0.1},
			NeatGene{GeneId: 2, IsEnabled: true, Type: _GENE_TYPE_NODE, Function: ACTIVATION_INVERSE},
			NeatGene{GeneId: 3, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "i2", To: "2", Weight: 0.25},
			NeatGene{GeneId: 4, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "2", To: "o1", Weight: 0.4},
			NeatGene{GeneId: 5, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "b", To: "2", Weight: 0.5},
			NeatGene{GeneId: 6, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "b", To: "o2", Weight: 0.5},
			NeatGene{GeneId: 7, IsEnabled: true, Type: _GENE_TYPE_CONNECTION, From: "2", To: "2", Weight: 0.5}, // Circular dependency.
		}},
	}

	//  Attempting to compute before preparing compute topology will fail.
	c.Check(func() { cppn.PrepareComputeTopology() }, Panics, `CPPN has a circular dependency in Genome: [{GeneId:1 IsEnabled:true Type:connection From:i1 To:o1 Weight:0.1 Function:} {GeneId:2 IsEnabled:true Type:node From: To: Weight:0 Function:inverse} {GeneId:3 IsEnabled:true Type:connection From:i2 To:2 Weight:0.25 Function:} {GeneId:4 IsEnabled:true Type:connection From:2 To:o1 Weight:0.4 Function:} {GeneId:5 IsEnabled:true Type:connection From:b To:2 Weight:0.5 Function:} {GeneId:6 IsEnabled:true Type:connection From:b To:o2 Weight:0.5 Function:} {GeneId:7 IsEnabled:true Type:connection From:2 To:2 Weight:0.5 Function:}]`)
}
