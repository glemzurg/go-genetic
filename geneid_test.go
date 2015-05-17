package genetic

import (
	. "gopkg.in/check.v1" // https://labix.org/gocheck
)

// Create a suite.
type GeneIdSuite struct{}

var _ = Suite(&GeneIdSuite{})

// Add the tests.

func (s *GeneIdSuite) Test_MaxGeneId(c *C) {
	var geneId uint64

	// The max gene id should start out as 0.
	c.Assert(_maxGeneId, Equals, uint64(0))

	// Get a new gene id.
	geneId = newGeneId()
	c.Assert(_maxGeneId, Equals, uint64(1))
	c.Assert(geneId, Equals, uint64(1))

	// Get a new gene id.
	geneId = newGeneId()
	c.Assert(_maxGeneId, Equals, uint64(2))
	c.Assert(geneId, Equals, uint64(2))

	// Get a new gene id.
	geneId = newGeneId()
	c.Assert(_maxGeneId, Equals, uint64(3))
	c.Assert(geneId, Equals, uint64(3))

	// Set the gene id to a desired value.
	setMaxGeneId(100)
	c.Assert(_maxGeneId, Equals, uint64(100))

	// Get a new gene id.
	geneId = newGeneId()
	c.Assert(_maxGeneId, Equals, uint64(101))
	c.Assert(geneId, Equals, uint64(101))

	// Get a new gene id.
	geneId = newGeneId()
	c.Assert(_maxGeneId, Equals, uint64(102))
	c.Assert(geneId, Equals, uint64(102))

	// Get a new gene id.
	geneId = newGeneId()
	c.Assert(_maxGeneId, Equals, uint64(103))
	c.Assert(geneId, Equals, uint64(103))
}
