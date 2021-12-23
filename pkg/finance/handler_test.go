package finance

import (
	"context"
	"fmt"
	"github.com/regel/tinkerbell/pkg/config"
	"github.com/regel/tinkerbell/pkg/finance/types"
	"net/http"
	"net/http/httptest"
	"sync"
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

const sampleHoldersResponse = `
<html>
  <table>
    <tbody data-reactid="24">
      <tr>
        <td>22.29%</td>
        <td>
          <span data-reactid="28">% of Shares Held by All Insider</span>
        </td>
      </tr>
      <tr>
        <td>110.64%</td>
        <td>
          <span data-reactid="32">% of Shares Held by Institutions</span>
        </td>
      </tr>
      <tr>
        <td>142.38%</td>
        <td>
          <span data-reactid="36">% of Float Held by Institutions</span>
        </td>
      </tr>
      <tr>
        <td>296</td>
        <td>
          <span data-reactid="40">Number of Institutions Holding Shares</span>
        </td>
      </tr>
    </tbody>
  </table>
<h3>
  <span data-reactid="43">Top Institutional Holders</span>
</h3>
<table>
  <thead data-reactid="45">
    <tr>
      <th>
        <span data-reactid="48">Holder</span>
      </th>
      <th>
        <span data-reactid="50">Shares</span>
      </th>
      <th>
        <span data-reactid="52">Date Reported</span>
      </th>
      <th>
        <span data-reactid="54">% Out</span>
      </th>
      <th>
        <span data-reactid="56">Value</span>
      </th>
    </tr>
  </thead>
  <tbody data-reactid="57">
    <tr>
      <td>FMR, LLC</td>
      <td>9,276,087</td>
      <td>
        <span data-reactid="62">Dec 30, 2020</span>
      </td>
      <td>13.26%</td>
      <td>174,761,479</td>
    </tr>
    <tr>
      <td>Blackrock Inc.</td>
      <td>9,217,335</td>
      <td>
        <span data-reactid="69">Dec 30, 2020</span>
      </td>
      <td>13.18%</td>
      <td>173,654,591</td>
    </tr>
    <tr>
      <td>Vanguard Group, Inc. (The)</td>
      <td>5,162,095</td>
      <td>
        <span data-reactid="76">Dec 30, 2020</span>
      </td>
      <td>7.38%</td>
      <td>97,253,869</td>
    </tr>
    <tr>
      <td>Senvest Management LLC</td>
      <td>5,050,915</td>
      <td>
        <span data-reactid="83">Dec 30, 2020</span>
      </td>
      <td>7.22%</td>
      <td>95,159,238</td>
    </tr>
    <tr>
      <td>Maverick Capital Ltd.</td>
      <td>4,658,607</td>
      <td>
        <span data-reactid="90">Dec 30, 2020</span>
      </td>
      <td>6.66%</td>
      <td>87,768,155</td>
    </tr>
    <tr>
      <td>Morgan Stanley</td>
      <td>4,275,838</td>
      <td>
        <span data-reactid="97">Dec 30, 2020</span>
      </td>
      <td>6.11%</td>
      <td>80,556,787</td>
    </tr>
    <tr>
      <td>Dimensional Fund Advisors LP</td>
      <td>3,934,919</td>
      <td>
        <span data-reactid="104">Dec 30, 2020</span>
      </td>
      <td>5.63%</td>
      <td>74,133,873</td>
    </tr>
    <tr>
      <td>Shaw D.E. &amp; Co., Inc.</td>
      <td>2,841,563</td>
      <td>
        <span data-reactid="111">Dec 30, 2020</span>
      </td>
      <td>4.06%</td>
      <td>53,535,046</td>
    </tr>
    <tr>
      <td>Susquehanna International Group, LLP</td>
      <td>2,487,366</td>
      <td>
        <span data-reactid="118">Dec 30, 2020</span>
      </td>
      <td>3.56%</td>
      <td>46,861,975</td>
    </tr>
    <tr>
      <td>State Street Corporation</td>
      <td>2,445,216</td>
      <td>
        <span data-reactid="125">Dec 30, 2020</span>
      </td>
      <td>3.50%</td>
      <td>46,067,869</td>
    </tr>
  </tbody>
</table>
<h3>
  <span data-reactid="130">Top Mutual Fund Holders</span>
</h3>
<table>
  <thead data-reactid="132">
    <tr>
      <th>
        <span data-reactid="135">Holder</span>
      </th>
      <th>
        <span data-reactid="137">Shares</span>
      </th>
      <th>
        <span data-reactid="139">Date Reported</span>
      </th>
      <th>
        <span data-reactid="141">% Out</span>
      </th>
      <th>
        <span data-reactid="143">Value</span>
      </th>
    </tr>
  </thead>
  <tbody data-reactid="144">
    <tr>
      <td>iShares Core S&amp;P Smallcap ETF</td>
      <td>3,665,529</td>
      <td>
        <span data-reactid="149">Feb 27, 2021</span>
      </td>
      <td>5.24%</td>
      <td>372,930,920</td>
    </tr>
    <tr>
      <td>Vanguard Total Stock Market Index Fund</td>
      <td>1,468,071</td>
      <td>
        <span data-reactid="156">Dec 30, 2020</span>
      </td>
      <td>2.10%</td>
      <td>27,658,457</td>
    </tr>
    <tr>
      <td>Morgan Stanley Inst Fund Inc-Inception Port</td>
      <td>1,415,967</td>
      <td>
        <span data-reactid="163">Dec 30, 2020</span>
      </td>
      <td>2.02%</td>
      <td>26,676,818</td>
    </tr>
    <tr>
      <td>iShares Russell 2000 ETF</td>
      <td>1,395,217</td>
      <td>
        <span data-reactid="170">Feb 27, 2021</span>
      </td>
      <td>1.99%</td>
      <td>141,949,377</td>
    </tr>
    <tr>
      <td>Vanguard Extended Market Index Fund</td>
      <td>817,583</td>
      <td>
        <span data-reactid="177">Dec 30, 2020</span>
      </td>
      <td>1.17%</td>
      <td>15,403,263</td>
    </tr>
    <tr>
      <td>Vanguard Small-Cap Index Fund</td>
      <td>633,843</td>
      <td>
        <span data-reactid="184">Dec 30, 2020</span>
      </td>
      <td>0.91%</td>
      <td>11,941,602</td>
    </tr>
    <tr>
      <td>iShares Russell 2000 Value ETF</td>
      <td>612,531</td>
      <td>
        <span data-reactid="191">Feb 27, 2021</span>
      </td>
      <td>0.88%</td>
      <td>62,318,903</td>
    </tr>
    <tr>
      <td>EQ Advisors Trust-EQ/Morgan Stanley Small Cap Growth Port</td>
      <td>602,280</td>
      <td>
        <span data-reactid="198">Dec 30, 2020</span>
      </td>
      <td>0.86%</td>
      <td>11,346,955</td>
    </tr>
    <tr>
      <td>Vanguard Strategic Equity Fund</td>
      <td>571,438</td>
      <td>
        <span data-reactid="205">Dec 30, 2020</span>
      </td>
      <td>0.82%</td>
      <td>10,765,891</td>
    </tr>
    <tr>
      <td>SPDR (R) Ser Tr-SPDR (R) S&amp;P (R) Retail ETF</td>
      <td>523,000</td>
      <td>
        <span data-reactid="212">Feb 27, 2021</span>
      </td>
      <td>0.75%</td>
      <td>53,210,020</td>
    </tr>
  </tbody>
</table>
</html>
`

const sampleIexChartResponse = `
{
    "AAPL": {
        "chart": [
            {
                "close": 120.13,
                "high": 123.6,
                "low": 118.62,
                "open": 121.75,
                "symbol": "AAPL",
                "volume": 178154975,
                "id": "HISTORICAL_PRICES",
                "key": "AAPL",
                "subkey": "",
                "date": "2021-03-04",
                "updated": 1614909622000,
                "changeOverTime": 0,
                "marketChangeOverTime": 0,
                "uOpen": 121.75,
                "uClose": 120.13,
                "uHigh": 123.6,
                "uLow": 118.62,
                "uVolume": 178154975,
                "fOpen": 121.75,
                "fClose": 120.13,
                "fHigh": 123.6,
                "fLow": 118.62,
                "fVolume": 178154975,
                "label": "Mar 4, 21",
                "change": 0,
                "changePercent": 0
            }
	]
    }
}
`

const sampleIexChartResponseNoContent = `
{
    "AAPL": {
        "chart": [
	]
    },
    "GME": {
        "chart": [
	]
    }
}
`

const sampleCoingeckoChartResponse = `
[
    [
        1614816000000,
        121.75,
        123.6,
        118.62,
        120.13
    ]
]
`

const sampleCoingeckoChartErrorResponse = `
{"status":500,"error":"Internal Server Error"}
`

func strftime(s string) time.Time {
	tm, _ := time.Parse("2006-01-02", s)
	return tm
}

func TestYahooHoldersResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var rsp string
		if r.URL.Path == "/quote/GME/holders" {
			rsp = sampleHoldersResponse
			w.Header()["Content-Type"] = []string{"text/html"}
		} else {
			panic("Cannot handle request")
		}

		fmt.Fprintln(w, rsp)
	}))
	defer ts.Close()

	context := context.Background()
	configuration := &config.Configuration{
		Provider:             "yahoo",
		YahooFinanceUrl:      ts.URL,
		YahooFinanceQueryUrl: ts.URL,
		DialTimeout:          time.Second,
		Bursts:               1,
		Tickers:              []string{"GME"},
		Debug:                false,
	}
	n, err := NewHandler(*configuration)
	require.NoError(t, err)

	brk, it, ft, err := n.GetHolders(context, "GME")
	require.NoError(t, err)
	require.NotNil(t, brk)
	require.NotNil(t, it)
	require.NotNil(t, ft)

	require.InDelta(t, 22.29, brk.PctSharesHeldbyAllInsider, 0.01, "% of Shares Held by All Insider must be the same")
	require.InDelta(t, 110.64, brk.PctSharesHeldbyInstitutions, 0.01, "% of Shares Held by Institutions must be the same")
	require.InDelta(t, 142.38, brk.PctFloatHeldbyInstitutions, 0.01, "% of Float Held by Institutions must be the same")
	require.EqualValues(t, 296, brk.NumberofInstitutionsHoldingShares, "Number of Institutions Holding Shares must be the same")

	var expected *types.HoldersTable
	expected = &types.HoldersTable{
		Ticker: "GME",
		Rows: []types.HoldersRow{
			types.HoldersRow{
				Holder:       "FMR, LLC",
				Shares:       9276087,
				DateReported: strftime("2020-12-30"),
				PctOut:       13.26,
				Value:        174761479,
			},
			types.HoldersRow{
				Holder:       "Blackrock Inc.",
				Shares:       9217335,
				DateReported: strftime("2020-12-30"),
				PctOut:       13.18,
				Value:        173654591,
			},
			types.HoldersRow{
				Holder:       "Vanguard Group, Inc. (The)",
				Shares:       5162095,
				DateReported: strftime("2020-12-30"),
				PctOut:       7.38,
				Value:        97253869,
			},
			types.HoldersRow{
				Holder:       "Senvest Management LLC",
				Shares:       5050915,
				DateReported: strftime("2020-12-30"),
				PctOut:       7.22,
				Value:        95159238,
			},
			types.HoldersRow{
				Holder:       "Maverick Capital Ltd.",
				Shares:       4658607,
				DateReported: strftime("2020-12-30"),
				PctOut:       6.66,
				Value:        87768155,
			},
			types.HoldersRow{
				Holder:       "Morgan Stanley",
				Shares:       4275838,
				DateReported: strftime("2020-12-30"),
				PctOut:       6.11,
				Value:        80556787,
			},
			types.HoldersRow{
				Holder:       "Dimensional Fund Advisors LP",
				Shares:       3934919,
				DateReported: strftime("2020-12-30"),
				PctOut:       5.63,
				Value:        74133873,
			},
			types.HoldersRow{
				Holder:       "Shaw D.E. & Co., Inc.",
				Shares:       2841563,
				DateReported: strftime("2020-12-30"),
				PctOut:       4.06,
				Value:        53535046,
			},
			types.HoldersRow{
				Holder:       "Susquehanna International Group, LLP",
				Shares:       2487366,
				DateReported: strftime("2020-12-30"),
				PctOut:       3.56,
				Value:        46861975,
			},
			types.HoldersRow{
				Holder:       "State Street Corporation",
				Shares:       2445216,
				DateReported: strftime("2020-12-30"),
				PctOut:       3.5,
				Value:        46067869,
			},
		},
	}
	require.EqualValues(t, expected, it, "Institutional Holders must be the same")

	expected = &types.HoldersTable{
		Ticker: "GME",
		Rows: []types.HoldersRow{
			types.HoldersRow{
				Holder:       "iShares Core S&P Smallcap ETF",
				Shares:       3665529,
				DateReported: strftime("2021-02-27"),
				PctOut:       5.24,
				Value:        372930920,
			},
			types.HoldersRow{
				Holder:       "Vanguard Total Stock Market Index Fund",
				Shares:       1468071,
				DateReported: strftime("2020-12-30"),
				PctOut:       2.1,
				Value:        27658457,
			},
			types.HoldersRow{
				Holder:       "Morgan Stanley Inst Fund Inc-Inception Port",
				Shares:       1415967,
				DateReported: strftime("2020-12-30"),
				PctOut:       2.02,
				Value:        26676818,
			},
			types.HoldersRow{
				Holder:       "iShares Russell 2000 ETF",
				Shares:       1395217,
				DateReported: strftime("2021-02-27"),
				PctOut:       1.99,
				Value:        141949377,
			},
			types.HoldersRow{
				Holder:       "Vanguard Extended Market Index Fund",
				Shares:       817583,
				DateReported: strftime("2020-12-30"),
				PctOut:       1.17,
				Value:        15403263,
			},
			types.HoldersRow{
				Holder:       "Vanguard Small-Cap Index Fund",
				Shares:       633843,
				DateReported: strftime("2020-12-30"),
				PctOut:       0.91,
				Value:        11941602,
			},
			types.HoldersRow{
				Holder:       "iShares Russell 2000 Value ETF",
				Shares:       612531,
				DateReported: strftime("2021-02-27"),
				PctOut:       0.88,
				Value:        62318903,
			},
			types.HoldersRow{
				Holder:       "EQ Advisors Trust-EQ/Morgan Stanley Small Cap Growth Port",
				Shares:       602280,
				DateReported: strftime("2020-12-30"),
				PctOut:       0.86,
				Value:        11346955,
			},
			types.HoldersRow{
				Holder:       "Vanguard Strategic Equity Fund",
				Shares:       571438,
				DateReported: strftime("2020-12-30"),
				PctOut:       0.82,
				Value:        10765891,
			},
			types.HoldersRow{
				Holder:       "SPDR (R) Ser Tr-SPDR (R) S&P (R) Retail ETF",
				Shares:       523000,
				DateReported: strftime("2021-02-27"),
				PctOut:       0.75,
				Value:        53210020,
			},
		},
	}
	require.EqualValues(t, expected, ft, "Fund Holders must be the same")
}

