package sender

import "github.com/Mordecaii/loy-certificate-api/internal/model"

type EventSender interface {
	Send(certificate *model.CertificateEvent) error
}
