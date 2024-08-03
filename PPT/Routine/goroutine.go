package main

import (
	"fmt"
	"time"
)

func sayHello() {
	fmt.Println("Hello from the goroutine!")
}

func main() {

	// Start a new goroutine
	go sayHello()

	// Give the goroutine some time to execute
	time.Sleep(1 * time.Second)

	fmt.Println("Main function")
}
