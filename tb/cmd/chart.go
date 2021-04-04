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
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"time"
)

const (
	dateFormat = "2006-01-02"
)

type chartData struct {
	Ohlc   []finance.Ohlc
	Ticker string
}

func newOhlcCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chart",
		Short: "Prints tables of price history to the current shell",
		Long: heredoc.Doc(`
			Query Yahoo finance Ohlc data points for selected tickers.
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

func addOhlcFlags(flags *flag.FlagSet) {
	addCommonFlags(flags)
	flags.String("start-time", "", heredoc.Doc(`
Start time of Ohlc time range. Format: 2006-01-02`))
	flags.String("end-time", "", heredoc.Doc(`
End time of Ohlc time range. Format: 2006-01-02`))
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
	startStr, err := cmd.Flags().GetString("start-time")
	if err != nil {
		return err
	}
	startTime, err := time.Parse(dateFormat, startStr)
	if err != nil {
		return err
	}
	endStr, err := cmd.Flags().GetString("end-time")
	if err != nil {
		return err
	}
	endTime, err := time.Parse(dateFormat, endStr)
	if err != nil {
		return err
	}
	chartChan := make(chan *chartData)
	for _, ticker := range configuration.Tickers {
		wg.Add(1)
		go func(t string, window string, start time.Time, end time.Time) {
			points, err := handler.GetOhlc(context, t, window, start, end)
			if err != nil {
				wg.Done()
				println(fmt.Sprintf("Error fetching '%s' data: %v", t, err))
				return
			}
			data := &chartData{
				Ohlc:   points,
				Ticker: t,
			}
			chartChan <- data
			wg.Done()
		}(ticker, interval, startTime, endTime)
	}
	go func() {
		wg.Wait()
		close(chartChan)
	}()

	PrintOhlc(chartChan)
	return nil
}

func PrintOhlc(chartChan chan *chartData) {
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
				row.Timestamp.Format("2006-01-02"),
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
