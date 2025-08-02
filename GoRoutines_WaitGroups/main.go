package main

import (
	"fmt"
	"log"
	"sync"
)

func printTask(id string, wg *sync.WaitGroup) {
	defer wg.Done() // using defer to ensure it should be executed even any error occur

	for i := 0; i < 5; i++ {
		fmt.Printf("Task %s - step %d\n", id, i)
	}
}

func main() {
	var waitGroup sync.WaitGroup
	tasks := []string{
		"task 1",
		"task 2",

	}
	//alwyas call Add(n) before executing the go routines but not inside the go routines
	waitGroup.Add(len(tasks))
	// length of task slice is passed as there is need to wait for all go routines compelte execution
	//  as i want to execute each task in separate goroutine

	for _, task := range tasks {
		go printTask(task, &waitGroup)
	}

	log.Println("main done")
	waitGroup.Wait()

}
