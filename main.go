package main

import (
	"context"
	"fmt"
	generators "github.com/buldumac/go-concurrency-patterns-notices/pattern-generators"
	tee "github.com/buldumac/go-concurrency-patterns-notices/tee-channel-pattern"
	"math/rand"
	"sync"
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

	fmt.Println("\n---delimiter---\n")
	fmt.Println("Tee Channel Pattern")
	// let's use Tee Channel Pattern

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	doneCh = make(chan interface{}) // we closed it previously, assign it to new channel

	someStringsAsInterfaces := generators.Take(doneCh, generators.Repeat(doneCh, "Hello", "how", "are", "you"), 7)
	someStringsAsStrings := generators.ToString[string](doneCh, someStringsAsInterfaces)

	wg := sync.WaitGroup{}
	ch1, ch2, ch3 := tee.TeePlease[string](ctx, someStringsAsStrings)
	var chans []<-chan string
	chans = append(chans, ch1, ch2, ch3)
	wg.Add(len(chans))
	for i, chx := range chans {
		chx := chx // starting from go1.22 loop variables have per-iteration scope
		// https://go.dev/blog/loopvar-preview#the-fix
		go func(i int) {
			defer wg.Done() // we can pass waitgroup as a parameter, but we'll go with closure, Done is thread-safe
			chanName := fmt.Sprintf("chan-%d", i)
			var vals []string
			for v := range chx {
				vals = append(vals, v)
			}
			fmt.Println(chanName, "values", vals, "length", len(vals))
		}(i + 1)
	}
	wg.Wait() // don't use time.Sleep, that's not gopher's way!

}
