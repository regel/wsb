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
	"github.com/regel/tinkerbell/pkg/config"
	"github.com/regel/tinkerbell/pkg/finance/iex"
	"github.com/regel/tinkerbell/pkg/finance/types"
	"github.com/regel/tinkerbell/pkg/finance/yahoo"
	"golang.org/x/time/rate"
	"net"
	"net/http"
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

func (h *Handler) GetOhlc(c context.Context, ticker string, interval string, startTime time.Time, endTime time.Time) ([]types.Ohlc, error) {
	var points []types.Ohlc
	var err error
	err = h.limiter.Wait(c)
	if err != nil {
		return nil, err
	}
	if h.iexCloudSecretToken != "" {
		points, err = iex.ReadOhlc(c, h.client, h.iexCloudQueryUrl, h.iexCloudSecretToken, ticker, interval, startTime, endTime)
	} else {
		points, err = yahoo.ReadOhlc(c, h.client, h.yahooFinanceQueryUrl, ticker, interval, startTime, endTime)
	}
	return points, err
}
