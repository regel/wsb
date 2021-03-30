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
	"sync"

	"context"
	"github.com/MakeNowJust/heredoc"
	"github.com/olekukonko/tablewriter"
	"github.com/regel/tinkerbell/finance"
	"github.com/regel/tinkerbell/pkg/config"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"os"
)

type holdersData struct {
	breakdown            *finance.HoldersBreakdown
	institutionalHolders *finance.HoldersTable
	fundHolders          *finance.HoldersTable
}

func newHoldersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hold",
		Short: "Prints tables of holders information to the current shell",
		Long: heredoc.Doc(`
			Query Yahoo finance holders information of selected tickers.
			Response includes:
			* Breakdown %
			* Institutional Holders names and %
			* Fund Holders names and %
			`),
		RunE: holders,
	}

	flags := cmd.Flags()
	addHoldersFlags(flags)
	return cmd
}

func addHoldersFlags(flags *flag.FlagSet) {
	addCommonFlags(flags)
}

func holders(cmd *cobra.Command, args []string) error {
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

	holdersChan := make(chan *holdersData)
	for _, ticker := range configuration.Tickers {
		wg.Add(1)
		go func(t string) {
			breakdown, institutionalHolders, fundHolders, err := handler.GetHolders(context, t)
			if err != nil {
				wg.Done()
				println(fmt.Sprintf("Error fetching '%s' data: %v", t, err))
				return
			}
			wg.Done()
			hd := &holdersData{
				breakdown:            breakdown,
				institutionalHolders: institutionalHolders,
				fundHolders:          fundHolders,
			}
			holdersChan <- hd
		}(ticker)
	}
	go func() {
		wg.Wait()
		close(holdersChan)
	}()

	PrintHolders(holdersChan)
	return nil
}

func PrintHolders(holdersChan chan *holdersData) {
	breakdownTable := tablewriter.NewWriter(os.Stdout)
	breakdownTable.SetHeader([]string{"Name",
		"% of Shares Held by All Insider",
		"% of Shares Held by Institutions",
		"% of Float Held by Institutions",
		"Number of Institutions Holding Shares",
	})
	breakdownTable.SetCaption(true, "Major Holders Breakdown.")
	for hd := range holdersChan {
		breakdownTable.Append([]string{
			hd.breakdown.Ticker,
			fmt.Sprintf("%.02f", hd.breakdown.PctSharesHeldbyAllInsider),
			fmt.Sprintf("%.02f", hd.breakdown.PctSharesHeldbyInstitutions),
			fmt.Sprintf("%.02f", hd.breakdown.PctFloatHeldbyInstitutions),
			fmt.Sprintf("%d", hd.breakdown.NumberofInstitutionsHoldingShares),
		})

		institutionalTable := tablewriter.NewWriter(os.Stdout)
		institutionalTable.SetHeader([]string{
			"Holder",
			"Shares",
			"% Out",
			"Value",
		})
		for _, row := range hd.institutionalHolders.Rows {
			institutionalTable.Append([]string{
				row.Holder,
				fmt.Sprintf("%d", row.Shares),
				fmt.Sprintf("%.02f", row.PctOut),
				fmt.Sprintf("%d", row.Value),
			})
		}
		institutionalTable.SetCaption(true, fmt.Sprintf("Top Institutional Holders '%s'.", hd.institutionalHolders.Ticker))
		institutionalTable.Render() // Send output

		fundTable := tablewriter.NewWriter(os.Stdout)
		fundTable.SetHeader([]string{
			"Holder",
			"Shares",
			"% Out",
			"Value",
		})
		for _, row := range hd.fundHolders.Rows {
			fundTable.Append([]string{
				row.Holder,
				fmt.Sprintf("%d", row.Shares),
				fmt.Sprintf("%.02f", row.PctOut),
				fmt.Sprintf("%d", row.Value),
			})
		}
		fundTable.SetCaption(true, fmt.Sprintf("Top Mutual Fund Holders '%s'.", hd.fundHolders.Ticker))
		fundTable.Render() // Send output
	}
	breakdownTable.Render() // Send output

}
