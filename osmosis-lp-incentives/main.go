package main

import (
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	conf := mustLoadConfig()
	setLogger(conf)
	log.Infof("%+v", conf)

	ctx, err := newContext(conf)
	if err != nil {
		log.Fatalf("failed to newContext: %v", err)
	}

	if err := run(ctx); err != nil {
		log.Fatal(err)
	}

	log.Infof("============ FINAL =============")
	log.Infof("Height: %v ~ %v", conf.StartHeight, conf.EndHeight)
	log.Infof("Coins:")
	for _, coin := range ctx.totalIncentives {
		log.Infof("%v %v", coin.Amount, coin.Denom)
	}
}

func setLogger(conf Config) {
	log.SetLevel(log.Level(conf.LogLevel))
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	})
}
