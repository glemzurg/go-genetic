package genetic

import (
	. "gopkg.in/check.v1" // https://labix.org/gocheck
	"math"
)

// Create a suite.
type ActivationFunctionSuite struct{}

var _ = Suite(&ActivationFunctionSuite{})

// Add the tests.

func (s *ActivationFunctionSuite) Test_Activate(c *C) {

	// Sigmoid curves from 0.0 to 1.0 in a curving "S"-like slope.
	c.Check(activate(ACTIVATION_SIGMOID, -100000000000.0), Equals, 0.0)
	c.Check(activate(ACTIVATION_SIGMOID, -101.0), Equals, 0.0) // Just beyond the input threhold.
	c.Check(activate(ACTIVATION_SIGMOID, 0.0), Equals, 0.5)
	c.Check(activate(ACTIVATION_SIGMOID, 100.0), Equals, 1.0)
	c.Check(activate(ACTIVATION_SIGMOID, 100000000000.0), Equals, 1.0)

	// Bipolar sigmoid curves from -1.0 to 1.0 in a curving "S"-like slope.
	c.Check(activate(ACTIVATION_BIPOLAR_SIGMOID, -100000000000.0), Equals, -1.0)
	c.Check(activate(ACTIVATION_BIPOLAR_SIGMOID, -101.0), Equals, -1.0) // Just beyond the input threhold.
	c.Check(activate(ACTIVATION_BIPOLAR_SIGMOID, -100.0), Equals, -1.0)
	c.Check(activate(ACTIVATION_BIPOLAR_SIGMOID, 0.0), Equals, 0.0)
	c.Check(activate(ACTIVATION_BIPOLAR_SIGMOID, 100.0), Equals, 1.0)
	c.Check(activate(ACTIVATION_BIPOLAR_SIGMOID, 100000000000.0), Equals, 1.0)

	// Guassian curves in a bell curve from 0.0 to 1.0 then back to 0.0.
	c.Check(activate(ACTIVATION_GAUSSIAN, -100000000000.0), Equals, 0.0)
	c.Check(activate(ACTIVATION_GAUSSIAN, -100.0), Equals, 0.0)
	c.Check(activate(ACTIVATION_GAUSSIAN, 0.0), Equals, 1.0)
	c.Check(activate(ACTIVATION_GAUSSIAN, 100.0), Equals, 0.0)
	c.Check(activate(ACTIVATION_GAUSSIAN, 100000000000.0), Equals, 0.0)

	// Inverse flips the valeu of the input.
	c.Check(activate(ACTIVATION_INVERSE, -100.0), Equals, 100.0)
	c.Check(activate(ACTIVATION_INVERSE, 0.0), Equals, 0.0)
	c.Check(activate(ACTIVATION_INVERSE, 100.0), Equals, -100.0)

	// Sine passes the value through a sine function.
	c.Check(activate(ACTIVATION_SINE, -100000000000.0), Equals, math.Sin(-100000000000.0))
	c.Check(activate(ACTIVATION_SINE, -100.0), Equals, math.Sin(-100.0))
	c.Check(activate(ACTIVATION_SINE, 0.0), Equals, math.Sin(0.0))
	c.Check(activate(ACTIVATION_SINE, 100.0), Equals, math.Sin(100.0))
	c.Check(activate(ACTIVATION_SINE, 100000000000.0), Equals, math.Sin(100000000000.0))

	// Cosine passes the value through a cosine function.
	c.Check(activate(ACTIVATION_COSINE, -100000000000.0), Equals, math.Cos(-100000000000.0))
	c.Check(activate(ACTIVATION_COSINE, -100.0), Equals, math.Cos(-100.0))
	c.Check(activate(ACTIVATION_COSINE, 0.0), Equals, math.Cos(0.0))
	c.Check(activate(ACTIVATION_COSINE, 100.0), Equals, math.Cos(100.0))
	c.Check(activate(ACTIVATION_COSINE, 100000000000.0), Equals, math.Cos(100000000000.0))

	// Tangent passes the value through a tangent function.
	c.Check(activate(ACTIVATION_TANGENT, -100000000000.0), Equals, math.Tan(-100000000000.0))
	c.Check(activate(ACTIVATION_TANGENT, -100.0), Equals, math.Tan(-100.0))
	c.Check(activate(ACTIVATION_TANGENT, 0.0), Equals, math.Tan(0.0))
	c.Check(activate(ACTIVATION_TANGENT, 100.0), Equals, math.Tan(100.0))
	c.Check(activate(ACTIVATION_TANGENT, 100000000000.0), Equals, math.Tan(100000000000.0))

	// Hyperbolic tangent passes the value through a hyperbolic tangent function.
	c.Check(activate(ACTIVATION_HYPERBOLIC_TANGENT, -100000000000.0), Equals, math.Tanh(-100000000000.0))
	c.Check(activate(ACTIVATION_HYPERBOLIC_TANGENT, -100.0), Equals, math.Tanh(-100.0))
	c.Check(activate(ACTIVATION_HYPERBOLIC_TANGENT, 0.0), Equals, math.Tanh(0.0))
	c.Check(activate(ACTIVATION_HYPERBOLIC_TANGENT, 100.0), Equals, math.Tanh(100.0))
	c.Check(activate(ACTIVATION_HYPERBOLIC_TANGENT, 100000000000.0), Equals, math.Tanh(100000000000.0))

	// Ramp produces diagonal slanting lines, going from 1.0 to -1.0.
	c.Check(activate(ACTIVATION_RAMP, -100000000000.5), Equals, 0.0)
	c.Check(activate(ACTIVATION_RAMP, -1.0), Equals, 1.0)
	c.Check(activate(ACTIVATION_RAMP, -0.75), Equals, 0.5)
	c.Check(activate(ACTIVATION_RAMP, -0.5), Equals, 0.0)
	c.Check(activate(ACTIVATION_RAMP, -0.25), Equals, -0.5)
	c.Check(activate(ACTIVATION_RAMP, -0.1), Equals, -0.8)
	c.Check(activate(ACTIVATION_RAMP, 0.0), Equals, 1.0)
	c.Check(activate(ACTIVATION_RAMP, 0.25), Equals, 0.5)
	c.Check(activate(ACTIVATION_RAMP, 0.5), Equals, 0.0)
	c.Check(activate(ACTIVATION_RAMP, 0.75), Equals, -0.5)
	c.Check(activate(ACTIVATION_RAMP, 100000000000.5), Equals, 0.0)

	// Step returns alternating 1.0 and -1.0 creating a square ridge pattern.
	c.Check(activate(ACTIVATION_STEP, -100000000000.5), Equals, -1.0)
	c.Check(activate(ACTIVATION_STEP, -1.25), Equals, 1.0)
	c.Check(activate(ACTIVATION_STEP, -1.0), Equals, -1.0)
	c.Check(activate(ACTIVATION_STEP, -0.75), Equals, -1.0)
	c.Check(activate(ACTIVATION_STEP, -0.5), Equals, -1.0)
	c.Check(activate(ACTIVATION_STEP, -0.25), Equals, -1.0)
	c.Check(activate(ACTIVATION_STEP, -0.1), Equals, -1.0)
	c.Check(activate(ACTIVATION_STEP, 0.0), Equals, 1.0)
	c.Check(activate(ACTIVATION_STEP, 0.25), Equals, 1.0)
	c.Check(activate(ACTIVATION_STEP, 0.5), Equals, 1.0)
	c.Check(activate(ACTIVATION_STEP, 0.75), Equals, 1.0)
	c.Check(activate(ACTIVATION_STEP, 0.9), Equals, 1.0)
	c.Check(activate(ACTIVATION_STEP, 1.0), Equals, -1.0)
	c.Check(activate(ACTIVATION_STEP, 100000000000.5), Equals, 1.0)

	// Spike returns a zig-zag between 1.0 and -1.0, a angular sine-like wave.
	c.Check(activate(ACTIVATION_SPIKE, -100000000000.5), Equals, 0.0)
	c.Check(activate(ACTIVATION_SPIKE, -1.25), Equals, -0.5)
	c.Check(activate(ACTIVATION_SPIKE, -1.0), Equals, -1.0)
	c.Check(activate(ACTIVATION_SPIKE, -0.75), Equals, -0.5)
	c.Check(activate(ACTIVATION_SPIKE, -0.5), Equals, 0.0)
	c.Check(activate(ACTIVATION_SPIKE, -0.25), Equals, 0.5)
	c.Check(activate(ACTIVATION_SPIKE, -0.1), Equals, 0.8)
	c.Check(activate(ACTIVATION_SPIKE, 0.0), Equals, 1.0)
	c.Check(activate(ACTIVATION_SPIKE, 0.25), Equals, 0.5)
	c.Check(activate(ACTIVATION_SPIKE, 0.5), Equals, 0.0)
	c.Check(activate(ACTIVATION_SPIKE, 0.75), Equals, -0.5)
	c.Check(activate(ACTIVATION_SPIKE, 0.9), Equals, -0.8)
	c.Check(activate(ACTIVATION_SPIKE, 1.0), Equals, -1.0)
	c.Check(activate(ACTIVATION_SPIKE, 1.25), Equals, -0.5)
	c.Check(activate(ACTIVATION_SPIKE, 100000000000.5), Equals, 0.0)

	// Invalid parameters.
	c.Assert(func() { activate("BOOGA", 0.0) }, Panics, `Unknown activation function: 'BOOGA'`)
}
