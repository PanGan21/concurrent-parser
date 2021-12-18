package parser

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"

	"github.com/PanGan21/concurrent-parser/models"
)

const NUM_TO_COMPARE = 20

// Read returns the values of a reader to a channel
func Read(ctx context.Context, r io.Reader) (<-chan *models.ExampleModel, <-chan error, error) {
	reader := csv.NewReader(r)
	modelChan := make(chan *models.ExampleModel)
	errChan := make(chan error)
	go func() {
		defer close(modelChan)
		defer close(errChan)
		for {
			data, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				errChan <- err
				return
			}

			select {
			case modelChan <- models.NewExampleModel(data):
			case <-ctx.Done():
				fmt.Println(ctx.Err().Error())
				return
			}
		}
	}()
	return modelChan, errChan, nil
}

// Transform transforms the exampleModels concurrently
func Transform(ctx context.Context, exampleModels <-chan *models.ExampleModel) (<-chan *models.ExampleModel, <-chan error, error) {
	transformedModels := make(chan *models.ExampleModel)
	errChan := make(chan error)
	go func() {
		defer close(transformedModels)
		defer close(errChan)
		for model := range exampleModels {
			if model.Attr1 == "error" {
				errChan <- errors.New("unrecoverable error")
				return
			}
			model.Transform()
			select {
			case transformedModels <- model:
			case <-ctx.Done():
				fmt.Println(ctx.Err().Error())
				return
			}

		}
	}()
	return transformedModels, errChan, nil
}

// Filter filters out the rows with Attr4 greater than NUM_TO_COMPARE
func Filter(ctx context.Context, exampleModels <-chan *models.ExampleModel) (<-chan *models.ExampleModel, <-chan error, error) {
	filteredModels := make(chan *models.ExampleModel)
	errChan := make(chan error)
	go func() {
		defer close(filteredModels)
		defer close(errChan)
		for model := range exampleModels {
			if model.Attr4 > NUM_TO_COMPARE {
				filteredModels <- model
			}
		}
	}()
	return filteredModels, errChan, nil
}
