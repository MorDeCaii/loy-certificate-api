package retranslator

import (
	"fmt"
	"github.com/Mordecaii/loy-certificate-api/internal/mocks"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func TestStart(t *testing.T) {

	ctrl := gomock.NewController(t)
	fmt.Println(ctrl)
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
	retranslator.Close()
}