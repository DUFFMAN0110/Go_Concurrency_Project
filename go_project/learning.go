package main

import "fmt"

func main() {
	var name string = "Nestor"
	name2 := "Hector"
	name3, age3 := "Cam", 420

	fmt.Println("Hello, world!")
	fmt.Println(name)
	fmt.Println(name2)
	fmt.Println(name3)
	fmt.Println(age3)

	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}

	// While-style loop
	i := 0
	for i < 5 {
		fmt.Println(i)
		i++
	}

	// Infinite loop
	count := 0
	for {
		fmt.Println("looping forever")
		fmt.Println("count: ", count)
		count++
		if count > 4800 {
			break
		}

	}

	// //
	// //
	// //
	// fmt.Println("Go Sub Routine Section:")

	// go sayHello() // runs sayHello() concurrently

	// func sayHello() {
	// 	fmt.Println("Hello from goroutine!")
	// }

	// // Channels (communication between goroutines)
	// ch := make(chan string)

	// go func() {
	// 	ch <- "ping"
	// }()

	// msg := <-ch
	// fmt.Println(msg)

	// Stage 1: Basic concurrency function
	package main

	import (
		"fmt"
		"time"
	)

	func greet(id int, ch chan string) {
		msg := fmt.Sprintf("Hello from goroutine %d", id)
		time.Sleep(time.Second) // simulate work
		ch <- msg
	}

	func main() {
		ch := make(chan string)
		for i := 1; i <= 3; i++ {
			go greet(i, ch)
		}
		for i := 1; i <= 3; i++ {
			fmt.Println(<-ch)
		}
	}

}
