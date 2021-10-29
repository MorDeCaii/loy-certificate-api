package repo

import "github.com/Mordecaii/loy-certificate-api/internal/model"

type EventRepo interface {
	Lock(n uint64) ([]model.CertificateEvent, error)
	Unlock(eventIDs []uint64) error

	Add(event []model.CertificateEvent) error
	Remove(eventIDs []uint64) error
}
