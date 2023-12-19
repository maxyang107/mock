package main

import (
	"fmt"
	"mock"
)

type Resp struct {
	Age   uint32  `json:"age"`
	Score float64 `json:"score"`
	Rank  int16   `json:"rank"`
	Name  string  `json:"name"`
	Img   string  `json:"img"`
	Email string  `json:"email"`
	Min   float64 `json:"min"`
	Max   int     `json:"max"`
}

func main() {
	var (
		re     = Resp{}
		mockOp = []mock.Option{}
	)
	mockOp = append(mockOp, mock.WithContainsFieldSourceString("Age", []string{"12", "15", "26"}))
	mockOp = append(mockOp, mock.WithContainsFieldSourceString("Score", []string{"12.56", "15.32", "26.1"}))
	mockOp = append(mockOp, mock.WithContainsFieldSourceString("Rank", []string{"1", "13", "2"}))
	mockOp = append(mockOp, mock.WithMinMaxByField("Min", 3, 7), mock.WithFloatLen("%.3f"))
	mockOp = append(mockOp, mock.WithMinMaxByField("Max", 3, 700))
	mockOp = append(mockOp, mock.WithContainsFieldSourceString("Img", []string{"http://www.baidu.com?adjae&daskjewa1321.jpe", "http://www.baidu.com?adjae&daskjewa1321.jpg", "http://www.baidu.com?adjae&daskjewa1321.png"}))
	mock.MockGen(&re, mockOp...)
	fmt.Println(re)
}