func TestYahooChartResponse(t *testing.T) {
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
		Provider:             "yahoo",
		YahooFinanceUrl:      ts.URL,
		YahooFinanceQueryUrl: ts.URL,
		DialTimeout:          time.Second,
		Bursts:               1,
		Tickers:              []string{"AAPL"},
		Debug:                false,
	}
	n, err := NewHandler(*configuration)
	require.NoError(t, err)

	loc, _ := time.LoadLocation("America/New_York")
	expected := types.Ohlc{
		Ticker:    "AAPL",
		Timestamp: time.Unix(1617307203, 0).In(loc),
		Open:      123.66000366210938,
		High:      124.18000030517578,
		Low:       122.48999786376953,
		Close:     123.0,
		Volume:    75089134,
	}

	from := time.Unix(1617307203, 0)
	to := time.Unix(1617307203, 0)
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

func TestYahooChartBatchResponse(t *testing.T) {
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
		Provider:             "yahoo",
		YahooFinanceUrl:      ts.URL,
		YahooFinanceQueryUrl: ts.URL,
		DialTimeout:          time.Second,
		Bursts:               1,
		Tickers:              []string{"AAPL"},
		Debug:                false,
	}
	n, err := NewHandler(*configuration)
	require.NoError(t, err)

	loc, _ := time.LoadLocation("America/New_York")
	expected := types.Ohlc{
		Ticker:    "AAPL",
		Timestamp: time.Unix(1617307203, 0).In(loc),
		Open:      123.66000366210938,
		High:      124.18000030517578,
		Low:       122.48999786376953,
		Close:     123.0,
		Volume:    75089134,
	}

	from := time.Unix(1617307203, 0)
	to := time.Unix(1617307203, 0)
	var wg sync.WaitGroup
	tickers := []string{"AAPL"}
	chartChan := make(chan *types.Chart)
	n.GetOhlcBatch(context, &wg, chartChan, tickers, "1d", from, to)
	go func() {
		wg.Wait()
		close(chartChan)
	}()

	out := <-chartChan
	// Make sure that the function does close the channel
	_, ok := (<-chartChan)

	// If we can receive on the channel then it is NOT closed
	if ok {
		t.Error("Channel is not closed")
	}
	require.Equal(t, 1, len(out.Ohlc), "Should contain one item")
	require.Equal(t, expected.Ticker, out.Ohlc[0].Ticker, "Ticker must be the same")
	require.Equal(t, expected.Volume, out.Ohlc[0].Volume, "Volume must be the same")
	require.InDelta(t, expected.Open, out.Ohlc[0].Open, 0.01, "Open must be the same")
	require.InDelta(t, expected.High, out.Ohlc[0].High, 0.01, "High must be the same")
	require.InDelta(t, expected.Low, out.Ohlc[0].Low, 0.01, "Low must be the same")
	require.InDelta(t, expected.Close, out.Ohlc[0].Close, 0.01, "Close must be the same")
}

