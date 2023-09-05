package main

import "fmt"

func main() {

	messages := make(chan string)

	//sending data from one go routine to another
	go func() { messages <- "ping" }()

	msg := <-messages
	fmt.Println(msg)
}
