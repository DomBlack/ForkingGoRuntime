package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("  main called")

	var barrier sync.WaitGroup

	barrier.Add(1)
	go func() {
		fmt.Println("  go routine called")
		barrier.Done()
	}()
	barrier.Wait()

	fmt.Println("  main finished")
}