func TestYahooChartNoContent(t *testing.T) {
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
		Provider:             "yahoo",
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

func TestIexCloudChartResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var rsp string
		if r.URL.Path == "/v1/stock/market/batch" {
			rsp = sampleIexChartResponse
			w.Header()["Content-Type"] = []string{"application/json"}
		} else {
			panic("Cannot handle request")
		}

		fmt.Fprintln(w, rsp)
	}))
	defer ts.Close()

	context := context.Background()
	configuration := &config.Configuration{
		Provider:            "iex",
		IexCloudQueryUrl:    ts.URL,
		IexCloudSecretToken: "SECRET_TOKEN",
		DialTimeout:         time.Second,
		Bursts:              1,
		Tickers:             []string{"AAPL"},
		Debug:               false,
	}
	n, err := NewHandler(*configuration)
	require.NoError(t, err)

	tm, _ := time.Parse("2006-01-02", "2021-03-04")
	expected := types.Ohlc{
		Ticker:    "AAPL",
		Timestamp: tm,
		Open:      121.75,
		High:      123.6,
		Low:       118.62,
		Close:     120.13,
		Volume:    178154975,
	}

	out, err := n.GetOhlc(context, "AAPL", "1d", tm, tm)
	require.NoError(t, err)

	require.Equal(t, 1, len(out), "Should contain one item")
	require.Equal(t, expected.Ticker, out[0].Ticker, "Ticker must be the same")
	require.Equal(t, expected.Volume, out[0].Volume, "Volume must be the same")
	require.InDelta(t, expected.Open, out[0].Open, 0.01, "Open must be the same")
	require.InDelta(t, expected.High, out[0].High, 0.01, "High must be the same")
	require.InDelta(t, expected.Low, out[0].Low, 0.01, "Low must be the same")
	require.InDelta(t, expected.Close, out[0].Close, 0.01, "Close must be the same")
}

func TestIexCloudChartBatchResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var rsp string
		if r.URL.Path == "/v1/stock/market/batch" {
			rsp = sampleIexChartResponse
			w.Header()["Content-Type"] = []string{"application/json"}
		} else {
			panic("Cannot handle request")
		}

		fmt.Fprintln(w, rsp)
	}))
	defer ts.Close()

	context := context.Background()
	configuration := &config.Configuration{
		Provider:            "iex",
		IexCloudQueryUrl:    ts.URL,
		IexCloudSecretToken: "SECRET_TOKEN",
		DialTimeout:         time.Second,
		Bursts:              1,
		Tickers:             []string{"AAPL"},
		Debug:               false,
	}
	n, err := NewHandler(*configuration)
	require.NoError(t, err)

	tm, _ := time.Parse("2006-01-02", "2021-03-04")
	expected := types.Ohlc{
		Ticker:    "AAPL",
		Timestamp: tm,
		Open:      121.75,
		High:      123.6,
		Low:       118.62,
		Close:     120.13,
		Volume:    178154975,
	}

	var wg sync.WaitGroup
	tickers := []string{"AAPL"}
	chartChan := make(chan *types.Chart)
	n.GetOhlcBatch(context, &wg, chartChan, tickers, "1d", tm, tm)
	go func() {
		wg.Wait()
		close(chartChan)
	}()

	out := <-chartChan
	// Make sure that the function does close the channel
	_, ok := (<-chartChan)

	// If we can receive on the channel then it is NOT closed
	if ok {
		t.Error("Channel is not closed")
	}
	require.Equal(t, 1, len(out.Ohlc), "Should contain one item")
	require.Equal(t, expected.Ticker, out.Ohlc[0].Ticker, "Ticker must be the same")
	require.Equal(t, expected.Volume, out.Ohlc[0].Volume, "Volume must be the same")
	require.InDelta(t, expected.Open, out.Ohlc[0].Open, 0.01, "Open must be the same")
	require.InDelta(t, expected.High, out.Ohlc[0].High, 0.01, "High must be the same")
	require.InDelta(t, expected.Low, out.Ohlc[0].Low, 0.01, "Low must be the same")
	require.InDelta(t, expected.Close, out.Ohlc[0].Close, 0.01, "Close must be the same")
}

