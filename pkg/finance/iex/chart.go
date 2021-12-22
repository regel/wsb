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
	"strings"
	"sync"
	"time"
)

const (
	batchMaxLen = 100
)

type Response struct {
	Chart []Chart `json:"chart"`
}

type Chart struct {
	Date    string  `json:"date"`
	Minute  string  `json:"minute,omitempty"`
	Updated int64   `json:"updated"`
	Volume  int64   `json:"volume"`
	Open    float64 `json:"open"`
	High    float64 `json:"high"`
	Low     float64 `json:"low"`
	Close   float64 `json:"close"`
}

func chunkSlice(slice []string, chunkSize int) [][]string {
	var chunks [][]string
	for {
		if len(slice) == 0 {
			break
		}

		// necessary check to avoid slicing beyond
		// slice capacity
		if len(slice) < chunkSize {
			chunkSize = len(slice)
		}

		chunks = append(chunks, slice[0:chunkSize])
		slice = slice[chunkSize:]
	}

	return chunks
}

func timeWithinRange(t time.Time, from time.Time, to time.Time) bool {
	return (t.Equal(from) || t.After(from)) && (t.Equal(to) || t.Before(to))
}

func getBatchUrl(baseUrl string, token string, tickers []string, interval string, from time.Time, to time.Time) string {
	base, err := url.Parse(baseUrl)
	if err != nil {
		panic("Can't parse IEX Cloud base url")
	}

	max := int(math.Ceil(time.Since(from).Hours() / 24))
	if max < 5 {
		max = 5
	}
	rangeValue := fmt.Sprintf("%dd", max)
	if to.Sub(from) <= time.Duration(24*time.Hour) {
		rangeValue = "1d"
	}
	values := url.Values{
		"token":   []string{token},
		"symbols": []string{strings.Join(tickers, ",")},
		"range":   []string{rangeValue},
		"types":   []string{"chart"},
	}
	relative := &url.URL{
		Path:     "/v1/stock/market/batch",
		RawQuery: values.Encode(),
	}

	return base.ResolveReference(relative).String()
}

func decodeChart(chart []Chart, ticker string, from time.Time, to time.Time) []types.Ohlc {
	points := make([]types.Ohlc, 0)
	for _, quote := range chart {
		timestamp, _ := time.Parse("2006-01-02", quote.Date)
		if quote.Minute != "" {
			var h, m int
			n, err := fmt.Sscanf(quote.Minute, "%d:%d", &h, &m)
			if err == nil && n == 2 {
				timestamp = timestamp.Add(time.Hour*time.Duration(h) +
					time.Minute*time.Duration(m))
			}
		}
		if timeWithinRange(timestamp, from, to) {
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
	return points
}

func (p Provider) GetOhlc(c context.Context, client *http.Client, ticker string, interval string, from time.Time, to time.Time) ([]types.Ohlc, error) {
	slice := []string{ticker}
	queryUrl := getBatchUrl(p.IexCloudQueryUrl, p.IexCloudSecretToken, slice, interval, from, to)
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
	chart := response[ticker].Chart
	return decodeChart(chart, ticker, from, to), nil
}

func (p Provider) GetOhlcBatch(wg *sync.WaitGroup, chartChan chan *types.Chart, c context.Context, client *http.Client, tickers []string, interval string, from time.Time, to time.Time) {
	chunks := chunkSlice(tickers, batchMaxLen)
	for _, chunk := range chunks {
		wg.Add(1)
		go func(slice []string, window string, from time.Time, to time.Time) {
			queryUrl := getBatchUrl(p.IexCloudQueryUrl, p.IexCloudSecretToken, slice, interval, from, to)
			req, err := http.NewRequest(http.MethodGet, queryUrl, nil)
			if err != nil {
				log.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(c, time.Duration(5*time.Second))
			defer cancel()
			req = req.WithContext(ctx)
			res, err := client.Do(req)
			if err != nil {
				return
			}
			defer res.Body.Close()
			if res.StatusCode != http.StatusOK {
				println(fmt.Sprintf("Non-OK HTTP status: %d", res.StatusCode))
			}

			response := map[string]Response{}
			err = json.NewDecoder(res.Body).Decode(&response)
			if err != nil {
				return
			}
			for ticker := range response {
				chart := response[ticker].Chart
				points := decodeChart(chart, ticker, from, to)
				out := &types.Chart{
					Ohlc:   points,
					Ticker: ticker,
				}
				chartChan <- out
			}
			wg.Done()
		}(chunk, interval, from, to)
	}

}

func (p Provider) BatchSupported() bool {
	return true
}
