//* Goroutine → lets a function run concurrently.
//* Channel → lets goroutines communicate/send data to each other.

//* In other words:

//* Goroutine = worker
//* Channel = pipe through which workers send data

package main

import "fmt"

func worker(id int, ch chan string) {
	ch <- fmt.Sprintf("Worker %d finished", id)
}

func main() {
	ch := make(chan string)

	//* The order can change because all three goroutines are running concurrently.
	go worker(1, ch)
	go worker(2, ch)
	go worker(3, ch)

	fmt.Println(<-ch)
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}
