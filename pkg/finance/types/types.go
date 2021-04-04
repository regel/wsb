package types

import (
	"time"
)

type Ohlc struct {
	Ticker    string
	Timestamp time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    int64
}

type HoldersBreakdown struct {
	Ticker                            string
	PctSharesHeldbyAllInsider         float64
	PctSharesHeldbyInstitutions       float64
	PctFloatHeldbyInstitutions        float64
	NumberofInstitutionsHoldingShares int64
}

type HoldersRow struct {
	Holder       string
	Shares       int64
	DateReported time.Time
	PctOut       float64
	Value        int64
}

// Top Institutional Holders
// Top Mutual Fund Holders
type HoldersTable struct {
	Ticker string
	Rows   []HoldersRow
}
