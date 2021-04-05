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
	"errors"
	"github.com/regel/tinkerbell/pkg/finance/types"
	"github.com/regel/tinkerbell/pkg/finance/yahoo"
	"net/http"
	"time"
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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ReadOhlc(c context.Context, client *http.Client, baseUrl string, ticker string, interval string, from time.Time, to time.Time) ([]types.Ohlc, error) {
	if !contains(getSupportedRanges(), interval) {
		return nil, errors.New("Invalid interval")
	}
	return yahoo.ReadOhlc(c, client, baseUrl, ticker, interval, from, to)
}
