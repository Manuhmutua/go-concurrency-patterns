package main

import (
	"fmt"
	"math/rand"
	"time"
)

func fanInTimeout1(input1, input2 <-chan string) <-chan string {
	c := make(chan string)
	go func() {
		for {
			for {
				select {
				case s := <-input1:
					c <- s
				case s := <-input2:
					c <- s
				}
			}
		}
	}()
	return c
}

func main() {
	c := fanInTimeout1(boringTimeout1("Joe"), boringTimeout1("Ann"))
	for i := 0; i < 10; i++ {
		select {
		case s := <-c:
			fmt.Println(s)
		case <-time.After(1 * time.Second):
			fmt.Println("You're too slow.")
			return
		}
	}
}

func boringTimeout1(msg string) <-chan string { // Returns receive-only channel of strings.
	c := make(chan string)
	go func() { // We launch the goroutine from inside the function.
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(2e3)) * time.Millisecond)
		}
	}()

	return c // Return the channel to the caller.
}
