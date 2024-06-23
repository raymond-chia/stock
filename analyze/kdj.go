package analyze

import (
	"math"

	"github.com/raymond-chia/stock/domain"
)

// https://www.futunn.com/hk/learn/detail-what-is-kdj-64858-220831019
// n 請見 rsv
func KDJ(data []domain.Data, period int) []domain.KDJData {
	// add default kdj
	// we will remove it later
	result := []domain.KDJData{{K: 50, D: 50, J: 50}}
	for i := range data {
		kdj := domain.KDJData{}
		kdj.Date = data[i].Date
		kdj.RSV = rsv(data[Max(i-period, 0) : i+1])
		kdj.K = 2.0/3.0*result[i].K + 1.0/3.0*kdj.RSV
		kdj.D = 2.0/3.0*result[i].D + 1.0/3.0*kdj.K
		kdj.J = 3*kdj.K - 2*kdj.D
		result = append(result, kdj)
	}
	return result[1:]
}

// RSV=（Ct-Ln）÷（Hn-Ln）×100
// Ct 為第 n 日收盤價
// Ln 為 n 日內最低價
// Hn 為 n 日內最高價
func rsv(data []domain.Data) float64 {
	c := data[len(data)-1].Close
	l := data[0].Low
	h := data[0].High
	for _, d := range data {
		l = math.Min(l, d.Low)
		h = math.Max(h, d.High)
	}
	return (c - l) / (h - l) * 100.0
}
