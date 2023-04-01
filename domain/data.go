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
	K   float64
	D   float64
	J   float64
	RSV float64
}
