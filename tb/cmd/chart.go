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

package cmd

import (
	"fmt"
	"os"
	"sync"

	"context"
	"github.com/MakeNowJust/heredoc"
	"github.com/olekukonko/tablewriter"
	"github.com/regel/tinkerbell/pkg/config"
	"github.com/regel/tinkerbell/pkg/finance"
	"github.com/regel/tinkerbell/pkg/finance/types"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"time"
)

const (
	dateFormat     = "2006-01-02"
	dateFormatLong = "2006-01-02T15:04:05"
)

func newOhlcCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chart",
		Short: "Prints tables of stock price history (OHLC) to the current shell",
		Long: heredoc.Doc(`
			Query finance Chart data points for selected tickers.
			An open-high-low-close chart (also OHLC) is a type of chart typically
			used to illustrate movements in the price of a financial instrument over time.

			Response includes:
			* Prices: Open, High, Low, Close
			* Volume
			`),
		RunE: chart,
	}

	flags := cmd.Flags()
	addOhlcFlags(flags)
	return cmd
}

func parseDate(date string) (time.Time, error) {
	dt, err := time.Parse(dateFormat, date)
	if err == nil {
		return dt, nil
	}
	return time.Parse(dateFormatLong, date)
}

func formatDate(t time.Time) string {
	h := t.Hour()
	m := t.Minute()
	s := t.Second()
	n := t.Nanosecond()
	if h == 0 && m == 0 && s == 0 && n == 0 {
		return t.In(time.UTC).Format(dateFormat)
	}
	return t.In(time.UTC).Format(dateFormatLong)
}

func addOhlcFlags(flags *flag.FlagSet) {
	addCommonFlags(flags)
	now := time.Now()
	minus7d := now.AddDate(0, 0, -7)
	flags.String("from", minus7d.Format(dateFormat), heredoc.Doc(`
Start time of Ohlc time range. Format: 2006-01-02, or 2006-01-02T15:04:05`))
	flags.String("to", now.Format(dateFormat), heredoc.Doc(`
End time of Ohlc time range. Format: 2006-01-02 or 2006-01-02T15:04:05`))
	flags.String("interval", "1d", heredoc.Doc(`
                Time interval range. Supported values: (1d, 5d, 1mo, 3mo, 6mo, 1y, 2y, 5y, 10y, ytd, max)`))
}

func chart(cmd *cobra.Command, args []string) error {
	var wg sync.WaitGroup
	printConfig, err := cmd.Flags().GetBool("print-config")
	if err != nil {
		return err
	}

	configuration, err := config.LoadConfiguration(cfgFile, cmd, printConfig)
	if err != nil {
		return fmt.Errorf("Error loading configuration: %s", err)
	}

	context := context.Background()
	handler, err := finance.NewHandler(*configuration)
	if err != nil {
		return fmt.Errorf("Error creating handler: %s", err)
	}
	interval, err := cmd.Flags().GetString("interval")
	if err != nil {
		return err
	}
	fromStr, err := cmd.Flags().GetString("from")
	if err != nil {
		return err
	}
	from, err := parseDate(fromStr)
	if err != nil {
		return err
	}
	toStr, err := cmd.Flags().GetString("to")
	if err != nil {
		return err
	}
	to, err := parseDate(toStr)
	if err != nil {
		return err
	}
	chartChan := make(chan *types.Chart)
	handler.GetOhlcBatch(context, &wg, chartChan, configuration.Tickers, interval, from, to)
	go func() {
		wg.Wait()
		close(chartChan)
	}()

	PrintOhlc(chartChan)
	return nil
}

func PrintOhlc(chartChan chan *types.Chart) {
	for data := range chartChan {
		history := tablewriter.NewWriter(os.Stdout)
		history.SetHeader([]string{
			"Date",
			"Open",
			"High",
			"Low",
			"Close",
			"Volume",
		})
		for _, row := range data.Ohlc {
			history.Append([]string{
				formatDate(row.Timestamp),
				fmt.Sprintf("%.02f", row.Open),
				fmt.Sprintf("%.02f", row.High),
				fmt.Sprintf("%.02f", row.Low),
				fmt.Sprintf("%.02f", row.Close),
				fmt.Sprintf("%d", row.Volume),
			})
		}
		history.SetCaption(true, fmt.Sprintf("History of '%s'.", data.Ticker))
		history.Render() // Send output
	}
}
