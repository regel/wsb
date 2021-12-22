package types

import (
	"context"
	"net/http"
	"sync"
	"time"
)

const (
	ProviderYahoo     string = "yahoo"
	ProviderIEX       string = "iex"
	ProviderCoingecko string = "coingecko"
)

type Provider interface {
	BatchSupported() bool
	GetOhlc(c context.Context, client *http.Client, ticker string, interval string, from time.Time, to time.Time) ([]Ohlc, error)
	GetOhlcBatch(wg *sync.WaitGroup, chartChan chan *Chart, c context.Context, client *http.Client, tickers []string, interval string, from time.Time, to time.Time)
	GetHolders(c context.Context, client *http.Client, ticker string) (*HoldersBreakdown, *HoldersTable, *HoldersTable, error)
}
