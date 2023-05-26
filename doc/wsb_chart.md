## wsb chart

Prints tables of stock price history (OHLC) to the current shell

### Synopsis

Query finance Chart data points for selected tickers.
An open-high-low-close chart (also OHLC) is a type of chart typically
used to illustrate movements in the price of a financial instrument over time.

Response includes:
* Prices: Open, High, Low, Close
* Volume


```
wsb chart [flags]
```

### Options

```
      --bursts int                       Permits bursts of at most N concurrent API calls (default 1)
      --coingecko-query-url string       The Most Comprehensive Cryptocurrency API (default "https://api.coingecko.com")
      --coingecko-secret-token string    Secret token to enable access to the Paid API
      --config string                    Config file
      --debug                            Print API calls to external tools to stdout
      --dial-timeout duration            Dial timeout to connect to external API sources (default 5s)
      --from string                      Start time of Ohlc time range. Format: 2006-01-02, or 2006-01-02T15:04:05 (default "2021-12-15")
  -h, --help                             help for chart
      --iex-cloud-query-url string       IEX Cloud is a platform that makes financial data and services accessible to everyone (default "https://cloud.iexapis.com")
      --iex-cloud-secret-token string    Secret token to enable access to IEX Cloud API
      --interval string                  Time interval range. Supported values: (1d, 5d, 1mo, 3mo, 6mo, 1y, 2y, 5y, 10y, ytd, max) (default "1d")
      --print-config                     Prints the configuration to stderr
      --provider string                  Provider of market data. Supported providers: 'yahoo' (default), 'iex', 'coingecko' (default "yahoo")
      --tickers strings                  Names of selected tickers
      --to string                        End time of Ohlc time range. Format: 2006-01-02 or 2006-01-02T15:04:05 (default "2021-12-22")
      --yahoo-finance-query-url string   Yahoo Finance Query Url (default "https://query2.finance.yahoo.com")
      --yahoo-finance-url string         Yahoo Finance Base Url (default "https://finance.yahoo.com")
```

### SEE ALSO

* [wsb](wsb.md)	 - The Go client to get stock market and cryptocurrencies market data

###### Auto generated by spf13/cobra on 22-Dec-2021