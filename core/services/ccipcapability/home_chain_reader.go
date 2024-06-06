package ccipcapability

import (
	"context"
)

// HomeChainReader is an interface for reading CCIP chain and OCR configurations from the home chain.
type HomeChainReader interface {
	GetAllChainConfigs(ctx context.Context) (map[uint64]ChainConfig, error)
	GetOCRConfig(ctx context.Context, donID uint32) (OCRConfig, error)
	IsHealthy() bool
}
