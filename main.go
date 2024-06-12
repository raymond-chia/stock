package main

import (
	"fmt"
	"time"

	"github.com/raymond-chia/stock/analyze"
	"github.com/raymond-chia/stock/crawler/yahoo"
	"github.com/raymond-chia/stock/domain"
)

const (
	// 日
	KDJDailyFilter = 30.0
	BBIDailyFilter = 2.0
	// 周
	KDJWeeklyFilter = 45.0
	BBIWeeklyFilter = 2.0
	// 月 算出來跟 yahoo 差距過大
)

type ID string

var IDs = map[ID]int{
	"0050": 0,
	"0056": 0,
	"2330": 0,
	"2603": 0,
	"2606": 0,
	"3037": 0,
	"3030": 0,
	"5425": 0,
	"5403": 0,
	"8155": 0,
	"3661": 0,
	"3665": 0,
	"2428": 0,
	"2474": 0,
	"3491": 0,
	"6435": 0,
	"3596": 0,
	"3592": 0,
	"3023": 0,
	"6138": 0,
	"6515": 0,
	"6488": 0,
	"3003": 0,
	"4915": 0,
	"2360": 0,
	"6683": 0,
	"6104": 0,
	"3217": 0,
	"3265": 0,
	"3526": 0,
	"6411": 0,
	"6752": 0,
	"6753": 0,
	"5371": 0,
	"6756": 0,
	"3081": 0,
	"2455": 0,
	"2379": 0,
	"6669": 0,
	"2376": 0,
	"2377": 0,
	"6279": 0,
	"6239": 0,
	"3558": 0,
	"2345": 0,
	"8478": 0,
	"6643": 0,
	"3374": 0,
	"3376": 0,
	"4919": 0,
	"1723": 0,
	"1760": 0,
	"3443": 0,
	"8261": 0,
	"8255": 0,
	"3532": 0,
	"2351": 0,
	"2357": 0,
	"2359": 0,
	"4739": 0,
	"4721": 0,
	"3680": 0,
	"8081": 0,
	"6510": 0,
	"3413": 0,
	"3035": 0,
	"1560": 0,
	"1582": 0,
	"6414": 0,
	"3533": 0,
	"3515": 0,
	"8299": 0,
	"3454": 0,
	"5474": 0,
	"4958": 0,
	"3324": 0,
	"3019": 0,
	"6491": 0,
	"5289": 0,
	"2421": 0,
	"6409": 0,
	"3527": 0,
	"2454": 0,
	"4966": 0,
	"2382": 0,
	"2356": 0,
	"2308": 0,
	"6412": 0,
	"6841": 0,
	"2404": 0,
	"2480": 0,
	"2464": 0,
	"2465": 0,
	"3131": 0,
	"3551": 0,
	"3583": 0,
	"5536": 0,
	"5009": 0,
	"2395": 0,
	"2439": 0,
	"2449": 0,
	"3227": 0,
	"6214": 0,
	"6166": 0,
	"3231": 0,
	"2453": 0,
	"2353": 0,
	"3029": 0,
	"6213": 0,
	"3563": 0,
	"4551": 0,
	"3044": 0,
	"4105": 0,
	"6271": 0,
	"3529": 0,
	"6829": 0,
	"3711": 0,
	"3587": 0,
	"8050": 0,
	"8046": 0,
	"6125": 0,
	"3693": 0,
	"3362": 0,
	"4938": 0,
	"9958": 0,
	"1514": 0,
	"1519": 0,
	"6806": 0,
	"2059": 0,
	"5904": 0,
	// 航運 可以參考 bdi ??
	"2609": 0,
	// 無人機
	"2645": 0,
	"8033": 0,
	"2634": 0,
}

