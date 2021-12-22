package yahoo

import (
	"github.com/regel/tinkerbell/pkg/finance/types"
)

// Local type implements the types.Provider interface
type Provider struct {
	YahooFinanceUrl      string
	YahooFinanceQueryUrl string
}

func NewProvider(YahooFinanceUrl string, YahooFinanceQueryUrl string) types.Provider {
	return &Provider{
		YahooFinanceUrl:      YahooFinanceUrl,
		YahooFinanceQueryUrl: YahooFinanceQueryUrl,
	}
}
