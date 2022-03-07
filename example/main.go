package main

import (
	"encoding/json"
	"fmt"

	"github.com/rohanraj7316/logger"
)

type Data struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type JsonTest struct {
	StatusCode string `json:"statusCode"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
	Data       []Data `json:"data"`
}

func checkErr(err error) bool {
	if err == nil {
		return false
	}
	fmt.Println(err)
	return true
}

func main() {
	lOptions := logger.DefaultOptions()
	err := logger.Configure(lOptions)
	checkErr(err)

	// testing string
	f := logger.Field{
		Key:   "string",
		Value: "test string",
	}
	logger.Error("testing string", f)

	// testring json
	j := JsonTest{
		StatusCode: "httpOK",
		Code:       200,
		Message:    "testing encryption",
		Data: []Data{
			{
				Name: "Unio",
				Age:  76,
			},
		},
	}
	fb, err := json.Marshal(&j)
	if !checkErr(err) {
		f := logger.Field{
			Key:   "json",
			Value: string(fb),
		}
		logger.Info("testing json", f)
	}
}
