package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {

	ctrlC := make(chan os.Signal, 1)
	signal.Notify(ctrlC, os.Interrupt)

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("Running update")
		case <-ctrlC:
			ticker.Stop()
			fmt.Println("\nFinished with CTRL + C")
			return
		}
	}
}
