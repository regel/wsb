package iex

import (
	"github.com/regel/wsb/pkg/finance/types"
)

// Local type implements the types.Provider interface
type Provider struct {
	IexCloudQueryUrl    string
	IexCloudSecretToken string
}

func NewProvider(IexCloudQueryUrl string, IexCloudSecretToken string) types.Provider {
	return &Provider{
		IexCloudQueryUrl:    IexCloudQueryUrl,
		IexCloudSecretToken: IexCloudSecretToken,
	}
}
