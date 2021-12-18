package utils

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/PanGan21/concurrent-parser/models"
)

// CheckError panics if err is not nil
func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// MergeFanIn merges ExampleModel channels into one output channel - Fan in pattern
func MergeFanIn(ctx context.Context, inputs ...<-chan *models.ExampleModel) <-chan *models.ExampleModel {
	out := make(chan *models.ExampleModel)

	wg := &sync.WaitGroup{}
	wg.Add(len(inputs))

	sendToOutput := func(input <-chan *models.ExampleModel) {
		defer wg.Done()
		for in := range input {
			// out <- in
			select {
			case out <- in:
			case <-ctx.Done():
				fmt.Println(ctx.Err().Error())
				return
			}
		}
	}

	go func() {
		for _, input := range inputs {
			go sendToOutput(input)
		}
	}()

	go func() {
		defer close(out)
		wg.Wait()
	}()

	return out
}

func MergeErr(ctx context.Context, inputs ...<-chan error) <-chan error {
	out := make(chan error)

	wg := &sync.WaitGroup{}
	wg.Add(len(inputs))

	sendToOutput := func(input <-chan error) {
		defer wg.Done()
		for err := range input {
			select {
			case out <- err:
			case <-ctx.Done():
				fmt.Println(ctx.Err().Error())
				return
			}
		}
	}

	go func() {
		for _, input := range inputs {
			go sendToOutput(input)
		}
	}()

	go func() {
		defer close(out)
		wg.Wait()
	}()

	return out
}
