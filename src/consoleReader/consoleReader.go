package consoleReader

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

// SupplierAdapter wraps function that writes
// to channel
type SupplierAdapter func() <-chan int

func (s SupplierAdapter) Supply() <-chan int {
	return s()
}

func NewConsoleReader() SupplierAdapter {
	return scanConsole
}

func scanConsole() <-chan int {
	out := make(chan int)

	go func() {

		defer close(out)

		log.Println("Start reading STDIN")
		fmt.Println("Usage: enter number and press Enter (only numbers will be accepted)")
		fmt.Println("To exit type 'exit' and press enter")

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {

			t := scanner.Text()

			if t == "exit" {
				log.Println("exit reading STDIN")
				break
			}

			// accepts integers only
			i, err := strconv.Atoi(t)
			if err != nil {
				fmt.Println("\nonly numbers are accepted")
				continue
			}

			log.Printf("Got data from STDIN: %d\n", i)
			out <- i
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
	}()

	return out
}
