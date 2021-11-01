package consumer

import (
	"context"
	"github.com/Mordecaii/loy-certificate-api/internal/app/repo"
	"github.com/Mordecaii/loy-certificate-api/internal/model"
	"time"
)

type Consumer interface {
	Start()
	Close()
}

type consumer struct {
	n      uint64
	events chan<- model.CertificateEvent

	repo repo.EventRepo

	batchSize uint64
	timeout   time.Duration

	ctx    context.Context
	cancel func()
}

type Config struct {
	n         uint64
	events    chan<- model.CertificateEvent
	repo      repo.EventRepo
	batchSize uint64
	timeout   time.Duration
}

func NewDbConsumer(
	ctx context.Context,
	n uint64,
	batchSize uint64,
	consumeTimeout time.Duration,
	repo repo.EventRepo,
	events chan<- model.CertificateEvent) Consumer {

	consumerCtx, cancel := context.WithCancel(ctx)

	return &consumer{
		n:         n,
		batchSize: batchSize,
		timeout:   consumeTimeout,
		repo:      repo,
		events:    events,
		ctx:       consumerCtx,
		cancel:    cancel,
	}
}

func (c *consumer) Start() {
	for i := uint64(0); i < c.n; i++ {
		go func() {
			ticker := time.NewTicker(c.timeout)
			for {
				select {
				case <-ticker.C:
					events, err := c.repo.Lock(c.batchSize)
					if err != nil {
						continue
					}
					for _, event := range events {
						c.events <- event
					}
				case <-c.ctx.Done():
					return
				}
			}
		}()
	}
}

func (c *consumer) Close() {
	c.cancel()
}