// 人工檢查本益比 ? 有時候高還是可以 ?
// ROE 15%+ (財報狗)
// 營業利益率 10%+ (奇摩基本)
// CDP 不高 (月線)
// DMI 紅尖買, 藍尖賣
// Yahoo 股市: 本益比, 股利, 財務, 基本
func main() {
	fmt.Println("總共:", len(IDs))

	daily, weekly, weeklyAndUnderAverage := crawl()

	fmt.Println("# 日篩選:")
	for id, name := range daily {
		if _, ok := weekly[id]; ok {
			continue
		}
		fmt.Println(id, name)
	}
	fmt.Println("# 周篩選:")
	for id, name := range weekly {
		fmt.Println(id, name)
	}
	// TODO 今年營收成長大於 0% ?
	fmt.Println("# 周篩選且價格低於平均:")
	for id, name := range weeklyAndUnderAverage {
		fmt.Println(id, name)
	}
	fmt.Println("- 日篩選:", len(daily))
	fmt.Println("- 周篩選:", len(weekly))
}

type Name string
type Filter map[ID]Name

func crawl() (Filter, Filter, Filter) {
	daily := Filter{}
	weekly := Filter{}
	weeklyAndUnderAverage := Filter{}

	for id := range IDs {
		missing, name, data, err := yahoo.Yahoo(string(id))
		if err != nil {
			fmt.Println("an error occurs:", err)
			continue
		}
		if missing {
			fmt.Println(id, "missing data")
		}

		if !filterDaily(data) {
			daily[id] = Name(name)
		}

		if !filterWeekly(data) {
			weekly[id] = Name(name)
			if underAverage(data) {
				weeklyAndUnderAverage[id] = Name(name)
			}
		}
	}
	return daily, weekly, weeklyAndUnderAverage
}

func filterDaily(data []domain.Data) bool {
	// data = data[analyze.Max(len(data)-180, 0):]
	kdj := analyze.KDJ(data, 9)
	// macd := analyze.MACD(data)
	bbi := analyze.BullBearIndex(data)

	i := len(data) - 1
	return kdj[i].K > KDJDailyFilter ||
		bbi[i].Diff > BBIDailyFilter
}

func filterWeekly(data []domain.Data) bool {
	weekly := []domain.Data{}
	t := time.Now().Add(time.Hour * 24 * 365)
	for i := len(data) - 1; i >= 0; i-- {
		d := data[i]
		if !d.Date.Add(time.Hour * 24 * 6).Before(t) {
			continue
		}
		t = d.Date
		weekly = append([]domain.Data{d}, weekly...)
	}

	kdj := analyze.KDJ(weekly, 9)
	// macd := analyze.MACD(weekly)
	bbi := analyze.BullBearIndex(weekly)

	i := len(weekly) - 1
	return kdj[i].K > KDJWeeklyFilter ||
		bbi[i].Diff > BBIWeeklyFilter
}

// // TODO
// func filterMonthly(data []domain.Data) bool {
// 	monthly := []domain.Data{}
// 	t := time.Now().Add(time.Hour * 24 * 365)
// 	for i := len(data) - 1; i >= 0; i-- {
// 		d := data[i]
// 		if !d.Date.Add(time.Hour * 24 * 27).Before(t) {
// 			continue
// 		}
// 		t = d.Date
// 		monthly = append([]domain.Data{d}, monthly...)
// 	}

// 	kdj := analyze.KDJ(monthly, 9)
// 	// macd := analyze.MACD(monthly)
// 	bbi := analyze.BullBearIndex(monthly)

// 	i := len(monthly) - 1
// 	// TODO fix monthly bbi
// 	_ = bbi
// 	return kdj[i].K > KDJMonthlyFilter
// 	// ||
// 	// bbi[i].Diff > BBIMonthlyFilter
// }

func underAverage(data []domain.Data) bool {
	count := 1460
	if len(data)-1461 < 0 {
		count = len(data) - 1
	}
	total := 0.0
	for i := len(data) - 1; i >= 0 && i >= len(data)-1461; i-- {
		total += data[i].Close
	}
	avg := total / float64(count)
	return data[len(data)-1].Close < avg
}
