package main

import (
	"fmt"

	"github.com/raymond-chia/stock/analyze"
	"github.com/raymond-chia/stock/crawler/yahoo"
)

const (
	KDJFilter = 50.0
	BBIFilter = 2.0
)

var IDs = map[int]int{
	2330: 0,
	3037: 0,
	5425: 0,
	8155: 0,
	3665: 0,
	2428: 0,
	2474: 0,
	3491: 0,
	6435: 0,
	3596: 0,
	3023: 0,
	3059: 0,
	6138: 0,
	6488: 0,
	3003: 0,
	4915: 0,
	2360: 0,
	6683: 0,
	6104: 0,
	3217: 0,
	3526: 0,
	6411: 0,
	6756: 0,
	3081: 0,
	2455: 0,
	2379: 0,
	6669: 0,
	2376: 0,
	2377: 0,
	6279: 0,
	3558: 0,
	2345: 0,
	8478: 0,
	6643: 0,
	3374: 0,
	4919: 0,
	1723: 0,
	1760: 0,
	3443: 0,
	8261: 0,
	3532: 0,
	2351: 0,
	3739: 0,
	4721: 0,
	3680: 0,
	8081: 0,
	6510: 0,
	3413: 0,
	3035: 0,
	1560: 0,
	6414: 0,
	3533: 0,
	3515: 0,
	8299: 0,
	3454: 0,
	4958: 0,
	3324: 0,
	3019: 0,
	6491: 0,
	5289: 0,
	2421: 0,
	6409: 0,
	3527: 0,
	2454: 0,
	4966: 0,
	2382: 0,
	2356: 0,
	2308: 0,
	6412: 0,
	6841: 0,
	2404: 0,
	2464: 0,
	3131: 0,
	3551: 0,
	3583: 0,
	5536: 0,
	1504: 0,
	2312: 0,
	2395: 0,
	2439: 0,
	2449: 0,
	3227: 0,
	5351: 0,
	6214: 0,
	6166: 0,
	3231: 0,
	2453: 0,
	2353: 0,
	3029: 0,
	6213: 0,
	3563: 0,
	4551: 0,
	3044: 0,
	4105: 0,
	6271: 0,
}

func main() {
	for id := range IDs {
		missing, name, data, err := yahoo.Yahoo(id)
		if err != nil {
			fmt.Println("an error occurs:", err)
			continue
		}
		if missing {
			fmt.Println(id, "missing data")
		}
		// data = data[analyze.Max(len(data)-180, 0):]
		kdj := analyze.KDJ(data, 9)
		macd := analyze.MACD(data)
		bbi := analyze.BullBearIndex(data)

		i := len(data) - 1
		if kdj[i].K > KDJFilter {
			continue
		}
		if bbi[i].Diff > BBIFilter {
			continue
		}
		fmt.Println(id, name)
		fmt.Printf("\tKDJ: %+v\n\tMACD: %+v\n\t多空乖離: %+v\n", kdj[i], macd[i], bbi[i])
	}
}
