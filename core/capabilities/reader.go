package capabilities

import (
	"context"
	"encoding/json"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/keystone_capability_registry"
	evmrelaytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type remoteRegistryReader struct {
	r types.ContractReader
}

type don struct {
	// ID is the unique identifier of the DON in the CR.
	ID uint32

	// IsPublic indicates whether this DON's capabilities can be accessed publicly.
	IsPublic bool

	// Nodes is the list of nodes in this DON, represented by their RageP2P public keys/IDs.
	Nodes [][]byte
}

type state struct {
	DONs []don
}

func (r *remoteRegistryReader) state(ctx context.Context) (state, error) {
	dons := []don{}
	err := r.r.GetLatestValue(ctx, "capabilityRegistry", "getDONs", nil, &dons)
	if err != nil {
		return state{}, err
	}

	return state{DONs: dons}, nil
}

type contractReaderFactory interface {
	NewContractReader(context.Context, []byte) (types.ContractReader, error)
}

func newRemoteRegistryReader(ctx context.Context, relayer contractReaderFactory, remoteRegistryAddress string) (*remoteRegistryReader, error) {
	contractReaderConfig := evmrelaytypes.ChainReaderConfig{
		Contracts: map[string]evmrelaytypes.ChainContractReader{
			"capabilityRegistry": {
				ContractABI: keystone_capability_registry.CapabilityRegistryABI,
				Configs: map[string]*evmrelaytypes.ChainReaderDefinition{
					"getDONs": {
						ChainSpecificName: "getDONs",
					},
				},
			},
		},
	}

	contractReaderConfigEncoded, err := json.Marshal(contractReaderConfig)
	if err != nil {
		return nil, err
	}

	cr, err := relayer.NewContractReader(ctx, contractReaderConfigEncoded)
	if err != nil {
		return nil, err
	}

	err = cr.Bind(ctx, []types.BoundContract{
		{
			Address: remoteRegistryAddress,
			Name:    "capabilityRegistry",
		},
	})
	if err != nil {
		return nil, err
	}

	return &remoteRegistryReader{r: cr}, err
}
