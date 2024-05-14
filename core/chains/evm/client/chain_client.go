package client

import (
	"context"
	"math/big"
	"time"

	evmconfig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethrpc "github.com/ethereum/go-ethereum/rpc"

	commonassets "github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	commonclient "github.com/smartcontractkit/chainlink/v2/common/client"
	"github.com/smartcontractkit/chainlink/v2/common/config"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

var _ Client = (*chainClient)(nil)

// TODO-1663: rename this to client, once the client.go file is deprecated.
type chainClient struct {
	multiNode commonclient.MultiNode[
		*big.Int,
		common.Hash,
		*evmtypes.Head,
		*RpcClient,
	]
	logger       logger.SugaredLogger
	chainType    config.ChainType
	clientErrors evmconfig.ClientErrors
}

func NewChainClient(
	lggr logger.Logger,
	selectionMode string,
	leaseDuration time.Duration,
	noNewHeadsThreshold time.Duration,
	nodes []commonclient.Node[*big.Int, *evmtypes.Head, *RpcClient],
	sendonlys []commonclient.Node[*big.Int, *evmtypes.Head, *RpcClient],
	chainID *big.Int,
	chainType config.ChainType,
	clientErrors evmconfig.ClientErrors,
) Client {
	multiNode := commonclient.NewMultiNode(
		lggr,
		selectionMode,
		leaseDuration,
		nodes,
		sendonlys,
		chainID,
		"EVM",
	)
	return &chainClient{
		multiNode:    multiNode,
		logger:       logger.Sugared(lggr),
		clientErrors: clientErrors,
	}
}

func (c *chainClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return nil, err
	}
	return rpc.BalanceAt(ctx, account, blockNumber)
}

// Request specific errors for batch calls are returned to the individual BatchElem.
// Ensure the same BatchElem slice provided by the caller is passed through the call stack
// to ensure the caller has access to the errors.
func (c *chainClient) BatchCallContext(ctx context.Context, b []ethrpc.BatchElem) error {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return err
	}
	return rpc.BatchCallContext(ctx, b)
}

// Similar to BatchCallContext, ensure the provided BatchElem slice is passed through
func (c *chainClient) BatchCallContextAll(ctx context.Context, b []ethrpc.BatchElem) error {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return err
	}
	// TODO: What should we do here? c.multiNode.DoAll()?
	return rpc.BatchCallContextAll(ctx, b)
}

// TODO-1663: return custom Block type instead of geth's once client.go is deprecated.
func (c *chainClient) BlockByHash(ctx context.Context, hash common.Hash) (b *types.Block, err error) {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return b, err
	}
	return rpc.BlockByHashGeth(ctx, hash)
}

// TODO-1663: return custom Block type instead of geth's once client.go is deprecated.
func (c *chainClient) BlockByNumber(ctx context.Context, number *big.Int) (b *types.Block, err error) {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return b, err
	}
	return rpc.BlockByNumberGeth(ctx, number)
}

func (c *chainClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return err
	}
	return rpc.CallContext(ctx, result, method, args...)
}

func (c *chainClient) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return nil, err
	}
	return rpc.CallContract(ctx, msg, blockNumber)
}

func (c *chainClient) PendingCallContract(ctx context.Context, msg ethereum.CallMsg) ([]byte, error) {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return nil, err
	}
	return rpc.PendingCallContract(ctx, msg)
}

// TODO-1663: change this to actual ChainID() call once client.go is deprecated.
func (c *chainClient) ChainID() (*big.Int, error) {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return nil, err
	}
	return rpc.chainID, nil
}

func (c *chainClient) Close() {
	_ = c.multiNode.Close()
}

func (c *chainClient) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return nil, err
	}
	return rpc.CodeAt(ctx, account, blockNumber)
}

func (c *chainClient) ConfiguredChainID() *big.Int {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return nil
	}
	return rpc.chainID
}

func (c *chainClient) Dial(ctx context.Context) error {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return err
	}
	return rpc.Dial(ctx)
}

func (c *chainClient) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return nil, err
	}
	return rpc.EstimateGas(ctx, call)
}
func (c *chainClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return nil, err
	}
	return rpc.FilterEvents(ctx, q)
}

func (c *chainClient) HeaderByHash(ctx context.Context, h common.Hash) (head *types.Header, err error) {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return head, err
	}
	return rpc.HeaderByHash(ctx, h)
}

func (c *chainClient) HeaderByNumber(ctx context.Context, n *big.Int) (head *types.Header, err error) {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return head, err
	}
	return rpc.HeaderByNumber(ctx, n)
}

