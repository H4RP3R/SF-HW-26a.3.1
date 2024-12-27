package main

import (
	"log"
	"time"

	ringBuf "github.com/H4RP3R/ring_buffer"
)

var (
	bufferDelay time.Duration
	bufferSize  int
)

// filterNegativeNumbers returns a channel that emits only non-negative
// numbers from the input channel.
func filterNegativeNumbers(done chan struct{}, inChan <-chan int, log *log.Logger) <-chan int {
	outChan := make(chan int)

	go func() {
		defer close(outChan)
		for num := range inChan {
			log.Printf("[filterNegativeNumbers] received value: %v", num)
			if num >= 0 {
				select {
				case outChan <- num:
					log.Printf("[filterNegativeNumbers] value %v has been passed further", num)
				case <-done:
					return
				}
			} else {
				log.Printf("[filterNegativeNumbers] value %v was discarded", num)
			}
		}
	}()

	return outChan
}

// filterMultiplesOfThree returns a channel that emits only numbers that
// are multiples of 3.
func filterMultiplesOfThree(done chan struct{}, inChan <-chan int, log *log.Logger) <-chan int {
	outChan := make(chan int)

	go func() {
		defer close(outChan)
		for num := range inChan {
			log.Printf("[filterMultiplesOfThree] received value: %v", num)
			if num != 0 && num%3 == 0 {
				select {
				case outChan <- num:
					log.Printf("[filterMultiplesOfThree] value %v has been passed further", num)
				case <-done:
					return
				}
			} else {
				log.Printf("[filterMultiplesOfThree] value %v was discarded", num)
			}
		}
	}()

	return outChan
}

// buffering collects numbers from the input channel into a buffer and emits
// them to the output channel at intervals specified by bufferDelay.
func buffering(done chan struct{}, inChan <-chan int, log *log.Logger) <-chan int {
	outChan := make(chan int)
	buffer, err := ringBuf.New[int](bufferSize)
	if err != nil {
		panic(err)
	}

	go func() {
		for num := range inChan {
			buffer.Push(num)
			log.Printf("[buffering] received value: %v", num)
		}
	}()

	ticker := time.NewTicker(bufferDelay)
	go func() {
		defer func() {
			ticker.Stop()
			close(outChan)
		}()
		for {
			select {
			case <-ticker.C:
				for !buffer.IsEmpty() {
					if num, ok := buffer.Pop(); ok {
						log.Printf("[buffering] value %v has been passed further", num)
						outChan <- num
					}
				}
			case <-done:
				return
			}
		}
	}()

	return outChan
}
