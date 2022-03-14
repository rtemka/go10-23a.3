package consoleReader

import (
	"bufio"
	"fmt"
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

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {

			t := scanner.Text()

			if t == "exit" {
				break
			}

			// accepts integers only
			i, err := strconv.Atoi(t)
			if err != nil {
				continue
			}
			out <- i
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
	}()

	return out
}
