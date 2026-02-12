//revive:disable:var-naming For now it's okay to have a generic name
package utils

import (
	"sync"
)

const (
	defaultBufSize = 5 // TODO later: choose appropriate buffer size (balance performance vs resource usage)
)

// ProcessInParallel Use for transformations that take an input array and return an output array by
// performing the same action on every input element. This action can be simple filtering or complete
// transformation into another type. Every input element can produce any number of output elements
// (including 0) and ordering is not kept.
// Think of this function like Flux::flatMap from Java
func ProcessInParallel[IN any, OUT any](
	inArr []IN,
	onEachElem func(in IN) []OUT,
) []OUT {
	return ProcessInParallelChan(inArr, func(in IN, outChan chan<- OUT) {
		for _, out := range onEachElem(in) {
			outChan <- out
		}
	})
}

// ProcessInParallelChan See ProcessInParallel, but the action performed on each element sends its output
// onto an open channel. Prefer this variant if said channel is being written to by other goroutines.
func ProcessInParallelChan[IN any, OUT any](
	inArr []IN,
	onEachElem func(in IN, outChan chan<- OUT),
) []OUT {
	var wg sync.WaitGroup
	outChan := make(chan OUT, defaultBufSize)
	doneCounterChan := make(chan struct{})

	for _, in := range inArr {
		wg.Go(func() {
			onEachElem(in, outChan)
			doneCounterChan <- struct{}{}
		})
	}

	wg.Go(func() {
		closeAfterSignalCount(len(inArr), doneCounterChan)
		close(outChan)
	})

	var outArr []OUT
	for out := range outChan {
		outArr = append(outArr, out)
	}

	wg.Wait()

	return outArr
}

func closeAfterSignalCount(target int, signalChannel chan struct{}) {
	defer close(signalChannel)

	if target == 0 {
		return
	}

	signalCount := 0
	for {
		select {
		case <-signalChannel:
			signalCount++
			if signalCount >= target {
				return
			}
		}
	}
}
