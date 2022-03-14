package consolePrinter

import "log"

// ConsumerAdapter wraps channel processing function
type ConsumerAdapter func(<-chan int)

func (c ConsumerAdapter) Consume(in <-chan int) {
	c(in)
}

func NewConsolePrinter() ConsumerAdapter {
	return print
}

func print(in <-chan int) {
	for i := range in {
		log.Printf("\nData received: %d\n", i)
	}
}
