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

package finance

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	ohlcUrl = "https://query2.finance.yahoo.com/v8/finance/chart/%s?interval=%s&period1=%d&period2=%d&region=US&corsDomain=com.finance.yahoo"
)

func getSupportedRanges() []string {
	return []string{
		"1d",
		"5d",
		"1mo",
		"3mo",
		"6mo",
		"1y",
		"2y",
		"5y",
		"10y",
		"ytd",
		"max",
	}
}

type Ohlc struct {
	Ticker    string
	Timestamp time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    int64
}

type Response struct {
	Chart Chart `json:"chart"`
}

type Chart struct {
	Result []Result `json:"result"`
}

type Result struct {
	Timestamps []int64   `json:"timestamp"`
	Indicators Indicator `json:"indicators"`
}

type Indicator struct {
	Quote []Quote `json:"quote"`
}

type Quote struct {
	Volume []int64   `json:"volume"`
	Open   []float64 `json:"open"`
	High   []float64 `json:"high"`
	Low    []float64 `json:"low"`
	Close  []float64 `json:"close"`
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ReadOhlc(c context.Context, client *http.Client, ticker string, interval string, startTime time.Time, endTime time.Time) ([]Ohlc, error) {
	if !contains(getSupportedRanges(), interval) {
		return nil, errors.New("Invalid interval")
	}
	queryUrl := fmt.Sprintf(ohlcUrl, ticker, interval, startTime.Unix(), endTime.Unix())
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

	response := &Response{}
	err = json.NewDecoder(res.Body).Decode(response)
	if err != nil {
		return nil, err
	}
	points := make([]Ohlc, 0)
	quote := response.Chart.Result[0].Indicators.Quote[0]
	for j, timestamp := range response.Chart.Result[0].Timestamps {
		ohlc := Ohlc{
			Ticker:    ticker,
			Timestamp: time.Unix(timestamp, 0),
			Volume:    quote.Volume[j],
			Open:      quote.Open[j],
			High:      quote.High[j],
			Low:       quote.Low[j],
			Close:     quote.Close[j],
		}
		points = append(points, ohlc)
	}
	return points, err
}
