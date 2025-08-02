package main

import (
	"fmt"
	"sync"
)

var msg string
var wg sync.WaitGroup
var mu sync.Mutex

func updateMessage(data string) {
	defer wg.Done()
	mu.Lock()
	msg = data
	mu.Unlock()
}

func main() {
	msg = "Hello Nepal!"
	wg.Add(2)
	go updateMessage("Hello Nepalese")
	go updateMessage("Hi Nepalese")

	wg.Wait()

	fmt.Println(msg)
}

//you might get proper output even there is chance for race condition occuring while running program by go run 
// use go run -race . , to check is there a chance to occur race condition
// if there is chance of race condition then warning will appear.
// In above code comment the mutex lock and unlock to see the race condition warning.
// there are two routines executing concurrently , don't know which will run when, and both are accessing same variable msg causing the chance of race condition
// so locking the variable while performing action so that another goroutine is not allow to access the variable at that instance of time. Once operation is completed
// then we should unlock the variable allowing other goroutine to access it if need
// mutex  mean mutual exclusion, avoid two jobs accessing same resource simultaneously
