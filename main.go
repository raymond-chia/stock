package main

import (
	"fmt"

	"github.com/raymond-chia/stock/analyze"
	"github.com/raymond-chia/stock/crawler/yahoo"
)

func main() {
	id := 3037
	missing, name, data, err := yahoo.Yahoo(id)
	if err != nil {
		panic(err)
	}
	if missing {
		fmt.Println(id, "missing data")
	}
	kdj := analyze.KDJ(data, 9)
	fmt.Println(name, kdj)
}
