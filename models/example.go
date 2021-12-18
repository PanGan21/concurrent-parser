package models

import (
	"encoding/json"
	"strconv"
	"strings"
)

type ExampleModel struct {
	Attr1 string `json:"attr1"`
	Attr2 string `json:"attr2"`
	Attr3 string `json:"attr3"`
	Attr4 int    `json:"attr4"`
}

// NewExampleModel returns a new ExampleModel
func NewExampleModel(data []string) *ExampleModel {
	num, err := strconv.Atoi(data[3])
	if err != nil {
		panic(err)
	}
	return &ExampleModel{Attr1: data[0], Attr2: data[1], Attr3: data[2], Attr4: num}
}

// ToJson returns the marsalled json bytes for an ExampleModel
func ToJson(exampleModel *ExampleModel) []byte {
	body, err := json.Marshal(&exampleModel)
	if err != nil {
		panic(err)
	}
	return body
}

// Transform transforms each attribute of an ExampleModel to uppercase
func (exampleModel *ExampleModel) Transform() {
	exampleModel.Attr1 = strings.ToUpper(exampleModel.Attr1)
	exampleModel.Attr2 = strings.ToUpper(exampleModel.Attr2)
	exampleModel.Attr3 = strings.ToUpper(exampleModel.Attr3)
}
