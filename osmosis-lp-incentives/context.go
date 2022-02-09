package main

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

// Context contains all ingredients for this application.
type Context struct {
	conf            Config
	rpcClient       rpcclient.Client
	totalIncentives types.Coins
}

func newContext(conf Config) (*Context, error) {
	rpcClient, err := client.NewClientFromNode(conf.NodeRPCAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to NewClientFromNode: %w", err)
	}

	return &Context{
		conf:            conf,
		rpcClient:       rpcClient,
		totalIncentives: types.NewCoins(),
	}, nil
}

// addIncentives cumulates incentive coins to the totalIncentives.
func (ctx *Context) addIncentives(coins types.Coins) {
	ctx.totalIncentives = ctx.totalIncentives.Add(coins...)
}
