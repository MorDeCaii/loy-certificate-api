package producer

import (
	"context"
	"fmt"
	"github.com/Mordecaii/loy-certificate-api/internal/app/repo"
	"github.com/Mordecaii/loy-certificate-api/internal/app/sender"
	"github.com/Mordecaii/loy-certificate-api/internal/model"
	"github.com/gammazero/workerpool"
	"time"
)

type Producer interface {
	Start()
	Close()
}

type producer struct {
	n       uint64
	timeout time.Duration
	repo    repo.EventRepo

	sender sender.EventSender
	events <-chan model.CertificateEvent

	workerPool *workerpool.WorkerPool

	ctx    context.Context
	cancel func()
}

func NewKafkaProducer(
	ctx context.Context,
	n uint64,
	repo repo.EventRepo,
	sender sender.EventSender,
	events <-chan model.CertificateEvent,
	workerPool *workerpool.WorkerPool,
) Producer {

	producerCtx, cancel := context.WithCancel(ctx)

	return &producer{
		n:          n,
		repo:       repo,
		sender:     sender,
		events:     events,
		workerPool: workerPool,
		ctx:        producerCtx,
		cancel:     cancel,
	}
}

func (p *producer) Start() {
	for i := uint64(0); i < p.n; i++ {
		go func() {
			for {
				select {
				case event := <-p.events:
					if err := p.sender.Send(&event); err != nil {
						fmt.Printf("Error while sending event {%s}: %s\n", event.String(), err)
						p.updateEvent(&event)
					} else {
						fmt.Printf("Event {%s} successfully sent to Kafka\n", event.String())
						p.cleanEvent(&event)

					}
				case <-p.ctx.Done():
					return
				}
			}
		}()
	}
}

func (p *producer) Close() {
	p.cancel()
}

func (p *producer) cleanEvent(event *model.CertificateEvent) {
	p.workerPool.Submit(func() {
		err := p.repo.Remove([]uint64{event.ID})
		if err != nil {
			fmt.Printf("Error while cleaning event: %s\n", err)
		}
	})
}

func (p *producer) updateEvent(event *model.CertificateEvent) {
	p.workerPool.Submit(func() {
		err := p.repo.Unlock([]uint64{event.ID})
		if err != nil {
			fmt.Printf("Error while updating event: %s\n", err)
		}
	})
}
