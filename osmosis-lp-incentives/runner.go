package main

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/types"
	osmotypes "github.com/osmosis-labs/osmosis/x/incentives/types"
	log "github.com/sirupsen/logrus"
	abci "github.com/tendermint/tendermint/abci/types"
)

func run(ctx *Context) error {
	height := ctx.conf.StartHeight

	for height <= ctx.conf.EndHeight {
		log.Debugf("getting block results for height %v", height)
		blockResults, err := getBlockResults(ctx, height)
		if err != nil {
			return fmt.Errorf("failed to getBlockResults for %v: %w", height, err)
		}

		// 'distribution' events can be found in BeginBlockEvents.
		log.Debugf("finding distribution events for height %v", height)
		found, err := findAndSumUpDistributionEvents(ctx, blockResults.BeginBlockEvents, blockResults.Height)
		if err != nil {
			return fmt.Errorf("failed to findAndSumUpDistributionEvents: %w", err)
		}

		height, err = jumpBlocks(ctx, height, found)
		if err != nil {
			return fmt.Errorf("failed to jumpBlocks: %w", err)
		}
	}

	return nil
}

// findAndSumUpDistributionEvents finds 'distribution' events which are created by x/incentive of Osmosis.
// If 'distribution' events are found, the amount of distributions is cumulated into Context.
func findAndSumUpDistributionEvents(ctx *Context, events []abci.Event, height int64) (bool, error) {
	found := false

	for _, event := range events {
		if event.Type == osmotypes.TypeEvtDistribution {
			// https://github.com/osmosis-labs/osmosis/blob/dfe74a35ccc7f8a182b7a9ca943e9bc720ab1f51/x/incentives/keeper/gauge.go#L357
			// Attributes[0].Key is 'receiver'
			if string(event.Attributes[0].Value) == ctx.conf.TargetAccAddr {
				found = true

				// Attributes[1].Key is 'amount'
				coinsStr := string(event.Attributes[1].Value)
				coins, err := types.ParseCoinsNormalized(coinsStr)
				if err != nil {
					return false, fmt.Errorf("failed to ParseCoinsNormalized(%v): %w", coinsStr, err)
				}

				ctx.addIncentives(coins)

				// get a block only for detailed logging
				block, err := getBlock(ctx, uint64(height))
				if err != nil {
					return false, fmt.Errorf("failed to getBlock(%v): %w", height, err)
				}
				log.Infof("epoch (%v): %v: %v", height, block.Block.Header.Time.Format(time.RFC3339), coinsStr)
			}
		}
	}

	return found, nil
}

// jumpBlocks calculates a next block height that is close to the next epoch block as much as possible.
// If the next epoch block is already close enough, jumpBlocks just returns the next block (from the current block).
func jumpBlocks(ctx *Context, curHeight uint64, foundEvents bool) (uint64, error) {
	curBlock, err := getBlock(ctx, curHeight)
	if err != nil {
		return 0, fmt.Errorf("failed to getBlock(%v): %w", curHeight, err)
	}
	curTime := curBlock.Block.Header.Time

	targetJumpTime := getTargetJumpTime(ctx, curTime, foundEvents)
	secToJump := targetJumpTime.Sub(curTime).Seconds()
	blocksToJump := uint64(secToJump / ctx.conf.SafeBlockIntervalSec)

	for {
		if blocksToJump == 0 {
			// do not jump. just go to the next block (height).
			return curHeight + 1, nil
		}

		newHeight := curHeight + blocksToJump
		newBlock, err := getBlock(ctx, newHeight)
		if err != nil {
			return 0, fmt.Errorf("failed to getBlock(%v): %w", newHeight, err)
		}
		newTime := newBlock.Block.Header.Time

		// if we jumped too much, retry jumping slowly.
		if newTime.After(targetJumpTime) {
			log.Debugf("jumped too much. adjusting blocksToJump: %d -> %d", blocksToJump, blocksToJump/2)
			blocksToJump /= 2
			continue
		}

		log.Infof("jumped from %v (%v) to %v (%v)", curHeight, curTime.Format(time.RFC3339), newHeight, newTime.Format(time.RFC3339))

		return newHeight, nil
	}
}

// getTargetJumpTime calculates the target time that we want to jump to.
// The target time will not be after the next epoch time because this algorithm uses enough margins.
func getTargetJumpTime(ctx *Context, curTime time.Time, foundEvents bool) time.Time {
	epochHour := int(ctx.conf.EpochHour)
	epochMin := int(ctx.conf.EpochMinute)

	if curTime.Hour() < epochHour ||
		(curTime.Hour() == epochHour && curTime.Minute() < epochMin) {
		// jump to the epoch hour:minute:00 in the same date.
		return time.Date(curTime.Year(), curTime.Month(), curTime.Day(), epochHour, epochMin, 0, 0, time.UTC)
	} else if curTime.Hour() > epochHour ||
		(curTime.Hour() == epochHour && curTime.Minute() > epochMin) ||
		foundEvents {
		// jump to the epoch hour:minute:00 in tomorrow.
		tomorrow := curTime.Add(24 * time.Hour)
		return time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), epochHour, epochMin, 0, 0, time.UTC)
	}

	// do not jump
	return curTime
}
