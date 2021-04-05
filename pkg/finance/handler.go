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
	"fmt"
	"github.com/regel/tinkerbell/pkg/config"
	"github.com/regel/tinkerbell/pkg/finance/iex"
	"github.com/regel/tinkerbell/pkg/finance/types"
	"github.com/regel/tinkerbell/pkg/finance/yahoo"
	"golang.org/x/time/rate"
	"net"
	"net/http"
	"sync"
	"time"
)

type Handler struct {
	yahooFinanceUrl      string
	yahooFinanceQueryUrl string
	iexCloudQueryUrl     string
	iexCloudSecretToken  string

	client  *http.Client
	limiter *rate.Limiter
}

// NewHandler creates a handler
func NewHandler(config config.Configuration) (*Handler, error) {
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	var cli = &http.Client{
		Transport: netTransport,
	}
	// Yahoo Finance usage is capped at 2,000 requests/hour
	limiter := rate.NewLimiter(rate.Every(time.Hour/2000), config.Bursts)
	if config.IexCloudSecretToken != "" {
		limiter = rate.NewLimiter(rate.Every(time.Second/100), config.Bursts)
	}
	h := &Handler{
		yahooFinanceUrl:      config.YahooFinanceUrl,
		yahooFinanceQueryUrl: config.YahooFinanceQueryUrl,
		iexCloudQueryUrl:     config.IexCloudQueryUrl,
		iexCloudSecretToken:  config.IexCloudSecretToken,
		client:               cli,
		limiter:              limiter,
	}
	return h, nil
}

func (h *Handler) GetHolders(c context.Context, ticker string) (*types.HoldersBreakdown, *types.HoldersTable, *types.HoldersTable, error) {
	err := h.limiter.Wait(c)
	if err != nil {
		return nil, nil, nil, err
	}
	return yahoo.GetHolders(c, h.client, h.yahooFinanceUrl, ticker)
}

func (h *Handler) GetOhlc(c context.Context, ticker string, interval string, from time.Time, to time.Time) ([]types.Ohlc, error) {
	var points []types.Ohlc
	var err error
	err = h.limiter.Wait(c)
	if err != nil {
		return nil, err
	}
	if h.iexCloudSecretToken != "" {
		points, err = iex.GetOhlc(c, h.client, h.iexCloudQueryUrl, h.iexCloudSecretToken, ticker, interval, from, to)
	} else {
		points, err = yahoo.GetOhlc(c, h.client, h.yahooFinanceQueryUrl, ticker, interval, from, to)
	}
	return points, err
}

func (h *Handler) GetOhlcBatch(c context.Context, wg *sync.WaitGroup, chartChan chan *types.Chart, tickers []string, interval string, from time.Time, to time.Time) {
	if h.iexCloudSecretToken != "" {
		iex.GetOhlcBatch(wg, chartChan, c, h.client, h.iexCloudQueryUrl, h.iexCloudSecretToken, tickers, interval, from, to)
		return
	}
	for _, ticker := range tickers {
		wg.Add(1)
		go func(t string, window string, from time.Time, to time.Time) {
			points, err := h.GetOhlc(c, t, window, from, to)
			if err != nil {
				wg.Done()
				println(fmt.Sprintf("Error fetching '%s' data: %v", t, err))
				return
			}
			chart := &types.Chart{
				Ohlc:   points,
				Ticker: t,
			}
			chartChan <- chart
			wg.Done()
		}(ticker, interval, from, to)
	}
}
