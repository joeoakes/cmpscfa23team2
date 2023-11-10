package main

import (
	"fmt"
	"log"
	"math/rand"

	"gorgonia.org/gorgonia"
	"gorgonia.org/tensor"
)

// go clean -modcache
// go env -w GOPRIVATE=github.com/google/flatbuffers
// go get -u go4.org/unsafe/assume-no-moving-gc
// go mod tidy

func main() {
	// Set up a seed to ensure reproducibility
	rand.Seed(42)

	// Create a new graph
	g := gorgonia.NewGraph()

	// Define the network architecture
	inputSize := 2  // two inputs for the XOR problem
	hiddenSize := 2 // size of the hidden layer
	outputSize := 1 // one output for the XOR problem

	// Create nodes for the inputs and the target output
	inputs := gorgonia.NewMatrix(g, tensor.Float64, gorgonia.WithShape(4, inputSize), gorgonia.WithName("inputs"))    // Ensure the shape is (4, inputSize)
	targets := gorgonia.NewMatrix(g, tensor.Float64, gorgonia.WithShape(4, outputSize), gorgonia.WithName("targets")) // Ensure the shape is (4, outputSize)

	// Initialize weights and biases for the hidden layer
	wHidden := gorgonia.NewMatrix(g, tensor.Float64, gorgonia.WithShape(inputSize, hiddenSize), gorgonia.WithName("wHidden"), gorgonia.WithInit(gorgonia.GlorotU(1)))
	bHidden := gorgonia.NewMatrix(g, tensor.Float64, gorgonia.WithShape(1, hiddenSize), gorgonia.WithName("bHidden"), gorgonia.WithInit(gorgonia.Zeroes())) // Shape (1, hiddenSize) for broadcasting

	// Initialize weights and biases for the output layer
	wOut := gorgonia.NewMatrix(g, tensor.Float64, gorgonia.WithShape(hiddenSize, outputSize), gorgonia.WithName("wOut"), gorgonia.WithInit(gorgonia.GlorotU(1)))
	bOut := gorgonia.NewMatrix(g, tensor.Float64, gorgonia.WithShape(1, outputSize), gorgonia.WithName("bOut"), gorgonia.WithInit(gorgonia.Zeroes())) // Shape (1, outputSize) for broadcasting

	// Initialize the datasets
	inputData := tensor.New(tensor.Of(tensor.Float64), tensor.WithShape(4, inputSize), tensor.WithBacking([]float64{0, 0, 0, 1, 1, 0, 1, 1}))
	targetData := tensor.New(tensor.Of(tensor.Float64), tensor.WithShape(4, outputSize), tensor.WithBacking([]float64{0, 1, 1, 0}))

	// Define the neuron operations
	hiddenLayer := gorgonia.Must(gorgonia.Add(gorgonia.Must(gorgonia.Mul(inputs, wHidden)), bHidden))
	hiddenLayer = gorgonia.Must(gorgonia.Rectify(hiddenLayer)) // ReLU activation
	output := gorgonia.Must(gorgonia.Add(gorgonia.Must(gorgonia.Mul(hiddenLayer, wOut)), bOut))
	pred := gorgonia.Must(gorgonia.Sigmoid(output))

	// Define the loss function (binary cross-entropy)
	loss := gorgonia.Must(gorgonia.BinaryXent(pred, targets))

	// Define our gradient descent optimizer with a fixed learning rate
	learningRate := 0.1
	solver := gorgonia.NewVanillaSolver(gorgonia.WithLearnRate(learningRate))

	// We must first compile the nodes of our graph which are learnable
	learnables := gorgonia.Nodes{wHidden, bHidden, wOut, bOut}
	if _, err := gorgonia.Grad(loss, learnables...); err != nil {
		log.Fatal(err)
	}

	// Create a VM to run the program on
	m := gorgonia.NewTapeMachine(g, gorgonia.BindDualValues(wHidden, bHidden, wOut, bOut))

	// Train the network
	for epoch := 0; epoch < 1000; epoch++ {
		// Let the input and target data into the corresponding nodes
		gorgonia.Let(inputs, inputData)
		gorgonia.Let(targets, targetData)

		if err := m.RunAll(); err != nil {
			log.Fatal(err)
		}

		// Gradients are automatically computed after running the VM.
		// Apply the gradients to update the weights and biases
		if err := solver.Step(gorgonia.NodesToValueGrads(learnables)); err != nil {
			log.Fatal(err)
		}

		// Reset the VM and the graph for the next iteration
		m.Reset()

		// Print out the loss every 100 epochs
		if epoch%100 == 0 {
			fmt.Printf("Epoch %d: Loss = %v\n", epoch, loss.Value())
		}
	}

	// Output the final trained weights and biases
	fmt.Println("Trained weights (hidden layer): \n", wHidden.Value())
	fmt.Println("Trained biases (hidden layer): \n", bHidden.Value())
	fmt.Println("Trained weights (output layer): \n", wOut.Value())
	fmt.Println("Trained biases (output layer): \n", bOut.Value())

	// Make predictions on the trained model
	for i := 0; i < 4; i++ {
		xT := tensor.New(tensor.Of(tensor.Float64), tensor.WithShape(1, inputSize), tensor.WithBacking(inputData.Data().([]float64)[i*inputSize:(i+1)*inputSize])) // Ensure the shape is (1, inputSize)
		gorgonia.Let(inputs, xT)
		if err := m.RunAll(); err != nil {
			log.Fatalf("Failed during prediction: %v", err)
		}

		fmt.Printf("Input: %v, Predicted: %v\n", xT, pred.Value())

		// Reset the VM after each prediction
		m.Reset()
	}
}
