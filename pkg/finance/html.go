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
	"golang.org/x/net/html"
	"log"
	"net/http"
	"strings"
	"time"
)

type Table struct {
	Rows [][]string
}

func ReadHtml(c context.Context, client *http.Client, url string) ([]*Table, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
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

	z := html.NewTokenizer(res.Body)
	tables := make([]*Table, 0)
	var table *Table
	var openTag string
	var raw string
	var row int
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return tables, z.Err()
		case html.TextToken:
			if openTag != "" {
				raw = raw + string(z.Text())
			}
		case html.StartTagToken, html.EndTagToken:
			tn, _ := z.TagName()
			tag := string(tn)
			if strings.EqualFold(tag, "td") || strings.EqualFold(tag, "th") {
				if tt == html.StartTagToken {
					openTag = tag
					raw = ""
				} else {
					openTag = ""
					trimmed := strings.Trim(raw, " \n")
					table.Rows[row] = append(table.Rows[row], trimmed)
				}
			}
			if tt == html.EndTagToken && strings.EqualFold(tag, "table") {
				tables = append(tables, table)
			}
			if tt == html.StartTagToken && strings.EqualFold(tag, "table") {
				table = &Table{
					Rows: make([][]string, 0),
				}
				row = 0
			}
			if tt == html.StartTagToken && strings.EqualFold(tag, "tr") {
				empty := make([]string, 0)
				table.Rows = append(table.Rows, empty)
			}
			if tt == html.EndTagToken && strings.EqualFold(tag, "tr") {
				row = row + 1
			}
		}
	}
}
