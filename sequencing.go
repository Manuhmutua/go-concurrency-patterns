package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Message struct {
	str string
	wait chan bool
}

func fanInSequence(msg1, msg2 <-chan Message) <-chan Message{
	c := make(chan Message)
	go func() {for {c <- <- msg1} }()
	go func() {for {c <- <- msg2} }()
	return c
}

func boringSequence(msg Message) <-chan Message { // Returns receive-only channel of strings.
	c := make(chan Message)
	waitForIt := make(chan bool)
	go func() { // We launch the goroutine from inside the function.
		for i := 0; ; i++ {
			c <- Message{ fmt.Sprintf("%s: %d", msg.str, i), waitForIt}
			time.Sleep(time.Duration(rand.Intn(2e3)) * time.Millisecond)
			<-waitForIt
		}
	}()

	return c // Return the channel to the caller.
}

func main()  {
	c := fanInSequence(boringSequence(Message{}), boringSequence(Message{}))
	for i := 0; i < 5 ; i++  {
		msg1 := <-c; fmt.Println("John")
		msg2 := <-c; fmt.Println("Ann")
		msg1.wait <- true
		msg2.wait <- true
	}
	fmt.Println("You're both boring; I'm leaving.")
}
