package genetic

import (
	"fmt"
	"sort"
	"strconv"
)

const (
	// The name of the bias node.
	NODE_BIAS = "b"

	// String to integer parsing parameters.
	_BASE_10     = 10
	_BIT_SIZE_64 = 64
)

// NeuralNetInOut is the inputs and outputs for a neural net. All neural nets in a single experiment must share the same
// inputs and outputs.
type NeuralNetInOut struct {
	// Public methods will be captured in the configuration details of an experiment.
	Inputs  []string
	Outputs []string
	// There is always an assumed bias called "b"
}

// Validate confirms that the neural net in/out is well-formed for running the experiment.
func (i *NeuralNetInOut) validate() {

	// Verify the parameters.
	if len(i.Inputs) == 0 {
		panic("NeuralNetInOut has no inputs.")
	}
	if len(i.Outputs) == 0 {
		panic("NeuralNetInOut has no outputs.")
	}

	// Sort the inputs and outputs for easy searching.
	sort.Strings(i.Inputs)
	sort.Strings(i.Outputs)

	// Inputs and output cannot share any names.
	for _, in := range i.Inputs {
		if inStrings(i.Outputs, in) {
			panic(fmt.Sprintf("NeuralNetInOut has both input and output named '%s'", in))
		}
	}

	// Neither inputs or outputs can be named "b", which is the name of the bias.
	if inStrings(i.Inputs, NODE_BIAS) {
		panic(fmt.Sprintf("NeuralNetInOut has input named same as the bias '%s'", NODE_BIAS))
	}
	if inStrings(i.Outputs, NODE_BIAS) {
		panic(fmt.Sprintf("NeuralNetInOut has output named same as the bias '%s'", NODE_BIAS))
	}

	// None of the inputs cannot be named as numbers. Those are reserved for hidden nodes named for their geneId.
	for _, in := range i.Inputs {
		var err error
		_, err = strconv.ParseUint(in, _BASE_10, _BIT_SIZE_64)
		// There is only a problem if we did NOT get an error. We successfully parsed the string to an integer which is bad.
		if err == nil {
			panic(fmt.Sprintf("NeuralNetInOut has input named as a number '%s'. Used for hidden nodes.", in))
		}
	}

	// None of the output cannot be named as numbers. Those are reserved for hidden nodes named for their geneId.
	for _, out := range i.Outputs {
		var err error
		_, err = strconv.ParseUint(out, _BASE_10, _BIT_SIZE_64)
		// There is only a problem if we did NOT get an error. We successfully parsed the string to an integer which is bad.
		if err == nil {
			panic(fmt.Sprintf("NeuralNetInOut has output named as a number '%s'. Used for hidden nodes.", out))
		}
	}
}
