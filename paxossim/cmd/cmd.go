package main

import (
	"fmt"
	"github.com/1xyz/paxossim/paxossim/queue"
)

func main() {
	q := queue.NewQueue()
	entries := []string{"5", "4", "3", "2", "1"}

	go func(entries []string) {
		for _, e := range entries {
			q.Enqueue(e)
			fmt.Printf("e %v\n", e)
		}
	}(entries)

	//time.Sleep(10 * time.Second)

	for i := 0; i < len(entries); i++ {
		actualEntry := q.WaitForItem()
		expectedEntry := entries[i]
		fmt.Printf("actualEntry %v i = %d expectedEntry = %v\n", actualEntry, i, expectedEntry)
	}
}
