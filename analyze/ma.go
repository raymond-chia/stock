package analyze

import "github.com/raymond-chia/stock/domain"

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

func sma(data []domain.Data, i, n int) float64 {
	result := 0.0
	for _, e := range data[Max(i-n+1, 0) : i+1] {
		result += e.Close
	}
	result /= float64(n)
	return result
}
