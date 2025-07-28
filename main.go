package main;

import (
	"flag"
	"fmt"
	"os"
	"sync"
)

func main() {
	var numTasks int
	flag.IntVar(&numTasks, "tasks", 1, "Number of concurrent tasks to run")
	flag.Parse()

	if numTasks < 1 || numTasks > 100 {
		fmt.Println("Error: Number of tasks must be between 1 and 100")
		os.Exit(1)
	}

	fmt.Printf("Starting %d concurrent tasks...\n", numTasks)

	var wg sync.WaitGroup
	wg.Add(numTasks)

	for i := 1; i <= numTasks; i++ {
		go func(taskID int) {
			defer wg.Done()
			RunTask(taskID)
		}(i)
	}

	wg.Wait()
	fmt.Println("All tasks completed.")
}
