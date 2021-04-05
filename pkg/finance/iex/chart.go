// Copyright The TB Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package iex

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/regel/tinkerbell/pkg/finance/types"
	"log"
	"math"
	"net/http"
	"net/url"
	"time"
)

type Response struct {
	Chart []Chart `json:"chart"`
}

type Chart struct {
	Date    string  `json:"date"`
	Updated int64   `json:"updated"`
	Volume  int64   `json:"volume"`
	Open    float64 `json:"open"`
	High    float64 `json:"high"`
	Low     float64 `json:"low"`
	Close   float64 `json:"close"`
}

func timeWithinRange(t time.Time, from time.Time, to time.Time) bool {
	return (t.Equal(from) || t.After(from)) && (t.Equal(to) || t.Before(to))
}

func getUrl(baseUrl string, token string, ticker string, interval string, startTime time.Time, endTime time.Time) string {
	base, err := url.Parse(baseUrl)
	if err != nil {
		panic("Can't parse IEX Cloud base url")
	}

	max := int(math.Ceil(time.Since(startTime).Hours() / 24))
	rangeValue := fmt.Sprintf("%dd", max)
	values := url.Values{
		"token":   []string{token},
		"symbols": []string{ticker},
		"range":   []string{rangeValue},
		"types":   []string{"chart"},
	}
	relative := &url.URL{
		Path:     "/v1/stock/market/batch",
		RawQuery: values.Encode(),
	}

	return base.ResolveReference(relative).String()
}

func ReadOhlc(c context.Context, client *http.Client, baseUrl string, token string, ticker string, interval string, startTime time.Time, endTime time.Time) ([]types.Ohlc, error) {
	queryUrl := getUrl(baseUrl, token, ticker, interval, startTime, endTime)
	req, err := http.NewRequest(http.MethodGet, queryUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(c, time.Duration(5*time.Second))
	defer cancel()
	req = req.WithContext(ctx)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Non-OK HTTP status: %d", res.StatusCode)
	}

	response := map[string]Response{}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	points := make([]types.Ohlc, 0)
	chart := response[ticker].Chart
	for _, quote := range chart {
		timestamp, _ := time.Parse("2006-01-02", quote.Date)
		if timeWithinRange(timestamp, startTime, endTime) {
			point := types.Ohlc{
				Ticker:    ticker,
				Timestamp: timestamp,
				Volume:    quote.Volume,
				Open:      quote.Open,
				High:      quote.High,
				Low:       quote.Low,
				Close:     quote.Close,
			}
			points = append(points, point)
		}
	}
	return points, err
}
