package tls_configuration

import (
	"crypto/tls"
)

type RootCertsConfigProducer struct {
}

func (producer *RootCertsConfigProducer) CreateConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: false,
	}
}

func (producer *RootCertsConfigProducer) Name() string {
	return "Root Certs"
}
