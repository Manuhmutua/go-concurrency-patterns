package main

import (
	"fmt"
	"math/rand"
	"time"
)

type MessageSelect struct {
	str string
	wait chan bool
}

func fanInSelect(msg1, msg2 <-chan MessageSelect) <-chan MessageSelect {
	c := make(chan MessageSelect)
	go func() {
		for {
			for {
				select {
				case s := <-msg1:
					c <- s
				case s := <-msg2:
					c <- s
				}
			}
		}
	}()
	return c
}

func boringSelect(msg MessageSelect) <-chan MessageSelect { // Returns receive-only channel of strings.
	c := make(chan MessageSelect)
	waitForIt := make(chan bool)
	go func() { // We launch the goroutine from inside the function.
		for i := 0; ; i++ {
			c <- MessageSelect{ fmt.Sprintf("%s: %d", msg.str, i), waitForIt}
			time.Sleep(time.Duration(rand.Intn(2e3)) * time.Millisecond)
			<-waitForIt
		}
	}()

	return c // Return the channel to the caller.
}

func main()  {
	c := fanInSelect(boringSelect(MessageSelect{}), boringSelect(MessageSelect{}))
	for i := 0; i < 5 ; i++  {
		msg1 := <-c; fmt.Println("John")
		msg2 := <-c; fmt.Println("Ann")
		msg1.wait <- true
		msg2.wait <- true
	}
	fmt.Println("You're both boring; I'm leaving.")
}