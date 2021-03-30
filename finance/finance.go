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
	"golang.org/x/time/rate"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	holdersUrl = "https://finance.yahoo.com/quote/%s/holders"
)

type Handler struct {
	client  *http.Client
	limiter *rate.Limiter
}

type HoldersBreakdown struct {
	Ticker                            string
	PctSharesHeldbyAllInsider         float64
	PctSharesHeldbyInstitutions       float64
	PctFloatHeldbyInstitutions        float64
	NumberofInstitutionsHoldingShares int64
}

type HoldersRow struct {
	Holder string
	Shares int64
	PctOut float64
	Value  int64
}

// Top Institutional Holders
// Top Mutual Fund Holders
type HoldersTable struct {
	Ticker string
	Rows   []HoldersRow
}

func trimPct(s string) (float64, error) {
	s2 := strings.TrimRight(s, "%")
	return strconv.ParseFloat(s2, 64)
}

func trimInt(s string) (int64, error) {
	s2 := strings.ReplaceAll(s, ",", "")
	return strconv.ParseInt(s2, 10, 64)
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
	// usage is capped at 2,000 requests/hour
	limiter := rate.NewLimiter(rate.Every(time.Hour/2000), config.Bursts)
	h := &Handler{
		client:  cli,
		limiter: limiter,
	}
	return h, nil
}

func (h *Handler) GetHolders(c context.Context, ticker string) (*HoldersBreakdown, *HoldersTable, *HoldersTable, error) {
	err := h.limiter.Wait(c)
	if err != nil {
		return nil, nil, nil, err
	}
	tables, err := ReadHtml(c, h.client, fmt.Sprintf(holdersUrl, ticker))
	if err != io.EOF {
		return nil, nil, nil, err
	}

	pctSharesHeldbyAllInsider, _ := trimPct(tables[0].Rows[0][0])
	pctSharesHeldbyInstitutions, _ := trimPct(tables[0].Rows[1][0])
	pctFloatHeldbyInstitutions, _ := trimPct(tables[0].Rows[2][0])
	numberofInstitutionsHoldingShares, _ := trimInt(tables[0].Rows[3][0])

	breakdown := &HoldersBreakdown{
		Ticker:                            ticker,
		PctSharesHeldbyAllInsider:         pctSharesHeldbyAllInsider,
		PctSharesHeldbyInstitutions:       pctSharesHeldbyInstitutions,
		PctFloatHeldbyInstitutions:        pctFloatHeldbyInstitutions,
		NumberofInstitutionsHoldingShares: numberofInstitutionsHoldingShares,
	}

	institutionalHolders := &HoldersTable{
		Ticker: ticker,
		Rows:   make([]HoldersRow, 0),
	}
	for j, row := range tables[1].Rows {
		if j == 0 {
			continue
		}
		holder := row[0]
		shares, _ := trimInt(row[1])
		pctOut, _ := trimPct(row[2])
		value, _ := trimInt(row[3])
		institutionalHolders.Rows = append(institutionalHolders.Rows,
			HoldersRow{
				Holder: holder,
				Shares: shares,
				PctOut: pctOut,
				Value:  value,
			})
	}
	fundHolders := &HoldersTable{
		Ticker: ticker,
		Rows:   make([]HoldersRow, 0),
	}
	for j, row := range tables[2].Rows {
		if j == 0 {
			continue
		}
		holder := row[0]
		shares, _ := trimInt(row[1])
		pctOut, _ := trimPct(row[2])
		value, _ := trimInt(row[3])
		fundHolders.Rows = append(fundHolders.Rows,
			HoldersRow{
				Holder: holder,
				Shares: shares,
				PctOut: pctOut,
				Value:  value,
			})
	}
	return breakdown, institutionalHolders, fundHolders, nil
}

func (h *Handler) GetOhlc(c context.Context, ticker string, interval string, startTime time.Time, endTime time.Time) ([]Ohlc, error) {
	err := h.limiter.Wait(c)
	if err != nil {
		return nil, err
	}
	points, err := ReadOhlc(c, h.client, ticker, interval, startTime, endTime)
	return points, err
}
