package domain

import "time"

type Data struct {
	Date   time.Time
	Open   float64
	Close  float64
	High   float64
	Low    float64
	Volume int
}

type KDJData struct {
	Date time.Time
	K    float64
	D    float64
	J    float64
	RSV  float64
}

type MACDData struct {
	Date  time.Time
	EMA12 float64
	EMA26 float64
	DIF   float64
	MACD  float64
	OSC   float64
}

type BullBearIndex struct {
	Date time.Time
	SMA3 float64
	BBI  float64
	Diff float64
}
