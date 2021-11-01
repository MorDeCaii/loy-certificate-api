package retranslator

import (
	"errors"
	"github.com/Mordecaii/loy-certificate-api/internal/mocks"
	"github.com/Mordecaii/loy-certificate-api/internal/model"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockEventRepo(ctrl)
	sender := mocks.NewMockEventSender(ctrl)

	repo.EXPECT().Lock(gomock.Any()).AnyTimes()

	cfg := Config{
		ChannelSize:    512,
		ConsumerCount:  2,
		ConsumeSize:    10,
		ConsumeTimeout: 10 * time.Second,
		ProducerCount:  2,
		WorkerCount:    2,
		Repo:           repo,
		Sender:         sender,
	}

	retranslator := NewRetranslator(cfg)
	retranslator.Start()
	time.Sleep(time.Second * 2)
	retranslator.Close()
}

func TestLockRemove(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockEventRepo(ctrl)
	sender := mocks.NewMockEventSender(ctrl)

	events := getEvents()

	lockEventsCall := repo.EXPECT().Lock(uint64(5)).Return(events, nil).Times(1)
	repo.EXPECT().Lock(gomock.Any()).Return(nil, errors.New("no new event available")).AnyTimes().After(lockEventsCall)
	sender.EXPECT().Send(gomock.Any()).Times(5).DoAndReturn(func(certificate *model.CertificateEvent) error {
		return nil
	})
	repo.EXPECT().Remove(gomock.Any()).Times(5)
	repo.EXPECT().Unlock(gomock.Any()).Times(0)

	cfg := Config{
		ChannelSize:    512,
		ConsumerCount:  2,
		ConsumeSize:    5,
		ConsumeTimeout: 1 * time.Second,
		ProducerCount:  2,
		WorkerCount:    2,
		Repo:           repo,
		Sender:         sender,
	}

	retranslator := NewRetranslator(cfg)
	retranslator.Start()
	time.Sleep(time.Second * 2)
	retranslator.Close()
}

func TestLockRemoveMatchIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockEventRepo(ctrl)
	sender := mocks.NewMockEventSender(ctrl)

	events := getEvents()

	var sentIDs, removedIDs []uint64
	lockEventsCall := repo.EXPECT().Lock(uint64(5)).Return(events, nil).Times(1)
	repo.EXPECT().Lock(gomock.Any()).Return(nil, errors.New("no new event available")).AnyTimes().After(lockEventsCall)
	sender.EXPECT().Send(gomock.Any()).Times(5).DoAndReturn(func(certificate *model.CertificateEvent) error {
		sentIDs = append(sentIDs, certificate.ID)
		return nil
	})
	repo.EXPECT().Unlock(gomock.Any()).Times(0)
	repo.EXPECT().Remove(gomock.Any()).Times(5).DoAndReturn(func(eventIDs []uint64) error {
		removedIDs = append(removedIDs, eventIDs...)
		return nil
	})

	if !gomock.InAnyOrder(sentIDs).Matches(removedIDs) {
		t.Errorf("Sent IDs aren't match with removed IDs")
	}

	cfg := Config{
		ChannelSize:    512,
		ConsumerCount:  2,
		ConsumeSize:    5,
		ConsumeTimeout: 1 * time.Second,
		ProducerCount:  2,
		WorkerCount:    2,
		Repo:           repo,
		Sender:         sender,
	}

	retranslator := NewRetranslator(cfg)
	retranslator.Start()
	time.Sleep(time.Second * 2)
	retranslator.Close()
}

func TestLockUnlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockEventRepo(ctrl)
	sender := mocks.NewMockEventSender(ctrl)

	events := getEvents()

	lockEventsCall := repo.EXPECT().Lock(uint64(5)).Return(events, nil).Times(1)
	repo.EXPECT().Lock(gomock.Any()).Return(nil, errors.New("no new event available")).AnyTimes().After(lockEventsCall)
	sender.EXPECT().Send(gomock.Any()).Times(5).DoAndReturn(func(certificate *model.CertificateEvent) error {
		return errors.New("test error")
	})
	repo.EXPECT().Remove(gomock.Any()).Times(0)
	repo.EXPECT().Unlock(gomock.Any()).Times(5)

	cfg := Config{
		ChannelSize:    512,
		ConsumerCount:  2,
		ConsumeSize:    5,
		ConsumeTimeout: 1 * time.Second,
		ProducerCount:  2,
		WorkerCount:    2,
		Repo:           repo,
		Sender:         sender,
	}

	retranslator := NewRetranslator(cfg)
	retranslator.Start()
	time.Sleep(time.Second * 2)
	retranslator.Close()
}

func TestLockUnlockMatchIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockEventRepo(ctrl)
	sender := mocks.NewMockEventSender(ctrl)

	events := getEvents()

	var sentIDs, unlockedIDs []uint64
	lockEventsCall := repo.EXPECT().Lock(uint64(5)).Return(events, nil).Times(1)
	repo.EXPECT().Lock(gomock.Any()).Return(nil, errors.New("no new event available")).AnyTimes().After(lockEventsCall)
	sender.EXPECT().Send(gomock.Any()).Times(5).DoAndReturn(func(certificate *model.CertificateEvent) error {
		sentIDs = append(sentIDs, certificate.ID)
		return errors.New("test error")
	})
	repo.EXPECT().Unlock(gomock.Any()).Times(5).DoAndReturn(func(eventIDs []uint64) error {
		unlockedIDs = append(unlockedIDs, eventIDs...)
		return nil
	})
	repo.EXPECT().Remove(gomock.Any()).Times(0)

	if !gomock.InAnyOrder(sentIDs).Matches(unlockedIDs) {
		t.Errorf("Sent IDs aren't match with unlocked IDs")
	}

	cfg := Config{
		ChannelSize:    512,
		ConsumerCount:  2,
		ConsumeSize:    5,
		ConsumeTimeout: 1 * time.Second,
		ProducerCount:  2,
		WorkerCount:    2,
		Repo:           repo,
		Sender:         sender,
	}

	retranslator := NewRetranslator(cfg)
	retranslator.Start()
	time.Sleep(time.Second * 2)
	retranslator.Close()
}

func getEvents() []model.CertificateEvent {
	return []model.CertificateEvent{
		{
			ID:     1,
			Type:   model.Created,
			Status: model.Deffered,
			Entity: &model.Certificate{ID: 1},
		},
		{
			ID:     2,
			Type:   model.Created,
			Status: model.Deffered,
			Entity: &model.Certificate{ID: 2},
		},
		{
			ID:     3,
			Type:   model.Created,
			Status: model.Deffered,
			Entity: &model.Certificate{ID: 3},
		},
		{
			ID:     4,
			Type:   model.Created,
			Status: model.Deffered,
			Entity: &model.Certificate{ID: 4},
		},
		{
			ID:     5,
			Type:   model.Created,
			Status: model.Deffered,
			Entity: &model.Certificate{ID: 5},
		},
	}
}
