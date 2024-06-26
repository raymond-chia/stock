package analyze

import "github.com/raymond-chia/stock/domain"

// TODO find out why this is very different from Yahoo
func MACD(data []domain.Data) []domain.MACDData {
	// padding
	// we will remove it later
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
		d.Date = data[i].Date
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

func dif(ema12, ema26 float64) float64 {
	return ema12 - ema26
}
