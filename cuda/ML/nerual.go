package main

// run on the local terminal:
// go mod tidy

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

// Sigmoid activation function and its derivative
func sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

func sigmoidDerivative(x float64) float64 {
	return x * (1 - x)
}

// Neuron represents a neuron with its weights and bias
type Neuron struct {
	weights []float64
	bias    float64
}

// Layer represents a neural network layer
type Layer struct {
	neurons []Neuron
	output  []float64
}

// NeuralNetwork represents a basic neural network
type NeuralNetwork struct {
	layers []Layer
}

// Feedforward propagates the input through the network
func (nn *NeuralNetwork) Feedforward(input []float64) []float64 {
	for i := range nn.layers {
		output := make([]float64, len(nn.layers[i].neurons))
		for j := range nn.layers[i].neurons {
			sum := nn.layers[i].neurons[j].bias
			for k, weight := range nn.layers[i].neurons[j].weights {
				sum += weight * input[k]
			}
			output[j] = sigmoid(sum)
		}
		input = output // The output of this layer is the input to the next layer
		nn.layers[i].output = output
	}
	return input // The final output is the output of the last layer
}

// Train the neural network using backpropagation
func (nn *NeuralNetwork) Train(input, target []float64, learningRate float64) {
	// Forward pass
	output := nn.Feedforward(input)

	// Calculate the output layer error
	outputLayerError := make([]float64, len(nn.layers[len(nn.layers)-1].neurons))
	for i := range outputLayerError {
		outputLayerError[i] = (target[i] - output[i]) * sigmoidDerivative(output[i])
	}

	// Backward pass
	for i := len(nn.layers) - 1; i >= 0; i-- {
		layer := nn.layers[i]
		errors := make([]float64, len(layer.neurons))
		if i != len(nn.layers)-1 {
			// If not the output layer, calculate the error based on the next layer's weights and errors
			for j := range layer.neurons {
				error := 0.0
				for k := range nn.layers[i+1].neurons {
					error += nn.layers[i+1].neurons[k].weights[j] * outputLayerError[k]
				}
				errors[j] = error * sigmoidDerivative(layer.output[j])
			}
		} else {
			errors = outputLayerError
		}

		// Update the weights and biases
		for j := range layer.neurons {
			neuron := &nn.layers[i].neurons[j]
			for k := range neuron.weights {
				if i == 0 {
					neuron.weights[k] += learningRate * errors[j] * input[k]
				} else {
					neuron.weights[k] += learningRate * errors[j] * nn.layers[i-1].output[k]
				}
			}
			neuron.bias += learningRate * errors[j]
		}
		outputLayerError = errors
	}
}

// NewNeuron creates a new neuron with random weights and bias
func NewNeuron(inputSize int) Neuron {
	rand.Seed(time.Now().UnixNano())
	weights := make([]float64, inputSize)
	for i := range weights {
		weights[i] = rand.Float64()
	}
	return Neuron{weights, rand.Float64()}
}

// NewLayer creates a new layer with a given number of neurons
func NewLayer(neuronCount, inputSize int) Layer {
	neurons := make([]Neuron, neuronCount)
	for i := range neurons {
		neurons[i] = NewNeuron(inputSize)
	}
	return Layer{neurons, nil}
}

// NewNeuralNetwork creates a new neural network with the specified architecture
func NewNeuralNetwork(architecture []int) NeuralNetwork {
	layers := make([]Layer, len(architecture)-1)
	for i := 0; i < len(architecture)-1; i++ {
		layers[i] = NewLayer(architecture[i+1], architecture[i])
	}
	return NeuralNetwork{layers}
}

// VisualizeNetwork generates a structured textual representation of the network.
func (nn *NeuralNetwork) VisualizeNetwork() {
	fmt.Println("Neural Network Architecture\n")

	// Iterate over each layer
	for i, layer := range nn.layers {
		// Header for the layer
		fmt.Printf("Layer %d (%s):\n", i, layerType(i, len(nn.layers)))

		// Iterate over each neuron in the layer
		for j, neuron := range layer.neurons {
			// Print details for the neuron
			fmt.Printf("  Neuron %d:\n", j)
			fmt.Printf("    Weights: %v\n", neuron.weights)
			fmt.Printf("    Bias: %f\n", neuron.bias)
		}
		fmt.Println(strings.Repeat("-", 30))
	}
}

// layerType determines the type of layer (input, hidden, output)
func layerType(index, totalLayers int) string {
	if index == 0 {
		return "Input"
	} else if index == totalLayers-1 {
		return "Output"
	}
	return "Hidden"
}

func main1() {
	// Example: XOR problem
	inputs := [][]float64{
		{0, 0},
		{0, 1},
		{1, 0},
		{1, 1},
	}
	targets := [][]float64{
		{0},
		{1},
		{1},
		{0},
	}

	// Create a neural network with 2 inputs, 2 neurons in the hidden layer, and 1 output
	nn := NewNeuralNetwork([]int{2, 2, 1})

	// Train the neural network
	for epoch := 0; epoch < 10000; epoch++ {
		for i, input := range inputs {
			nn.Train(input, targets[i], 0.1)
		}
	}

	// Test the neural network
	for _, input := range inputs {
		output := nn.Feedforward(input)
		fmt.Printf("Input: %v, Predicted: %v\n", input, output)
	}

	fmt.Println()
	// After training and testing the neural network, visualize its architecture.
	nn.VisualizeNetwork()
}
