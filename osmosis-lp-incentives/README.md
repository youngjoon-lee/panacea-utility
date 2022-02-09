# Osmosis LP Incentives Aggregator

This application aggregates historical Osmosis LP incentives of a specific Osmosis account.

## Background

Osmosis performs LP incentive distributions to all liquidity providers at every 1-day epoch (17:16:xx UTC).
However, Osmosis doesn't perform distributions as transactions because those jobs are done in the `x/incentives` BeginBlocker.

Thus, in order to aggregate daily LP incentives that a specific account receives,
this application iterates over blocks in a specific range and investigate begin-block events.

## How to use

Build source codes.
```bash
go build ./...
```

Run the application with environment variables.
```bash
LOG_LEVEL="debug" \  # panic, fatal, error, warn, info (default), debug, trace
NODE_RPC_ADDR="https://my-osmosis.rpc:443" \  # e.g. localhost:26657
START_HEIGHT="2383300" \  # investigation range
END_HEIGHT="2611465" \
SAFE_BLOCK_INTERVAL_SEC="7.0" \  # set this slightly larger than the actual average block interval (6.x sec)
EPOCH_HOUR="17" \    # the 'hour' part of the epoch time (17:16:xx in general)
EPOCH_MINUTE="16" \  # the 'minute' part of the epoch time (17:16:xx in general)
TARGET_ACCOUNT="osmo1..." \  # the account that we want to know about
./osmosis-lp-incentives
```

After the application investigates all blocks in the specified range, it logs the final summary.
```bash
INFO[2022-01-22T15:56:17+09:00] ============ FINAL =============             
INFO[2022-01-22T15:56:17+09:00] Height: 2383300 ~ 2384400                    
INFO[2022-01-22T15:56:17+09:00] Coins:                                       
INFO[2022-01-22T15:56:17+09:00] 235908091 ibc/3BCCC93AD5DF58D11A6F8A05FA8BC801CBA0BA61A981F57E91B8B598BF8061CB 
INFO[2022-01-22T15:56:17+09:00] 9440521 uosmo   
```