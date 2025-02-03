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

	// 成交量不低
	// rsi 替代 kdj

	VolumeThresholdWeekly  = 2
	VolumeThresholdMonthly = 3
)

type ID string

// 定存股
// 玉山金 2884 利息 & 股票
// 0056
// 聯華 1229 利息 & 股票
// 鈊象 3293
// 00679B 美債 ... 利息低
var IDs = map[ID]int{
	"0050": 0,
	"0056": 0,
	"2330": 0,
	"2603": 0,
	"2606": 0,
	"3037": 0,
	"3030": 0,
	"5425": 0,
	"3661": 0,
	"3665": 0,
	"2428": 0,
	"2474": 0,
	"3491": 0,
	"3596": 0,
	"3592": 0,
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
	"6753": 0,
	"5371": 0,
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
	"3443": 0,
	"8255": 0,
	"3532": 0,
	"2351": 0,
	"2357": 0,
	"2359": 0,
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
	"4958": 0,
	"3324": 0,
	"3019": 0,
	"6491": 0,
	"5289": 0,
	"2421": 0,
	"6409": 0,
	"2454": 0,
	"4966": 0,
	"2382": 0,
	"2356": 0,
	"2308": 0,
	"6412": 0,
	"2404": 0,
	"2480": 0,
	"2464": 0,
	"3131": 0,
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
	"2353": 0,
	"3029": 0,
	"6213": 0,
	"3563": 0,
	"4551": 0,
	"3044": 0,
	"4105": 0,
	"6271": 0,
	"3529": 0,
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
	// # 航運 可以參考 bdi ??
	"2609": 0,
	// # 無人機
	"2645": 0,
	"2634": 0,
}

// 每股盈餘 (奇摩基本) x 營收年增% x 本益比 = 預期股價 ??
// ROE 15%+ (財報狗)
// 營業利益率 10%+ (奇摩基本)
// RSI 當 RSI 高於 70, 表示股價處於超買區, 可能會回調; 當 RSI 低於 30, 表示股價處於超賣區, 可能會反彈
// DMI +DI尖賣; -DI尖買
// DMI +DI 上漲程度; -DI 下降程度; AD https://tw.stock.yahoo.com/news/%E6%8A%80%E8%A1%93%E5%88%86%E6%9E%90-dmi%E6%8C%87%E6%A8%99-dmi%E5%8B%95%E5%90%91%E6%8C%87%E6%A8%99-%E5%A4%9A%E7%A9%BA%E6%96%B9%E5%90%91-%E8%B6%A8%E5%8B%A2%E5%8B%95%E8%83%BD-130256849.html
// MACD 長期趨勢. 負正黃金交叉; 正負死亡交叉
// CDP 短線趨勢. 分開買; 收縮賣 ??
// Yahoo 股市: 本益比, 股利, 財務, 基本
// Goodinfo 財務評分表
// 股價預估值 ?? https://www.findbillion.com/twstock/1231
func main() {
	fmt.Println("總共:", len(IDs))

	result := crawl()

	fmt.Println("# 日篩選:")
	for id, name := range result.Daily {
		if _, ok := result.Weekly[id]; ok {
			continue
		}
		fmt.Println(id, name)
	}
	fmt.Println("# 周篩選:")
	for id, name := range result.Weekly {
		fmt.Println(id, name)
	}
	// TODO 今年營收成長大於 0% ?
	fmt.Println("# 周篩選且價格低於平均:")
	for id, name := range result.WeeklyAndUnderAverage {
		fmt.Println(id, name)
	}
	fmt.Println("# 日近期有 DMI -DI 尖:")
	for id, name := range result.DailyDMI {
		fmt.Println(id, name)
	}
	fmt.Println("# 周成交增長", VolumeThresholdWeekly)
	for id, name := range result.weeklyVolume {
		fmt.Println(id, name)
	}
	fmt.Println("# 月成交增長", VolumeThresholdMonthly)
	for id, name := range result.monthlyVolume {
		fmt.Println(id, name)
	}
}

type Name string
type IDToName map[ID]Name

type CrawlResult struct {
	Daily                 IDToName
	Weekly                IDToName
	WeeklyAndUnderAverage IDToName
	DailyDMI              IDToName
	weeklyVolume          IDToName
	monthlyVolume         IDToName
}

func crawl() CrawlResult {
	daily := IDToName{}
	weekly := IDToName{}
	weeklyAndUnderAverage := IDToName{}
	dailyDMI := IDToName{}
	weeklyVolume := IDToName{}
	monthlyVolume := IDToName{}

	for id := range IDs {
		missing, _name, data, err := yahoo.Yahoo(string(id))
		if err != nil {
			fmt.Println("an error occurs:", err)
			continue
		}
		if missing {
			fmt.Println(id, "missing data")
		}
		name := Name(_name)
		weeklyData := weeklyData(data)

		if !filterDaily(data) {
			daily[id] = name
		}
		if !filterWeekly(weeklyData) {
			weekly[id] = name
			if underAverage(data) {
				weeklyAndUnderAverage[id] = name
			}
		}

		if dmiPeak(data) {
			dailyDMI[id] = name
		}

		if this, previous := volumeSurge(data, 5); this > previous*VolumeThresholdWeekly {
			weeklyVolume[id] = name
		}
		if this, previous := volumeSurge(data, 20); this > previous*VolumeThresholdMonthly {
			monthlyVolume[id] = name
		}
	}
	return CrawlResult{
		Daily:                 daily,
		Weekly:                weekly,
		WeeklyAndUnderAverage: weeklyAndUnderAverage,
		DailyDMI:              dailyDMI,
		weeklyVolume:          weeklyVolume,
		monthlyVolume:         monthlyVolume,
	}
}

func filterDaily(data []domain.Data) bool {
	kdj := analyze.KDJ(data, 9)
	bbi := analyze.BullBearIndex(data)

	i := len(data) - 1
	return kdj[i].K > KDJDailyFilter ||
		bbi[i].Diff > BBIDailyFilter
}

func filterWeekly(data []domain.Data) bool {
	kdj := analyze.KDJ(data, 9)
	bbi := analyze.BullBearIndex(data)

	i := len(data) - 1
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

func weeklyData(data []domain.Data) []domain.Data {
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
	return weekly
}

func dmiPeak(data []domain.Data) bool {
	dmi := analyze.DMI(data, 14)
	for i := 0; i < 7; i++ {
		if analyze.DMIMinusPeak(dmi, i) {
			for j := 0; j < i; j++ {
				if analyze.DMIPlusPeak(dmi, j) {
					return false
				}
			}
			return true
		}
	}
	return false
}

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

// ignore holidays
func volumeSurge(data []domain.Data, interval int) (int, int) {
	l := len(data)
	this := 0
	previous := 0
	for _, d := range data[l-interval : l] {
		this += d.Volume
	}
	for _, d := range data[l-interval*2 : l-interval] {
		previous += d.Volume
	}
	return this, previous
}
