package consolePrinter

// ConsumerAdapter wraps channel processing function
type ConsumerAdapter func(<-chan int)

func (c ConsumerAdapter) Consume(in <-chan int) {
	c(in)
}

func NewConsolePrinter() ConsumerAdapter {
	return print
}

func print(in <-chan int) {
	for range in {
	}
}
