package capabilities

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/keystone_capability_registry"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	evmrelaytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

func stringToByteWord(in string) [32]byte {
	out := [32]byte{}
	copy(out[:], []byte(in))
	return out
}

var dataStreamsReportCapability = kcr.CapabilityRegistryCapability{
	LabelledName: stringToByteWord("data-streams-report"),
	Version:      stringToByteWord("1.0.0"),
	ResponseType: uint8(0),
}

var writeChainCapability = kcr.CapabilityRegistryCapability{
	LabelledName: stringToByteWord("write-chain"),
	Version:      stringToByteWord("1.0.1"),
	ResponseType: uint8(1),
}

func startNewChainWithRegistry(t *testing.T) (*kcr.CapabilityRegistry, *bind.TransactOpts, *backends.SimulatedBackend) {
	owner := testutils.MustNewSimTransactor(t)

	oneEth, _ := new(big.Int).SetString("100000000000000000000", 10)
	gasLimit := ethconfig.Defaults.Miner.GasCeil * 2 // 60 M blocks

	simulatedBackend := backends.NewSimulatedBackend(core.GenesisAlloc{owner.From: {
		Balance: oneEth,
	}}, gasLimit)
	simulatedBackend.Commit()

	capabilityRegistryAddress, _, capabilityRegistry, err := kcr.DeployCapabilityRegistry(owner, simulatedBackend)
	require.NoError(t, err, "DeployCapabilityRegistry failed")

	fmt.Println("Deployed CapabilityRegistry at", capabilityRegistryAddress.Hex())
	simulatedBackend.Commit()

	return capabilityRegistry, owner, simulatedBackend
}

type crFactory struct {
	lggr      logger.Logger
	logPoller logpoller.LogPoller
	client    evmclient.Client
}

func (c *crFactory) NewContractReader(ctx context.Context, cfg []byte) (types.ContractReader, error) {
	crCfg := &evmrelaytypes.ChainReaderConfig{}
	if err := json.Unmarshal(cfg, crCfg); err != nil {
		return nil, err
	}
	return evm.NewChainReaderService(ctx, c.lggr, c.logPoller, c.client, *crCfg)
}

func newChainReaderFactory(t *testing.T, simulatedBackend *backends.SimulatedBackend) *crFactory {
	lggr := logger.TestLogger(t)
	client := evmclient.NewSimulatedBackendClient(
		t,
		simulatedBackend,
		testutils.SimulatedChainID,
	)
	db := pgtest.NewSqlxDB(t)
	lp := logpoller.NewLogPoller(
		logpoller.NewORM(testutils.SimulatedChainID, db, lggr),
		client,
		lggr,
		logpoller.Opts{
			PollPeriod:               100 * time.Millisecond,
			FinalityDepth:            2,
			BackfillBatchSize:        3,
			RpcBatchSize:             2,
			KeepFinalizedBlocksDepth: 1000,
		},
	)
	return &crFactory{
		lggr:      lggr,
		client:    client,
		logPoller: lp,
	}
}

func TestReader(t *testing.T) {
	ctx := testutils.Context(t)
	reg, owner, sim := startNewChainWithRegistry(t)

	_, err := reg.AddCapability(owner, writeChainCapability)
	require.NoError(t, err, "AddCapability failed for %s", writeChainCapability.LabelledName)
	sim.Commit()

	factory := newChainReaderFactory(t, sim)
	_, err = newRemoteRegistryReader(ctx, factory, "")
	require.NoError(t, err)
}
