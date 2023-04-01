package yahoo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/raymond-chia/stock/domain"
)

func Yahoo(id int) (bool, string, []domain.Data, error) {
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
	resp, err := http.Get(target)
	if err != nil {
		return false, "", nil, fmt.Errorf("fail to reach Yahoo with error[%w]", err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", nil, fmt.Errorf("fail to read from Yahoo with error[%w]", err)
	}
	b = bytes.TrimPrefix(b, []byte("("))
	b = bytes.TrimSuffix(b, []byte(");"))
	raw := map[string]any{}
	err = json.Unmarshal(b, &raw)
	if err != nil {
		return false, "", nil, fmt.Errorf("fail to unmarshal response from Yahoo with error[%w]", err)
	}

	mem := raw["mem"].(map[string]any)
	name := mem["name"].(string)

	type RawData struct {
		Date   int     `json:"t"`
		Open   float64 `json:"o"`
		Close  float64 `json:"c"`
		High   float64 `json:"h"`
		Low    float64 `json:"l"`
		Volume int     `json:"v"`
	}
	b, err = json.Marshal(raw["ta"])
	if err != nil {
		return false, "", nil, fmt.Errorf("fail to marshal ta data from Yahoo with error[%w]", err)
	}
	rawData := []RawData{}
	err = json.Unmarshal(b, &rawData)
	if err != nil {
		return false, "", nil, fmt.Errorf("fail to unmarshal ta data from Yahoo with error[%w]", err)
	}
	result := []domain.Data{}
	missing := false
	for _, rawData := range rawData {
		switch {
		case rawData.Date == 0, rawData.Volume == 0:
			missing = true
			continue
		case rawData.Open == 0, rawData.Close == 0, rawData.High == 0, rawData.Low == 0:
			missing = true
			continue
		}
		rawDate := fmt.Sprint(rawData.Date)
		date, err := time.Parse("20060102", rawDate)
		if err != nil {
			missing = true
			continue
		}
		result = append(result, domain.Data{
			Date:   date,
			Open:   rawData.Open,
			Close:  rawData.Close,
			High:   rawData.High,
			Low:    rawData.Low,
			Volume: rawData.Volume,
		})
	}

	return missing, name, result, nil
}
