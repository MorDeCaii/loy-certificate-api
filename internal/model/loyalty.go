package model

import "fmt"

type Certificate struct {
	ID uint64
}

type EventType uint8

type EventStatus uint8

const (
	Created EventType = iota
	Updated
	Removed

	Deffered EventStatus = iota
	Processed
)

type CertificateEvent struct {
	ID     uint64
	Type   EventType
	Status EventStatus
	Entity *Certificate
}

func (c *CertificateEvent) String() string {
	return fmt.Sprintf("ID: %d", c.ID)
}
