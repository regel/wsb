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

package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	homeDir, _            = homedir.Dir()
	configSearchLocations = []string{
		".",
		filepath.Join(homeDir, ".wsb"),
		"/usr/local/etc/wsb",
		"/etc/wsb",
	}
)

type Configuration struct {
	Provider             string        `mapstructure:"provider"`
	YahooFinanceUrl      string        `mapstructure:"yahoo-finance-url"`
	YahooFinanceQueryUrl string        `mapstructure:"yahoo-finance-query-url"`
	IexCloudQueryUrl     string        `mapstructure:"iex-cloud-query-url"`
	IexCloudSecretToken  string        `mapstructure:"iex-cloud-secret-token"`
	CoingeckoQueryUrl    string        `mapstructure:"coingecko-query-url"`
	CoingeckoSecretToken string        `mapstructure:"coingecko-secret-token"`
	DialTimeout          time.Duration `mapstructure:"dial-timeout"`
	Bursts               int           `mapstructure:"bursts"`
	Tickers              []string      `mapstructure:"tickers"`
	Debug                bool          `mapstructure:"debug"`
}

func PrintDelimiterLineToWriter(w io.Writer, delimiterChar string) {
	delim := make([]string, 120)
	for i := 0; i < 120; i++ {
		delim[i] = delimiterChar
	}
	fmt.Fprintln(w, strings.Join(delim, ""))
}

func LoadConfiguration(cfgFile string, cmd *cobra.Command, printConfig bool) (*Configuration, error) {
	v := viper.New()

	cmd.Flags().VisitAll(func(flag *flag.Flag) {
		flagName := flag.Name
		if flagName != "config" && flagName != "help" {
			if err := v.BindPFlag(flagName, flag); err != nil {
				// can't really happen
				panic(fmt.Sprintln(errors.Wrapf(err, "Error binding flag '%s'", flagName)))
			}
		}
	})

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.SetEnvPrefix("WSB")

	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		v.SetConfigName("wsb")
		if cfgFile, ok := os.LookupEnv("WSB_CONFIG_DIR"); ok {
			v.AddConfigPath(cfgFile)
		} else {
			for _, searchLocation := range configSearchLocations {
				v.AddConfigPath(searchLocation)
			}
		}
	}

	if err := v.ReadInConfig(); err != nil {
		if cfgFile != "" {
			// Only error out for specified config file. Ignore for default locations.
			return nil, errors.Wrap(err, "Error loading config file")
		}
	} else {
		if printConfig {
			fmt.Fprintln(os.Stderr, "Using config file:", v.ConfigFileUsed())
		}
	}

	cfg := &Configuration{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(err, "Error unmarshaling configuration")
	}

	if printConfig {
		printCfg(cfg)
	}

	return cfg, nil
}

func printCfg(cfg *Configuration) {
	PrintDelimiterLineToWriter(os.Stderr, "-")
	fmt.Fprintln(os.Stderr, " Configuration")
	PrintDelimiterLineToWriter(os.Stderr, "-")

	e := reflect.ValueOf(cfg).Elem()
	typeOfCfg := e.Type()

	for i := 0; i < e.NumField(); i++ {
		var pattern string
		switch e.Field(i).Kind() {
		case reflect.Bool:
			pattern = "%s: %t\n"
		default:
			pattern = "%s: %s\n"
		}
		fmt.Fprintf(os.Stderr, pattern, typeOfCfg.Field(i).Name, e.Field(i).Interface())
	}

	PrintDelimiterLineToWriter(os.Stderr, "-")
}