func TestIexCloudChartNoContent(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var rsp string
		if r.URL.Path == "/v1/stock/market/batch" {
			rsp = sampleIexChartResponseNoContent
			w.Header()["Content-Type"] = []string{"application/json"}
		} else {
			panic("Cannot handle request")
		}

		fmt.Fprintln(w, rsp)
	}))
	defer ts.Close()

	context := context.Background()
	configuration := &config.Configuration{
		Provider:            "iex",
		IexCloudQueryUrl:    ts.URL,
		IexCloudSecretToken: "SECRET_TOKEN",
		DialTimeout:         time.Second,
		Bursts:              1,
		Tickers:             []string{"AAPL"},
		Debug:               false,
	}
	n, err := NewHandler(*configuration)
	require.NoError(t, err)

	from := time.Unix(1617283800, 0)
	to := time.Unix(1617283800, 0)
	out, err := n.GetOhlc(context, "AAPL", "1d", from, to)
	require.NoError(t, err)

	require.Empty(t, out)
}

func TestCoinGeckoChartResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var rsp string
		if r.URL.Path == "/api/v3/coins/bitcoin/ohlc" {
			rsp = sampleCoingeckoChartResponse
			w.Header()["Content-Type"] = []string{"application/json"}
		} else {
			panic("Cannot handle request")
		}

		fmt.Fprintln(w, rsp)
	}))
	defer ts.Close()

	context := context.Background()
	configuration := &config.Configuration{
		Provider:          "coingecko",
		CoingeckoQueryUrl: ts.URL,
		DialTimeout:       time.Second,
		Bursts:            1,
		Tickers:           []string{"bitcoin"},
		Debug:             false,
	}
	n, err := NewHandler(*configuration)
	require.NoError(t, err)

	tm, _ := time.Parse("2006-01-02", "2021-03-04")
	expected := types.Ohlc{
		Ticker:    "bitcoin",
		Timestamp: tm,
		Open:      121.75,
		High:      123.6,
		Low:       118.62,
		Close:     120.13,
		Volume:    0,
	}

	out, err := n.GetOhlc(context, "bitcoin", "1d", tm, tm)
	require.NoError(t, err)

	require.Equal(t, 1, len(out), "Should contain one item")
	require.Equal(t, expected.Ticker, out[0].Ticker, "Ticker must be the same")
	require.Equal(t, expected.Volume, out[0].Volume, "Volume must be the same")
	require.InDelta(t, expected.Open, out[0].Open, 0.01, "Open must be the same")
	require.InDelta(t, expected.High, out[0].High, 0.01, "High must be the same")
	require.InDelta(t, expected.Low, out[0].Low, 0.01, "Low must be the same")
	require.InDelta(t, expected.Close, out[0].Close, 0.01, "Close must be the same")
}

func TestCoinGeckoChartErrorResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var rsp string
		if r.URL.Path == "/api/v3/coins/xxx/ohlc" {
			rsp = sampleCoingeckoChartErrorResponse
			w.Header()["Content-Type"] = []string{"application/json"}
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			panic("Cannot handle request")
		}

		fmt.Fprintln(w, rsp)
	}))
	defer ts.Close()

	context := context.Background()
	configuration := &config.Configuration{
		Provider:          "coingecko",
		CoingeckoQueryUrl: ts.URL,
		DialTimeout:       time.Second,
		Bursts:            1,
		Tickers:           []string{"xxx"},
		Debug:             false,
	}
	n, err := NewHandler(*configuration)
	require.NoError(t, err)

	tm, _ := time.Parse("2006-01-02", "2021-03-04")

	out, err := n.GetOhlc(context, "xxx", "1d", tm, tm)
	require.Error(t, err)
	require.Nil(t, out)
}
