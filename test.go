package main

import (
	"fmt"
)

// worker is a function that takes a channel as an argument
// and sends integers (0 through 4) into the channel.
func worker(ch chan int) {
	for i := 0; i < 5; i++ {
		ch <- i // Send the value of i into the channel ch
	}
	close(ch) // Close the channel to indicate that no more values will be sent
}

func main() {
	ch := make(chan int) // Create a channel of type int

	go worker(ch) // Start the worker function as a Goroutine

	// The range loop receives values from the channel ch
	// The loop automatically ends when the channel is closed
	for item := range ch {
		fmt.Printf("Received: %d\n", item) // Print the received value
	}

	fmt.Println("Channel closed, main function ends.")
}
