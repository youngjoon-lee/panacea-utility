module github.com/medibloc/panacea-utility/osmosis-lp-incentives

go 1.16

require (
	github.com/cosmos/cosmos-sdk v0.44.5
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/osmosis-labs/osmosis v1.0.3-0.20211217093711-dfe74a35ccc7
	github.com/sirupsen/logrus v1.8.1
	github.com/tendermint/tendermint v0.34.14
)

replace (
	github.com/cosmos/cosmos-sdk => github.com/osmosis-labs/cosmos-sdk v0.43.0-rc3.0.20211209072213-711e78b4f6b4
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/tendermint/tendermint => github.com/tendermint/tendermint v0.34.14
	github.com/tendermint/tm-db => github.com/osmosis-labs/tm-db v0.6.5-0.20210911033928-ba9154613417
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
)
