package analyze

import "github.com/raymond-chia/stock/domain"

func MACD(data []domain.Data) []domain.MACDData {
	ema12 := []float64{data[0].Close}
	ema26 := []float64{data[0].Close}
	for i, d := range data {
		ema12 = append(ema12, ema(d.Close, ema12[i], 12))
		ema26 = append(ema26, ema(d.Close, ema26[i], 26))
	}
	ema12 = ema12[1:]
	ema26 = ema26[1:]

	result := []domain.MACDData{}
	for i := range ema12 {
		d := domain.MACDData{}
		d.EMA12 = ema12[i]
		d.EMA26 = ema26[i]
		d.DIF = dif(d.EMA12, d.EMA26)
		if i == 0 {
			d.MACD = ema(d.DIF, d.DIF, 9)
			d.OSC = d.DIF - d.MACD
		} else {
			d.MACD = ema(d.DIF, result[i-1].DIF, 9)
			d.OSC = d.DIF - d.MACD
		}
		result = append(result, d)
	}
	return result
}

// a = 2 / (n + 1)
// n 代表天數
//
// St = a x Pt + (1 - a) x St-1
func ema(now, past float64, n int) float64 {
	a := float64(n) + 1.0
	a = 2.0 / a

	now = a * now
	past = (1 - a) * past
	return now + past
}

func dif(ema12, ema26 float64) float64 {
	return ema12 - ema26
}
