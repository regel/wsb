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
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	defaultYahooBaseUrl  = "https://finance.yahoo.com"
	defaultYahooQueryUrl = "https://query2.finance.yahoo.com"
)

var (
	cfgFile string
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tb",
		Short: "The Go client to get stock market data",
		Long: heredoc.Doc(`
			Get finance data
			* Price history
			* holder information
			for the given ticker names.`),
		SilenceUsage: true,
	}

	cmd.AddCommand(newHoldersCmd())
	cmd.AddCommand(newOhlcCmd())
	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(newGenerateDocsCmd())

	return cmd
}

// Execute runs the application
func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func addCommonFlags(flags *pflag.FlagSet) {
	flags.StringVar(&cfgFile, "config", "", "Config file")
	flags.String("yahoo-finance-url", defaultYahooBaseUrl, heredoc.Doc(`
		Yahoo Finance Base Url`))
	flags.String("yahoo-finance-query-url", defaultYahooQueryUrl, heredoc.Doc(`
		Yahoo Finance Query Url`))
	flags.StringSlice("tickers", []string{}, heredoc.Doc(`
		Names of selected tickers`))
	flags.Bool("print-config", false, heredoc.Doc(`
		Prints the configuration to stderr`))
	flags.Bool("debug", false, heredoc.Doc(`
		Print API calls to external tools to stdout`))
	flags.Duration("dial-timeout", 5*time.Second, heredoc.Doc(`
		Dial timeout to connect to external API sources`))
	flags.Int("bursts", 1, heredoc.Doc(`
		Permits bursts of at most N concurrent API calls`))

}
