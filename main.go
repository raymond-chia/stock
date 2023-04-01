package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

func main() {
	id := 3037
	missing, name, result := yahoo(id)
	if missing {
		fmt.Println(id, "missing data")
	}
	fmt.Println(name, result)
}

func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

type Data struct {
	Date   time.Time `json:"t"`
	Open   float64   `json:"o"`
	Close  float64   `json:"c"`
	High   float64   `json:"h"`
	Low    float64   `json:"l"`
	Volume int       `json:"v"`
}

func yahoo(id int) (bool, string, []Data) {
	v := url.Values{
		"v":        {"1"},
		"type":     {"ta"},
		"mkt":      {"10"},
		"sym":      {fmt.Sprint(id)},
		"perd":     {"d"},
		"_":        {fmt.Sprint(time.Now().UnixMilli())},
		"callback": {""},
	}
	target := "https://tw.quote.finance.yahoo.net/quote/q?" + v.Encode()
	resp := Must(http.Get(target))
	defer resp.Body.Close()
	b := Must(io.ReadAll(resp.Body))
	b = bytes.TrimPrefix(b, []byte("("))
	b = bytes.TrimSuffix(b, []byte(");"))
	raw := map[string]any{}
	Must(0, json.Unmarshal(b, &raw))

	mem := raw["mem"].(map[string]any)
	name := mem["name"].(string)

	type RawData struct {
		Date int `json:"t"`
		Data
	}
	b = Must(json.Marshal(raw["ta"]))
	rawData := []RawData{}
	Must(0, json.Unmarshal(b, &rawData))
	result := []Data{}
	missing := false
	for _, rawData := range rawData {
		switch {
		case rawData.Date == 0, rawData.Volume == 0:
			missing = true
		case rawData.Open == 0, rawData.Close == 0, rawData.High == 0, rawData.Low == 0:
			missing = true
		}
		date := fmt.Sprint(rawData.Date)
		data := rawData.Data
		data.Date = Must(time.Parse("20060102", date))
		result = append(result, data)
	}

	return missing, name, result
}