func (c *chainClient) HeadByHash(ctx context.Context, h common.Hash) (*evmtypes.Head, error) {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return nil, err
	}
	return rpc.BlockByHash(ctx, h)
}

func (c *chainClient) HeadByNumber(ctx context.Context, n *big.Int) (*evmtypes.Head, error) {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return nil, err
	}
	return rpc.BlockByNumber(ctx, n)
}

func (c *chainClient) IsL2() bool {
	// TODO: Where should this come from?
	return c.multiNode.IsL2()
}

func (c *chainClient) LINKBalance(ctx context.Context, address common.Address, linkAddress common.Address) (*commonassets.Link, error) {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return nil, err
	}
	return rpc.LINKBalance(ctx, address, linkAddress)
}

func (c *chainClient) LatestBlockHeight(ctx context.Context) (*big.Int, error) {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return nil, err
	}
	return rpc.LatestBlockHeight(ctx)
}

func (c *chainClient) NodeStates() map[string]string {
	// TODO: Should nodeState be public and returned here?
	return c.multiNode.NodeStates()
}

func (c *chainClient) PendingCodeAt(ctx context.Context, account common.Address) (b []byte, err error) {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return b, err
	}
	return rpc.PendingCodeAt(ctx, account)
}

// TODO-1663: change this to evmtypes.Nonce(int64) once client.go is deprecated.
func (c *chainClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return b, err
	}
	n, err := rpc.PendingSequenceAt(ctx, account)
	return uint64(n), err
}

func (c *chainClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return err
	}
	return rpc.SendTransaction(ctx, tx)
}

func (c *chainClient) SendTransactionReturnCode(ctx context.Context, tx *types.Transaction, fromAddress common.Address) (commonclient.SendTxReturnCode, error) {
	err := c.SendTransaction(ctx, tx)
	returnCode := ClassifySendError(err, c.clientErrors, c.logger, tx, fromAddress, c.IsL2())
	return returnCode, err
}

func (c *chainClient) SequenceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (evmtypes.Nonce, error) {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return b, err
	}
	return rpc.SequenceAt(ctx, account, blockNumber)
}

func (c *chainClient) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (s ethereum.Subscription, err error) {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return s, err
	}
	return rpc.SubscribeFilterLogs(ctx, q, ch)
}

func (c *chainClient) SubscribeNewHead(ctx context.Context) (chan<- *evmtypes.Head, ethereum.Subscription, error) {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return nil, nil, err
	}

	// TODO: Implement this
	rpc.SubscribeToHeads(ctx)
	csf := newChainIDSubForwarder(c.ConfiguredChainID(), ch)
	err := csf.start(c.multiNode.Subscribe(ctx, csf.srcCh, "newHeads"))
	if err != nil {
		return nil, err
	}
	return csf, nil
}

func (c *chainClient) SuggestGasPrice(ctx context.Context) (p *big.Int, err error) {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return p, err
	}
	return rpc.SuggestGasPrice(ctx)
}

func (c *chainClient) SuggestGasTipCap(ctx context.Context) (t *big.Int, err error) {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return t, err
	}
	return rpc.SuggestGasTipCap(ctx)
}

func (c *chainClient) TokenBalance(ctx context.Context, address common.Address, contractAddress common.Address) (*big.Int, error) {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return p, err
	}
	return rpc.TokenBalance(ctx, address, contractAddress)
}

func (c *chainClient) TransactionByHash(ctx context.Context, txHash common.Hash) (*types.Transaction, error) {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return p, err
	}
	return rpc.TransactionByHash(ctx, txHash)
}

// TODO-1663: return custom Receipt type instead of geth's once client.go is deprecated.
func (c *chainClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (r *types.Receipt, err error) {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return r, err
	}
	//return rpc.TransactionReceipt(ctx, txHash)
	return rpc.TransactionReceiptGeth(ctx, txHash)
}

func (c *chainClient) LatestFinalizedBlock(ctx context.Context) (*evmtypes.Head, error) {
	rpc, err := c.multiNode.SelectRPC()
	if err != nil {
		return p, err
	}
	return rpc.LatestFinalizedBlock(ctx)
}

func (c *chainClient) CheckTxValidity(ctx context.Context, from common.Address, to common.Address, data []byte) *SendError {
	msg := ethereum.CallMsg{
		From: from,
		To:   &to,
		Data: data,
	}
	return SimulateTransaction(ctx, c, c.logger, c.chainType, msg)
}
