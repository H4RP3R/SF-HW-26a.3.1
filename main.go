package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"pipeline/logger"
)

var ErrInvalidInput = errors.New("invalid input: please enter a number")

var logsTarget string

type LoggerStatus string

type Stage func(done chan struct{}, inChan <-chan int, log *log.Logger) <-chan int

// pipeline represents a series of stages that process a stream of
// integers.
type pipeline struct {
	stages []Stage
	log    *log.Logger
}

// AddStage appends a new stage to the pipeline's list of stages.
func (p *pipeline) AddStage(stage Stage) {
	p.stages = append(p.stages, stage)
	p.log.Printf("New stage added (stages=%d)", len(p.stages))
}

// Run starts the pipeline, processing data from the dataSource channel and
// returning a new channel that emits the final pipeline products.
func (p *pipeline) Run(done chan struct{}, dataSource <-chan int) <-chan int {
	c := dataSource
	for _, stage := range p.stages {
		c = stage(done, c, p.log)
	}

	p.log.Println("Start the pipeline")
	return c
}

// NewPipeLine creates and returns a new instance of a pipeline.
func NewPipeLine(logger *log.Logger) *pipeline {
	p := &pipeline{}
	p.log = logger

	p.log.Println("Created new pipeline")
	return p
}

// readInput reads user input from the terminal and emits each number to the
// outChan until the done channel is closed.
func readInput(done chan struct{}) <-chan int {
	outChan := make(chan int)
	reader := bufio.NewReader(os.Stdin)

	go func() {
		defer close(outChan)
		for {
			input, err := reader.ReadString('\n')
			if err != nil {
				log.Println(err)
			}
			input = strings.TrimSpace(input)
			num, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println(ErrInvalidInput)
			} else {
				select {
				case outChan <- num:
				case <-done:
					return
				}
			}
		}
	}()

	return outChan
}

// waitForInterrupt creates a goroutine that waits for a SIGINT signal and
// exits the program when it's received.
func waitForInterrupt(done chan struct{}) {
	var wg sync.WaitGroup
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Press Ctrl+C to exit...")
		fmt.Printf("buffer size: %d, delay: %v\n", bufferSize, bufferDelay)
		<-sigChan
		close(done)
		fmt.Println("\nBye!")
	}()

	wg.Wait()
}

// display consumes a channel of pipeline products and prints them to the
// console until the done channel is closed.
func display(done chan struct{}, products <-chan int) {
	go func() {
		for {
			select {
			case num := <-products:
				fmt.Printf("processed: %d\n", num)
			case <-done:
				return
			}
		}
	}()
}

func main() {
	flag.DurationVar(&bufferDelay, "delay", 15*time.Second, "buffer delay")
	flag.IntVar(&bufferSize, "size", 24, "buffer size")
	flag.StringVar(&logsTarget, "log", "none", "destination for log output")
	flag.Parse()

	done := make(chan struct{})

	logger, err := logger.New(logsTarget)
	if err != nil {
		log.Fatal("Cant't implement logger: ", err)
	}
	p := NewPipeLine(logger)
	p.AddStage(filterMultiplesOfThree)
	p.AddStage(filterNegativeNumbers)
	p.AddStage(buffering)

	dataSource := readInput(done)
	products := p.Run(done, dataSource)
	display(done, products)

	waitForInterrupt(done)
}
