package worker

import "log"

// IntProcessor represents
// integers processing instance
type IntProcessor struct {
	id     int
	filter func(i int) (int, bool)
}

// NewIntProcessor returns new instance of IntProcessor
func NewIntProcessor(id int, f func(i int) (int, bool)) *IntProcessor {
	return &IntProcessor{id: id, filter: f}
}

// Filter applies filtering function
// on the incoming stream of integers
// and writes result to the output channel
func (p *IntProcessor) Filter(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		defer p.stopMsg()

		for {
			select {

			case <-done:
				return

			case i, ok := <-in:
				if !ok {
					return
				}
				// filtering logic
				if v, ok := p.filter(i); ok {
					log.Printf("\nWorker #%d -> Processed value: %d -> Result: Passed\n", p.id, v)
					select {
					case out <- v:
					case <-done:
						return
					}
				} else {
					log.Printf("\nWorker #%d -> Processed value: %d -> Result: Filtered out\n", p.id, v)
				}
			}
		}

	}()
	return out
}

func (p *IntProcessor) stopMsg() {
	log.Printf("\nWorker #%d -> stoped\n", p.id)
}
