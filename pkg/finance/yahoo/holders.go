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
	"fmt"
	"github.com/regel/tinkerbell/pkg/common"
	"github.com/regel/tinkerbell/pkg/finance/types"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func trimPct(s string) (float64, error) {
	s2 := strings.TrimRight(s, "%")
	return strconv.ParseFloat(s2, 64)
}

func trimInt(s string) (int64, error) {
	s2 := strings.ReplaceAll(s, ",", "")
	return strconv.ParseInt(s2, 10, 64)
}

func (p Provider) GetHolders(c context.Context, client *http.Client, ticker string) (*types.HoldersBreakdown, *types.HoldersTable, *types.HoldersTable, error) {
	holdersUrl := p.YahooFinanceUrl + fmt.Sprintf("/quote/%s/holders", ticker)
	tables, err := common.ReadHtml(c, client, holdersUrl)
	if err != io.EOF {
		return nil, nil, nil, err
	}

	pctSharesHeldbyAllInsider, _ := trimPct(tables[0].Rows[0][0])
	pctSharesHeldbyInstitutions, _ := trimPct(tables[0].Rows[1][0])
	pctFloatHeldbyInstitutions, _ := trimPct(tables[0].Rows[2][0])
	numberofInstitutionsHoldingShares, _ := trimInt(tables[0].Rows[3][0])

	breakdown := &types.HoldersBreakdown{
		Ticker:                            ticker,
		PctSharesHeldbyAllInsider:         pctSharesHeldbyAllInsider,
		PctSharesHeldbyInstitutions:       pctSharesHeldbyInstitutions,
		PctFloatHeldbyInstitutions:        pctFloatHeldbyInstitutions,
		NumberofInstitutionsHoldingShares: numberofInstitutionsHoldingShares,
	}

	institutionalHolders := &types.HoldersTable{
		Ticker: ticker,
		Rows:   make([]types.HoldersRow, 0),
	}
	for j, row := range tables[1].Rows {
		if j == 0 {
			continue
		}
		holder := row[0]
		shares, _ := trimInt(row[1])
		reported, _ := time.Parse("Jan 2, 2006", row[2])
		pctOut, _ := trimPct(row[3])
		value, _ := trimInt(row[4])
		institutionalHolders.Rows = append(institutionalHolders.Rows,
			types.HoldersRow{
				Holder:       holder,
				Shares:       shares,
				DateReported: reported,
				PctOut:       pctOut,
				Value:        value,
			})
	}
	fundHolders := &types.HoldersTable{
		Ticker: ticker,
		Rows:   make([]types.HoldersRow, 0),
	}
	for j, row := range tables[2].Rows {
		if j == 0 {
			continue
		}
		holder := row[0]
		shares, _ := trimInt(row[1])
		reported, _ := time.Parse("Jan 2, 2006", row[2])
		pctOut, _ := trimPct(row[3])
		value, _ := trimInt(row[4])
		fundHolders.Rows = append(fundHolders.Rows,
			types.HoldersRow{
				Holder:       holder,
				Shares:       shares,
				DateReported: reported,
				PctOut:       pctOut,
				Value:        value,
			})
	}
	return breakdown, institutionalHolders, fundHolders, nil
}
