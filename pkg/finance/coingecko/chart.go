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

package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/regel/wsb/pkg/finance/types"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

func timeWithinRange(t time.Time, from time.Time, to time.Time) bool {
	return (t.Equal(from) || t.After(from)) && (t.Equal(to) || t.Before(to))
}

func getUrl(baseUrl string, ticker string, interval string, from time.Time, to time.Time) string {
	base, err := url.Parse(baseUrl)
	if err != nil {
		panic("Can't parse Coingecko base url")
	}
	now := time.Now()
	// 1/7/14/30/90/180/365/max
	time_ago := "1"
	if now.Sub(from) > time.Duration(365*24*time.Hour) {
		time_ago = "max"
	} else if now.Sub(from) > time.Duration(180*24*time.Hour) {
		time_ago = "365"
	} else if now.Sub(from) > time.Duration(90*24*time.Hour) {
		time_ago = "180"
	} else if now.Sub(from) > time.Duration(30*24*time.Hour) {
		time_ago = "90"
	} else if now.Sub(from) > time.Duration(14*24*time.Hour) {
		time_ago = "30"
	} else if now.Sub(from) > time.Duration(7*24*time.Hour) {
		time_ago = "14"
	} else if now.Sub(from) > time.Duration(1*24*time.Hour) {
		time_ago = "7"
	}

	values := url.Values{
		"days":        []string{time_ago},
		"vs_currency": []string{"usd"},
	}
	relative := &url.URL{
		Path:     fmt.Sprintf("/api/v3/coins/%s/ohlc", ticker),
		RawQuery: values.Encode(),
	}

	return base.ResolveReference(relative).String()
}

func (p Provider) GetOhlc(c context.Context, client *http.Client, ticker string, interval string, from time.Time, to time.Time) ([]types.Ohlc, error) {
	queryUrl := getUrl(p.CoingeckoQueryUrl, ticker, interval, from, to)
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

	var response [][]float64
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	points := make([]types.Ohlc, 0)
	for _, point := range response {
		t := time.Unix(int64(point[0]/1000.0), 0)
		if timeWithinRange(t, from, to) {
			ohlc := types.Ohlc{
				Ticker:    ticker,
				Timestamp: t,
				Volume:    0,
				Open:      point[1],
				High:      point[2],
				Low:       point[3],
				Close:     point[4],
			}
			points = append(points, ohlc)
		}
	}
	return points, err
}

func (p Provider) BatchSupported() bool {
	return false
}

func (p Provider) GetOhlcBatch(wg *sync.WaitGroup, chartChan chan *types.Chart, c context.Context, client *http.Client, tickers []string, interval string, from time.Time, to time.Time) {
	// not implemented
}
