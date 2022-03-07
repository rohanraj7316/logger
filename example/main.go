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
	err := logger.Configure()
	checkErr(err)

	// Type: number
	f := logger.Field{
		Key:   "number",
		Value: 2,
	}
	logger.Info("number", f)

	// Type: bool
	f = logger.Field{
		Key:   "bool",
		Value: true,
	}
	logger.Info("bool", f)

	// Type: nil
	f = logger.Field{
		Key:   "nil",
		Value: nil,
	}
	logger.Info("nil", f)

	// Type: string
	f = logger.Field{
		Key:   "string",
		Value: "test string",
	}
	logger.Info("string", f)

	// Type: struct
	s := JsonTest{
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
	f = logger.Field{
		Key:   "struct",
		Value: s,
	}
	logger.Error("struct", f)

	// Type: []byte
	fb, err := json.Marshal(&s)
	if !checkErr(err) {
		f = logger.Field{
			Key:   "[]byte",
			Value: fb,
		}
		logger.Error("[]byte", f)
	}

	// Type: map[string]interface{}
	mI := map[string]interface{}{
		"status":     "200",
		"statusCode": "httpOK",
		"message":    "testing encryption",
		"data": []map[string]interface{}{
			{
				"name": "Unio",
				"age":  76,
			},
		},
	}
	f = logger.Field{
		Key:   "map[string]interface",
		Value: mI,
	}
	logger.Error("struct", f)
}
