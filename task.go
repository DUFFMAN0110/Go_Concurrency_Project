package main;

import (
	"fmt"
	"math/rand"
	"time"
)

// RunTask simulates a task that takes random time to complete
// The task is printing prime numbers 1 100
func RunTask(id int) {
	fmt.Printf("Task %d started\n", id)
	duration := time.Duration(rand.Intn(1000)+500) * time.Millisecond
	time.Sleep(duration)
	fmt.Printf("Task %d completed in %v\n", id, duration)
}