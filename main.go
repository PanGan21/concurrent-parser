package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/PanGan21/concurrent-parser/models"
	"github.com/PanGan21/concurrent-parser/parser"
	"github.com/PanGan21/concurrent-parser/utils"
)

const (
	parrallelExec = 3
)

func main() {
	file, err := os.Open("./data.csv")
	utils.CheckError(err)
	defer file.Close()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	exampleModelList := make([]<-chan *models.ExampleModel, 0)
	errorList := make([]<-chan error, 0)

	models, errorChan, err := parser.Read(ctx, file)
	utils.CheckError(err)
	errorList = append(errorList, errorChan)

	// Fan out to multiple workers
	// Number of workers started = parrallelExec
	for i := 0; i < parrallelExec; i++ {
		filteredModels, errorChan, err := parser.Filter(ctx, models)
		utils.CheckError(err)
		errorList = append(errorList, errorChan)

		transformedModels, errorChan, err := parser.Transform(ctx, filteredModels)
		utils.CheckError(err)
		errorList = append(errorList, errorChan)

		exampleModelList = append(exampleModelList, transformedModels)
	}

	// Fan in the results of the workers
	for i := range utils.MergeFanIn(ctx, exampleModelList...) {
		fmt.Println(i)
	}

	for err := range utils.MergeErr(ctx, errorList...) {
		log.Fatal(err)
	}

}
