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

package yahoo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/regel/tinkerbell/pkg/finance/types"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Response struct {
	Chart Chart `json:"chart"`
}

type Chart struct {
	Result []Result `json:"result"`
}

type Meta struct {
	Timezone             string `json:"timezone"`
	ExchangeTimezoneName string `json:"exchangeTimezoneName"`
}

type Result struct {
	Meta       Meta      `json:"meta"`
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

func timeWithinRange(t time.Time, from time.Time, to time.Time) bool {
	return (t.Equal(from) || t.After(from)) && (t.Equal(to) || t.Before(to))
}

func getUrl(baseUrl string, ticker string, interval string, from time.Time, to time.Time) string {
	base, err := url.Parse(baseUrl)
	if err != nil {
		panic("Can't parse Yahoo Finance base url")
	}
	if to.Sub(from) < time.Duration(24*time.Hour) {
		to = from.AddDate(0, 0, 1)
	}
	values := url.Values{
		"interval":   []string{interval},
		"period1":    []string{strconv.FormatInt(from.Unix(), 10)},
		"period2":    []string{strconv.FormatInt(to.Unix(), 10)},
		"region":     []string{"US"},
		"corsDomain": []string{"com.finance.yahoo"},
	}
	relative := &url.URL{
		Path:     "/v8/finance/chart/" + ticker,
		RawQuery: values.Encode(),
	}

	return base.ResolveReference(relative).String()
}

func GetOhlc(c context.Context, client *http.Client, baseUrl string, ticker string, interval string, from time.Time, to time.Time) ([]types.Ohlc, error) {
	queryUrl := getUrl(baseUrl, ticker, interval, from, to)
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
	points := make([]types.Ohlc, 0)
	quote := response.Chart.Result[0].Indicators.Quote[0]
	loc, err := time.LoadLocation(response.Chart.Result[0].Meta.ExchangeTimezoneName)
	if err != nil {
		return nil, err
	}
	for j, timestamp := range response.Chart.Result[0].Timestamps {
		t := time.Unix(timestamp, 0).In(loc)
		if timeWithinRange(t, from, to) {
			ohlc := types.Ohlc{
				Ticker:    ticker,
				Timestamp: t,
				Volume:    quote.Volume[j],
				Open:      quote.Open[j],
				High:      quote.High[j],
				Low:       quote.Low[j],
				Close:     quote.Close[j],
			}
			points = append(points, ohlc)
		}
	}
	return points, err
}
