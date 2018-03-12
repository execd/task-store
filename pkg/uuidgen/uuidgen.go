package uuidgen

import "github.com/satori/go.uuid"

// UUIDGen : interface for a uuid generator
type UUIDGen interface {
	GenV4() (uuid.UUID, error)
}

// UUIDGenImpl : implementation for a uuid generator
type UUIDGenImpl struct{}

// NewUUIDGenImpl : build a new UUIDGenImpl
func NewUUIDGenImpl() *UUIDGenImpl {
	return &UUIDGenImpl{}
}

// GenV4 : generate a v4 uuid
func (u *UUIDGenImpl) GenV4() (uuid.UUID, error) {
	return uuid.NewV4()
}
