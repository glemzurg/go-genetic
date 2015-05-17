package genetic

import (
	"fmt"
	"math"
)

const (
	// The activation functions used in CPPN nodes.
	// Reference: http://www.computing.dcu.ie/~humphrys/Notes/Neural/sigmoid.html
	ACTIVATION_SIGMOID = "sigmoid"
	// Reference: http://www.cs.ucf.edu/~hastings/index.php?content=ann
	ACTIVATION_BIPOLAR_SIGMOID    = "bipolar_sigmoid"
	ACTIVATION_INVERSE            = "inverse"
	ACTIVATION_SINE               = "sine"
	ACTIVATION_COSINE             = "cosine"
	ACTIVATION_TANGENT            = "tangent"
	ACTIVATION_HYPERBOLIC_TANGENT = "hyperbolic_tangent"
	ACTIVATION_GAUSSIAN           = "gaussian"
	ACTIVATION_RAMP               = "ramp"
	ACTIVATION_STEP               = "step"
	ACTIVATION_SPIKE              = "spike"

	// The math.Exp() chokes if the input goes way out of range.
	// For consistent output from activation functions, include function-specific input thresholds.
	_ACTIVATION_SIGMOID_INPUT_THRESHOLD = -100.0
)

// activate runs the given activation function on the input
func activate(function string, input float64) (output float64) {
	switch function {
	case ACTIVATION_SIGMOID:
		output = activationSigmoid(input)
	case ACTIVATION_BIPOLAR_SIGMOID:
		output = activationBipolarSigmoid(input)
	case ACTIVATION_GAUSSIAN:
		output = activationGaussian(input)
	case ACTIVATION_INVERSE:
		output = activationInverse(input)
	case ACTIVATION_SINE:
		output = activationSine(input)
	case ACTIVATION_COSINE:
		output = activationCosine(input)
	case ACTIVATION_TANGENT:
		output = activationTangent(input)
	case ACTIVATION_HYPERBOLIC_TANGENT:
		output = activationHyperbolicTangent(input)
	case ACTIVATION_RAMP:
		output = activationRamp(input)
	case ACTIVATION_STEP:
		output = activationStep(input)
	case ACTIVATION_SPIKE:
		output = activationSpike(input)
	default:
		panic(fmt.Errorf("Unknown activation function: '%s'", function))
	}
	return
}

// activationSigmoid is the sigmoid activation function. It graduall curves from 0.0 to 1.0 in an "S" shape.
// Reference: http://en.wikipedia.org/wiki/Sigmoid_function
// Reference: http://www.computing.dcu.ie/~humphrys/Notes/Neural/sigmoid.html
func activationSigmoid(input float64) (output float64) {
	// The math.Exp() doesn't handle extreme inputs well.
	// If we are about to go there just return the limit we're approaching.
	if input < _ACTIVATION_SIGMOID_INPUT_THRESHOLD {
		return 0.0
	}
	return 1.0 / (1.0 + math.Exp(-1.0*input))
}

// activationBipolarSigmoid is the bipoloer sigmoid activation function. It gradually curves from -1.0 to 1.0 in an "S" shape.
// Reference: http://www.cs.ucf.edu/~hastings/index.php?content=ann (although curve is in the wrong direction)
func activationBipolarSigmoid(input float64) (output float64) {
	// The math.Exp() doesn't handle extreme inputs well.
	// If we are about to go there just return the limit we're approaching.
	if input < _ACTIVATION_SIGMOID_INPUT_THRESHOLD {
		return -1.0
	}
	return (1.0 - math.Exp(-1.0*input)) / (1.0 + math.Exp(-1.0*input))
}

// activationGaussian is the guassian activation function. It gradually curves in a bell curve from 0.0 to 1.0 then back to 0.0.
// Reference: http://www.cs.ucf.edu/~hastings/index.php?content=ann
func activationGaussian(input float64) (output float64) {
	return math.Exp(-1.0 * (input * input))
}

// activationInverse is the inverse activation function. It flips the input sign.
// Reference: http://www.cs.ucf.edu/~hastings/index.php?content=ann
func activationInverse(input float64) (output float64) {
	return -1.0 * input
}

// activationSine is the sine activation function. It passes the value through a sine function.
// Reference: http://www.cs.ucf.edu/~hastings/index.php?content=ann
func activationSine(input float64) (output float64) {
	return math.Sin(input)
}

// activationCosine is the cosine activation function. It passes the value through a cosine function.
// Reference: http://www.cs.ucf.edu/~hastings/index.php?content=ann
func activationCosine(input float64) (output float64) {
	return math.Cos(input)
}

// activationTangent is the tangent activation function. It passes the value through a tangent function.
// Reference: http://www.cs.ucf.edu/~hastings/index.php?content=ann
func activationTangent(input float64) (output float64) {
	return math.Tan(input)
}

// activationHyperbolicTangent is the hyperbolic tangent activation function. It passes the value through a hyperbolic tangent function.
// Reference: http://www.cs.ucf.edu/~hastings/index.php?content=ann
func activationHyperbolicTangent(input float64) (output float64) {
	return math.Tanh(input)
}

// activationRamp is the ramp activation function. It produces slanting lines from 1.0 to -1.0
// Reference: http://www.cs.ucf.edu/~hastings/index.php?content=ann
func activationRamp(input float64) (output float64) {
	return 1.0 - 2.0*(input-math.Floor(input))
}

// activationStep is the step activation function. It produces alternating values of 1.0 and -1.0 for each integer, creating a square ridge pattern.
// Reference: http://www.cs.ucf.edu/~hastings/index.php?content=ann
func activationStep(input float64) (output float64) {
	// Alternate between 1.0 and -1.0 for each integer.
	var floor float64 = math.Floor(input)
	if math.Mod(floor, 2.0) == 0.0 {
		return 1.0
	}
	return -1.0
}

// activationSpike is the spike activation function. It produces point sine wave zigzaging between -1.0 and 1.0.
// Reference: http://www.cs.ucf.edu/~hastings/index.php?content=ann
func activationSpike(input float64) (output float64) {
	// Alternate between 1.0 and -1.0 for each integer.
	var absFloor float64 = math.Abs(math.Floor(input))
	if math.Mod(absFloor, 2.0) == 0.0 {
		return 1.0 - 2.0*(input-math.Floor(input))
	}
	return -1.0 + 2.0*(input-math.Floor(input))
}
