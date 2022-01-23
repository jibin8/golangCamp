package main

import (
	"fmt"

	. "github.com/hhxsv5/go-redis-memory-analysis"
)

func main() {
	analysis, err := NewAnalysisConnection("127.0.0.1", 6379, "")
	if err != nil {
		r := fmt.Sprintf("connect redis err:%s", err.Error())
		fmt.Println(r)
		return
	}
	defer analysis.Close()

	analysis.Start([]string{"#", ":"})

	err = analysis.SaveReports("./reports")
	if err != nil {
		r := fmt.Sprintf("save result data err:%s", err.Error())
		fmt.Println(r)
	}
}
