package main

import (
	"context"

	rpctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// getBlock retrieves a block at a specific height via Tendermint RPC.
func getBlock(ctx *Context, height uint64) (*rpctypes.ResultBlock, error) {
	h := int64(height)
	return ctx.rpcClient.Block(context.Background(), &h)
}

// getBlockResults retrieves block results at a specific height via Tendermint RPC.
func getBlockResults(ctx *Context, height uint64) (*rpctypes.ResultBlockResults, error) {
	h := int64(height)
	return ctx.rpcClient.BlockResults(context.Background(), &h)
}
