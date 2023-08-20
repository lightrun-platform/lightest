package tls_configuration

import (
	"crypto/tls"
)

type TlsConfigProducer interface {
	CreateConfig() *tls.Config
	Name() string
}

type HandshakeError struct {
	Message string
}

func (e HandshakeError) Error() string {
	return e.Message
}
