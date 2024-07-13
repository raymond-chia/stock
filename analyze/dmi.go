package analyze

import (
	"math"

	"github.com/raymond-chia/stock/domain"
)

// https://zh.wikipedia.org/wiki/%E5%8B%95%E5%90%91%E6%8C%87%E6%95%B8
// https://tw.tradingview.com/support/solutions/43000502250/
func DMI(data []domain.Data, period int) []domain.DMI {
	var plusDM, minusDM, tr []float64
	for i := 1; i < len(data); i++ {
		upMove := data[i].High - data[i-1].High
		downMove := data[i-1].Low - data[i].Low
		switch {
		case upMove > downMove && upMove > 0:
			plusDM = append(plusDM, upMove)
			minusDM = append(minusDM, 0)
		case downMove > upMove && downMove > 0:
			plusDM = append(plusDM, 0)
			minusDM = append(minusDM, downMove)
		default:
			plusDM = append(plusDM, 0)
			minusDM = append(minusDM, 0)
		}
		// https://zh.wikipedia.org/wiki/%E7%9C%9F%E5%AF%A6%E6%B3%A2%E5%8B%95%E5%B9%85%E5%BA%A6%E5%9D%87%E5%80%BC
		tr = append(tr, math.Max(data[i].High, data[i-1].Close)-math.Min(data[i].Low, data[i-1].Close))
	}

	result := []domain.DMI{}
	plusDM14 := 0.0
	minusDM14 := 0.0
	tr14 := 0.0
	// len(Data) = len(plusDM) + 1
	lenDiff := 1
	for i := period - 1; i < len(data)-lenDiff; i++ {
		if i == period-1 {
			plusDM14 = Sum(plusDM[i-period+1 : i+1])
			minusDM14 = Sum(minusDM[i-period+1 : i+1])
			tr14 = Sum(tr[i-period+1 : i+1])
		} else {
			plusDM14 = ema(plusDM[i], plusDM14, period)
			minusDM14 = ema(minusDM[i], minusDM14, period)
			tr14 = ema(tr[i], tr14, period)
		}

		plusDI := (plusDM14 / tr14) * 100
		minusDI := (minusDM14 / tr14) * 100

		result = append(result, domain.DMI{
			Date:  data[i+lenDiff].Date,
			Plus:  plusDI,
			Minus: minusDI,
		})
	}

	return result
}

func DMIPlusPeak(dmi []domain.DMI, offset int) bool {
	first := dmi[len(dmi)-3-offset].Plus
	second := dmi[len(dmi)-2-offset].Plus
	third := dmi[len(dmi)-1-offset].Plus
	return second > first && second > third && (Max(second-first, second-third) > 1.0)
}

func DMIMinusPeak(dmi []domain.DMI, offset int) bool {
	first := dmi[len(dmi)-3-offset].Minus
	second := dmi[len(dmi)-2-offset].Minus
	third := dmi[len(dmi)-1-offset].Minus
	return second > first && second > third && (Max(second-first, second-third) > 1.0)
}
