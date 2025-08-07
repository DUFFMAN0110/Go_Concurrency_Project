package main

import (
	"fmt"
)

// simulate a task that takes time
func doTask(id int, ch chan string) {
	// Simulate work
	ch <- fmt.Sprintf("Task %d completed", id)
}

func main() {
	// Number of concurrent tasks
	numTasks := 50
	ch := make(chan string)

	// Launch goroutines
	for i := 1; i <= numTasks; i++ {
		go doTask(i, ch)
	}

	// Collect results
	for i := 0; i < numTasks; i++ {
		fmt.Println(<-ch)
	}
}