package utils

import (
	"sync"
)

const (
	defaultBufSize = 5
)

// ProcessInParallel Use for transformations that take an input array and return an output array by
// performing the same action on every input element. This action can be simple filtering or complete
// transformation into another type. Every input element can produce any number of output elements
// (including 0) and ordering is not kept.
func ProcessInParallel[IN any, OUT any](
	inArr []IN,
	action func(in IN, outChan chan<- OUT),
) []OUT {
	var wg sync.WaitGroup
	outChan := make(chan OUT, defaultBufSize)
	doneCounterChan := make(chan struct{})

	for _, in := range inArr {
		wg.Go(func() {
			action(in, outChan)
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
