package main

import (
	"context"
	"github.com/Mordecaii/loy-certificate-api/internal/app/retranslator"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	ctx := context.Background()
	sigs := make(chan os.Signal, 1)

	cfg := retranslator.Config{
		ChannelSize:   512,
		ConsumerCount: 2,
		ConsumeSize:   10,
		ProducerCount: 28,
		WorkerCount:   2,
	}

	retranslator := retranslator.NewRetranslator(ctx, cfg)
	retranslator.Start()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
}
