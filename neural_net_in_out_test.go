package genetic

import (
	. "gopkg.in/check.v1" // https://labix.org/gocheck
)

// Create a suite.
type NeuralNetInOutSuite struct{}

var _ = Suite(&NeuralNetInOutSuite{})

// Add the tests.

func (s *NeuralNetInOutSuite) Test_NeuralNetInOut_Validate(c *C) {
	var inOut NeuralNetInOut

	// We always expect the well-formated in-out to be sorted.
	var expectedInOut NeuralNetInOut = NeuralNetInOut{
		Inputs:  []string{"i1", "i2", "i3"},
		Outputs: []string{"o1", "o2", "o3"},
	}

	// First, a well-formed in/out will be ok with no changes.
	inOut = NeuralNetInOut{
		Inputs:  []string{"i1", "i2", "i3"},
		Outputs: []string{"o1", "o2", "o3"},
	}
	inOut.validate() // No panic.
	c.Assert(inOut, DeepEquals, expectedInOut)

	// An unsorted in/out
	inOut = NeuralNetInOut{
		Inputs:  []string{"i1", "i3", "i2"},
		Outputs: []string{"o3", "o2", "o1"},
	}
	inOut.validate() // No panic.
	c.Assert(inOut, DeepEquals, expectedInOut)

	// A name collision.
	inOut = NeuralNetInOut{
		Inputs:  []string{"i1", "o2", "i3"},
		Outputs: []string{"o1", "o2", "o3"},
	}
	c.Assert(func() { inOut.validate() }, Panics, `NeuralNetInOut has both input and output named 'o2'`)

	// A bias collision.
	inOut = NeuralNetInOut{
		Inputs:  []string{"i1", "b", "i3"},
		Outputs: []string{"o1", "o2", "o3"},
	}
	c.Assert(func() { inOut.validate() }, Panics, `NeuralNetInOut has input named same as the bias 'b'`)

	// A bias collision.
	inOut = NeuralNetInOut{
		Inputs:  []string{"i1", "i2", "i3"},
		Outputs: []string{"o1", "b", "o3"},
	}
	c.Assert(func() { inOut.validate() }, Panics, `NeuralNetInOut has output named same as the bias 'b'`)

	// A numeric name.
	inOut = NeuralNetInOut{
		Inputs:  []string{"i1", "2", "i3"},
		Outputs: []string{"o1", "o2", "o3"},
	}
	c.Assert(func() { inOut.validate() }, Panics, `NeuralNetInOut has input named as a number '2'. Used for hidden nodes.`)

	// A numeric name.
	inOut = NeuralNetInOut{
		Inputs:  []string{"i1", "i2", "i3"},
		Outputs: []string{"o1", "2", "o3"},
	}
	c.Assert(func() { inOut.validate() }, Panics, `NeuralNetInOut has output named as a number '2'. Used for hidden nodes.`)

	// No inputs.
	inOut = NeuralNetInOut{
		Inputs:  nil,
		Outputs: []string{"o1", "o2", "o3"},
	}
	c.Assert(func() { inOut.validate() }, Panics, `NeuralNetInOut has no inputs.`)

	// No inputs.
	inOut = NeuralNetInOut{
		Inputs:  []string{},
		Outputs: []string{"o1", "o2", "o3"},
	}
	c.Assert(func() { inOut.validate() }, Panics, `NeuralNetInOut has no inputs.`)

	// No outputs.
	inOut = NeuralNetInOut{
		Inputs:  []string{"i1", "i2", "i3"},
		Outputs: nil,
	}
	c.Assert(func() { inOut.validate() }, Panics, `NeuralNetInOut has no outputs.`)

	// No outputs.
	inOut = NeuralNetInOut{
		Inputs:  []string{"i1", "i2", "i3"},
		Outputs: []string{},
	}
	c.Assert(func() { inOut.validate() }, Panics, `NeuralNetInOut has no outputs.`)

}
