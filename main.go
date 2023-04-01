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
	data = data[analyze.Max(len(data)-180, 0):]
	kdj := analyze.KDJ(data, 9)
	macd := analyze.MACD(data)
	fmt.Println(name)
	for i := range data {
		fmt.Printf("%v %+v %+v\n", data[i].Date, kdj[i], macd[i])
	}
}
