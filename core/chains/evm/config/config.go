package config

import (
	"math/big"
	"net/url"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"

	commonassets "github.com/smartcontractkit/chainlink-common/pkg/assets"

	commonconfig "github.com/smartcontractkit/chainlink/v2/common/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

type EVM interface {
	HeadTracker() HeadTracker
	BalanceMonitor() BalanceMonitor
	Transactions() Transactions
	GasEstimator() GasEstimator
	OCR() OCR
	OCR2() OCR2
	Workflow() Workflow
	NodePool() NodePool

	AutoCreateKey() bool
	BlockBackfillDepth() uint64
	BlockBackfillSkip() bool
	BlockEmissionIdleWarningThreshold() time.Duration
	ChainID() *big.Int
	ChainType() commonconfig.ChainType
	FinalityDepth() uint32
	FinalityTagEnabled() bool
	FlagsContractAddress() string
	LinkContractAddress() string
	LogBackfillBatchSize() uint32
	LogKeepBlocksDepth() uint32
	BackupLogPollerBlockDelay() uint64
	LogPollInterval() time.Duration
	LogPrunePageSize() uint32
	MinContractPayment() *commonassets.Link
	MinIncomingConfirmations() uint32
	NonceAutoSync() bool
	OperatorFactoryAddress() string
	RPCDefaultBatchSize() uint32
	NodeNoNewHeadsThreshold() time.Duration

	IsEnabled() bool
	TOMLString() (string, error)
}

type OCR interface {
	ContractConfirmations() uint16
	ContractTransmitterTransmitTimeout() time.Duration
	ObservationGracePeriod() time.Duration
	DatabaseTimeout() time.Duration
	DeltaCOverride() time.Duration
	DeltaCJitterOverride() time.Duration
}

type OCR2 interface {
	Automation() OCR2Automation
}

type OCR2Automation interface {
	GasLimit() uint32
}

type HeadTracker interface {
	HistoryDepth() uint32
	MaxBufferSize() uint32
	SamplingInterval() time.Duration
	FinalityTagBypass() bool
	MaxAllowedFinalityDepth() uint32
}

type BalanceMonitor interface {
	Enabled() bool
}

type ClientErrors interface {
	NonceTooLow() string
	NonceTooHigh() string
	ReplacementTransactionUnderpriced() string
	LimitReached() string
	TransactionAlreadyInMempool() string
	TerminallyUnderpriced() string
	InsufficientEth() string
	TxFeeExceedsCap() string
	L2FeeTooLow() string
	L2FeeTooHigh() string
	L2Full() string
	TransactionAlreadyMined() string
	Fatal() string
	ServiceUnavailable() string
}

type Transactions interface {
	ForwardersEnabled() bool
	ReaperInterval() time.Duration
	ResendAfterThreshold() time.Duration
	ReaperThreshold() time.Duration
	MaxInFlight() uint32
	MaxQueued() uint64
	AutoPurge() AutoPurgeConfig
}

type AutoPurgeConfig interface {
	Enabled() bool
	Threshold() uint32
	MinAttempts() uint32
	DetectionApiUrl() *url.URL
}

//go:generate mockery --quiet --name GasEstimator --output ./mocks/ --case=underscore
type GasEstimator interface {
	BlockHistory() BlockHistory
	LimitJobType() LimitJobType

	EIP1559DynamicFees() bool
	BumpPercent() uint16
	BumpThreshold() uint64
	BumpTxDepth() uint32
	BumpMin() *assets.Wei
	FeeCapDefault() *assets.Wei
	LimitDefault() uint64
	LimitMax() uint64
	LimitMultiplier() float32
	LimitTransfer() uint64
	PriceDefault() *assets.Wei
	TipCapDefault() *assets.Wei
	TipCapMin() *assets.Wei
	PriceMax() *assets.Wei
	PriceMin() *assets.Wei
	Mode() string
	PriceMaxKey(gethcommon.Address) *assets.Wei
}

type LimitJobType interface {
	OCR() *uint32
	OCR2() *uint32
	DR() *uint32
	FM() *uint32
	Keeper() *uint32
	VRF() *uint32
}

type BlockHistory interface {
	BatchSize() uint32
	BlockHistorySize() uint16
	BlockDelay() uint16
	CheckInclusionBlocks() uint16
	CheckInclusionPercentile() uint16
	EIP1559FeeCapBufferBlocks() uint16
	TransactionPercentile() uint16
}

type Workflow interface {
	FromAddress() *types.EIP55Address
	ForwarderAddress() *types.EIP55Address
}

type NodePool interface {
	PollFailureThreshold() uint32
	PollInterval() time.Duration
	SelectionMode() string
	SyncThreshold() uint32
	LeaseDuration() time.Duration
	NodeIsSyncingEnabled() bool
	FinalizedBlockPollInterval() time.Duration
	Errors() ClientErrors
}

// TODO BCF-2509 does the chainscopedconfig really need the entire app config?
//
//go:generate mockery --quiet --name ChainScopedConfig --output ./mocks/ --case=underscore
type ChainScopedConfig interface {
	EVM() EVM
}
