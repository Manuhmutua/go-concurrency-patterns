package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

type Conn struct {
	str string
	wait chan bool
}


func fanInQuit(input1, input2 <-chan string) <-chan string {
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
	// worker
	quit := make(chan string)
	c := fanInQuit(boringQuit("Joe", quit), boringQuit("Ann", quit))
	for i := rand.Intn(10); i >= 0; i-- {
		fmt.Println(<-c)
	}
	quit <- "Bye"
	fmt.Printf("Sender says %q\n", <- quit )
}

func boringQuit(msg string, quit chan string) <-chan string { // Returns receive-only channel of strings.

	//
	c := make(chan string)
	go func() { // We launch the goroutine from inside the function.
		for i := 0; ; i++ {
			select {
			case c <- fmt.Sprintf("%s %d", msg, i):
				// do nothing
			case <-quit:
				// do something amazing
				time.Sleep(time.Duration(rand.Intn(5e3)) * time.Millisecond)
				quit <- "See you"
				return
			}
		}
	}()

	return c // Return the channel to the caller.
}

func Query(conns []Conn, query string) sql.Result {
	ch := make(chan sql.Result, len(conns)) //buffered
	for _, Conn := range conns{
		go func(c Conn) {
			ch <- c.DoQuery(query)
		}(Conn)
		
	}
	return  <- ch
}