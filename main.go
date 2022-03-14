package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"pipeline/src/consolePrinter"
	"pipeline/src/consoleReader"
	"pipeline/src/ringBuffer"
	"pipeline/src/worker"
)

// buffer settings
const (
	bufferSize          = 20
	bufferFlushInterval = time.Second * 10
)

// supplier is responsible for data supplying
// for the pipeline
type supplier interface {
	Supply() <-chan int
}

// consumer is responsible for data consuming
// from the pipeline
type consumer interface {
	Consume(<-chan int)
}

// processor is the main working unit
// of pipeline (e.g. doing something with data
// and pass it down the line)
type processor interface {
	Filter(done <-chan struct{}, in <-chan int) <-chan int
}

// buffer aggregates data from the pipeline
type buffer interface {
	Read() (int, error)
	Write(int) error
	IsFull() bool
}

// bufferController is managing underlying buffer.
// It controls in/out channels
// and manages buffer flush intervals
type bufferController struct {
	b             buffer
	flushInterval time.Duration
}

// controller is conducts of entire pipeline.
// It launch/cancel all processes relative to pipeline
type controller struct {
	bufCtl bufferController
	// this can be done with slice of processors
	negFilter processor // negative integers filter
	d3Filter  processor // divisible by three and zero integers filter
	done      chan struct{}
}

// manageBuffer writes incoming data to underlying buffer
// and flushes it to the channel after provided period of time
func (bc *bufferController) manageBuffer(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)

	// flush closure to read all data from the buffer
	flush := func() {
		for {
			el, err := bc.b.Read()
			if err != nil {
				break
			}
			out <- el
		}
		log.Println("buffer has been flushed")
	}

	go func() {
		defer close(out)
		defer flush()

		for {
			// if buffer is full we wait
			// for the flush time or stop signal
			if bc.b.IsFull() {
				select {
				case <-done:
					return
				case <-time.After(bc.flushInterval):
					flush()
					continue
				}
			}

			// if buffer is not full
			// we wtrite incoming data to buffer
			// and listen for the stop signals
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				err := bc.b.Write(v)
				if err != nil {
					fmt.Println(err)
				}
				log.Println("written to the buffer", v)
			case <-time.After(bc.flushInterval):
				flush()
			}
		}
	}()

	return out
}

// positive returns provided number
// and result of checking if it's positive
func positive(x int) (int, bool) {
	return x, x >= 0
}

// diviByThreeNotZero returns provided number
// and result of checking if it's divisible by 3
// and not zero
func diviByThreeNotZero(x int) (int, bool) {
	return x, x%3 == 0 && x != 0
}

func main() {

	buf, err := ringBuffer.NewBuffer(bufferSize)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	bc := bufferController{
		b:             buf,
		flushInterval: bufferFlushInterval,
	}

	// supplier and consumer for the pipeline
	var supplier supplier = consoleReader.NewConsoleReader()
	var consumer consumer = consolePrinter.NewConsolePrinter()

	// controller assembling
	ctl := controller{
		bufCtl:    bc,
		negFilter: worker.NewIntProcessor(1, positive),
		d3Filter:  worker.NewIntProcessor(2, diviByThreeNotZero),
		done:      make(chan struct{}),
	}

	// in this case done channel is not necessary
	// since we stops only if supplier closes supply line.
	// But if we had some stop logic (like os.Interrupt signal) then done channel
	// would help us to stop pipeline in any given point in time.
	// In that case we need to pass done channel to supplier and consumer
	// so they know that we stop processing/producing
	defer close(ctl.done)

	// input from supplier
	pipelineInput := supplier.Supply()

	// processing
	stageOne := ctl.negFilter.Filter(ctl.done, pipelineInput)
	stageTwo := ctl.d3Filter.Filter(ctl.done, stageOne)
	pipelineOutput := ctl.bufCtl.manageBuffer(ctl.done, stageTwo)

	// output to consumer
	consumer.Consume(pipelineOutput)
}
