package finance

import (
	"context"
	"fmt"
	"github.com/regel/tinkerbell/pkg/config"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const sampleChartResponse = `
{
	"chart": {
		"result": [{
			"meta": {
				"currency": "USD",
				"symbol": "AAPL",
				"exchangeName": "NMS",
				"instrumentType": "EQUITY",
				"firstTradeDate": 345479400,
				"regularMarketTime": 1617307203,
				"gmtoffset": -14400,
				"timezone": "EDT",
				"exchangeTimezoneName": "America/New_York",
				"regularMarketPrice": 123.0,
				"chartPreviousClose": 122.15,
				"priceHint": 2,
				"currentTradingPeriod": {
					"pre": {
						"timezone": "EDT",
						"start": 1617264000,
						"end": 1617283800,
						"gmtoffset": -14400
					},
					"regular": {
						"timezone": "EDT",
						"start": 1617283800,
						"end": 1617307200,
						"gmtoffset": -14400
					},
					"post": {
						"timezone": "EDT",
						"start": 1617307200,
						"end": 1617321600,
						"gmtoffset": -14400
					}
				},
				"dataGranularity": "1d",
				"range": "",
				"validRanges": ["1d", "5d", "1mo", "3mo", "6mo", "1y", "2y", "5y", "10y", "ytd", "max"]
			},
			"timestamp": [1617307203],
			"indicators": {
				"quote": [{
					"close": [123.0],
					"high": [124.18000030517578],
					"open": [123.66000366210938],
					"volume": [75089134],
					"low": [122.48999786376953]
				}],
				"adjclose": [{
					"adjclose": [123.0]
				}]
			}
		}],
		"error": null
	}
}
`

const sampleChartNoContent = `
{
	"chart": {
		"result": [{
			"meta": {
				"currency": "USD",
				"symbol": "AAPL",
				"exchangeName": "NMS",
				"instrumentType": "EQUITY",
				"firstTradeDate": 345479400,
				"regularMarketTime": 1617307203,
				"gmtoffset": -14400,
				"timezone": "EDT",
				"exchangeTimezoneName": "America/New_York",
				"regularMarketPrice": 123.0,
				"chartPreviousClose": 122.15,
				"priceHint": 2,
				"currentTradingPeriod": {
					"pre": {
						"timezone": "EDT",
						"start": 1617264000,
						"end": 1617283800,
						"gmtoffset": -14400
					},
					"regular": {
						"timezone": "EDT",
						"start": 1617283800,
						"end": 1617307200,
						"gmtoffset": -14400
					},
					"post": {
						"timezone": "EDT",
						"start": 1617307200,
						"end": 1617321600,
						"gmtoffset": -14400
					}
				},
				"dataGranularity": "1d",
				"range": "",
				"validRanges": ["1d", "5d", "1mo", "3mo", "6mo", "1y", "2y", "5y", "10y", "ytd", "max"]
			},
			"indicators": {
				"quote": [{}],
				"adjclose": [{}]
			}
		}],
		"error": null
	}
}
`

func TestChartResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var rsp string
		if r.URL.Path == "/v8/finance/chart/AAPL" {
			rsp = sampleChartResponse
			w.Header()["Content-Type"] = []string{"application/json"}
		} else {
			panic("Cannot handle request")
		}

		fmt.Fprintln(w, rsp)
	}))
	defer ts.Close()

	context := context.Background()
	configuration := &config.Configuration{
		YahooFinanceUrl:      ts.URL,
		YahooFinanceQueryUrl: ts.URL,
		DialTimeout:          time.Second,
		Bursts:               1,
		Tickers:              []string{"AAPL"},
		Debug:                false,
	}
	n, err := NewHandler(*configuration)
	require.NoError(t, err)

	expected := Ohlc{
		Ticker:    "AAPL",
		Timestamp: time.Unix(1617348113, 0),
		Open:      123.66000366210938,
		High:      124.18000030517578,
		Low:       122.48999786376953,
		Close:     123.0,
		Volume:    75089134,
	}

	from := time.Unix(1617348113, 0)
	to := time.Unix(1617348113, 0)
	out, err := n.GetOhlc(context, "AAPL", "1d", from, to)
	require.NoError(t, err)

	require.Equal(t, 1, len(out), "Should contain one item")
	require.Equal(t, expected.Ticker, out[0].Ticker, "Ticker must be the same")
	require.Equal(t, expected.Volume, out[0].Volume, "Volume must be the same")
	require.InDelta(t, expected.Open, out[0].Open, 0.01, "Open must be the same")
	require.InDelta(t, expected.High, out[0].High, 0.01, "High must be the same")
	require.InDelta(t, expected.Low, out[0].Low, 0.01, "Low must be the same")
	require.InDelta(t, expected.Close, out[0].Close, 0.01, "Close must be the same")
}

func TestChartNoContent(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var rsp string
		if r.URL.Path == "/v8/finance/chart/AAPL" {
			rsp = sampleChartNoContent
			w.Header()["Content-Type"] = []string{"application/json"}
		} else {
			panic("Cannot handle request")
		}

		fmt.Fprintln(w, rsp)
	}))
	defer ts.Close()

	context := context.Background()
	configuration := &config.Configuration{
		YahooFinanceUrl:      ts.URL,
		YahooFinanceQueryUrl: ts.URL,
		DialTimeout:          time.Second,
		Bursts:               1,
		Tickers:              []string{"AAPL"},
		Debug:                false,
	}
	n, err := NewHandler(*configuration)
	require.NoError(t, err)

	from := time.Unix(1617283800, 0)
	to := time.Unix(1617283800, 0)
	out, err := n.GetOhlc(context, "AAPL", "1d", from, to)
	require.NoError(t, err)

	require.Empty(t, out)
}
