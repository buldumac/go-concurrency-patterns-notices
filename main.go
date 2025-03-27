package main

import (
	"fmt"
	generators "github.com/buldumac/go-concurrency-patterns-notices/pattern-generators"
	"math/rand"
	"time"
)

func main() {
	// simple generator
	doneCh := make(chan interface{})
	stream1 := generators.Generator(doneCh, 10, 20, 30, 40, 50, 60)
	for vv := range stream1 {
		fmt.Println(vv)
	}

	fmt.Println("\n---delimiter---\n")

	// take only 3 values from the generator
	// the beauty here is how we can chain generators
	stream2 := generators.Take(doneCh, generators.Generator(doneCh, 10, 20, 30, 40, 50), 3)
	for vv := range stream2 {
		fmt.Println(vv)
	}

	fmt.Println("\n---delimiter---\n")

	// we want 50 of A!
	stream3 := generators.Take(doneCh, generators.Repeat(doneCh, "A"), 50)
	areThereReally50OfA := make([]interface{}, 0)
	for vv := range stream3 {
		fmt.Println(vv)
		areThereReally50OfA = append(areThereReally50OfA, vv)
	}
	fmt.Println(fmt.Sprintf("We got %d of A values", len(areThereReally50OfA)))

	fmt.Println("\n---delimiter---\n")

	// Now let's see how repeat function generator useful could ever be
	randSource := rand.New(rand.NewSource(time.Now().UnixMilli()))
	randFunc := func() interface{} { return randSource.Int() }

	stream4 := generators.RepeatFunc(doneCh, randFunc)
	for vv := range generators.Take(doneCh, stream4, 30) {
		fmt.Println(vv)
	}

	// can we read 2 more random numbers from the stream ? Sure thing!
	fmt.Println("Two more random numbers:")
	fmt.Println(<-stream4)
	fmt.Println(<-stream4)
	fmt.Println()

	// How can we stop it ? We don't need it anymore!
	// Verify if channel is open:
	val, isChOpen := <-stream4
	fmt.Println("Channel is open?", isChOpen, "value received", val)
	// just close done channel! :)
	close(doneCh)
	val, isChOpen = <-stream4
	fmt.Println("Channel is open?", isChOpen, "value received", val)
	// now the channel is closed, no more indefinitely running generator of data, because it generates interfaces
	// the default value is nil, it's an empty interface
}
