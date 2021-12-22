package coingecko

import (
	"github.com/regel/tinkerbell/pkg/finance/types"
)

// Local type implements the types.Provider interface
type Provider struct {
	CoingeckoQueryUrl    string
	CoingeckoSecretToken string
}

func NewProvider(CoingeckoQueryUrl string, CoingeckoSecretToken string) types.Provider {
	return &Provider{
		CoingeckoQueryUrl:    CoingeckoQueryUrl,
		CoingeckoSecretToken: CoingeckoSecretToken,
	}
}
