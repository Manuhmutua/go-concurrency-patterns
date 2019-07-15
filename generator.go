package main

import (
	"fmt"
	"math/rand"
	"time"
)

// go-concurrency-patterns
func main() {
	c := boringGenerator("boring!") // function return channel.
	for i := 0; i < 5; i++ {
		fmt.Printf("You Say: %q\n", <-c)
	}
	fmt.Println("You're boring; I'm leaving")

	joe := boringGenerator("Joe")
	ann := boringGenerator("Ann")
	for i := 0; i < 5 ; i++ {
		fmt.Println(<-joe)
		fmt.Println(<-ann)
	}
	fmt.Println("You're boring; I'm leaving")
}

func boringGenerator(msg string) <-chan string { // Returns receive-only channel of strings.
	c := make(chan string)
	go func() { // We launch the goroutine from inside the function.
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()

	return c // Return the channel to the caller.
}

