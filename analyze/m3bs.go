package analyze

import "github.com/raymond-chia/stock/domain"

// https://peggy0501.pixnet.net/blog/post/16240317
func BullBearIndex(data []domain.Data) []domain.BullBearIndex {
	result := []domain.BullBearIndex{}
	for i := range data {
		bbi := domain.BullBearIndex{}
		bbi.Date = data[i].Date
		bbi.SMA3 = sma(data, i, 3)
		// bbi.BBI = (sma(data, i, 3) + sma(data, i, 6) + sma(data, i, 12) + sma(data, i, 24)) / 4.0
		bbi.BBI = (sma(data, i, 3) + sma(data, i, 6) + sma(data, i, 9) + sma(data, i, 12)) / 4.0
		bbi.Diff = bbi.SMA3 - bbi.BBI

		result = append(result, bbi)
	}
	return result
}
