package tls_configuration

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"prerequisites-tester/config"
)

type PinnedCertConfigProducer struct{}

func (producer *PinnedCertConfigProducer) CreateConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true, // We're handling verification ourselves via pinning
		VerifyConnection: func(cs tls.ConnectionState) error {
			return checkPinnedHash(&cs, config.GetConfig().Certificate)
		},
	}
}

func (producer *PinnedCertConfigProducer) Name() string {
	return "Pinned Cert"
}

func checkPinnedHash(connState *tls.ConnectionState, pinnedHash string) error {
	for _, peerCert := range connState.PeerCertificates {
		derKey, err := x509.MarshalPKIXPublicKey(peerCert.PublicKey)
		if err != nil {
			return HandshakeError{Message: fmt.Sprint("Failed to marshal public key: %v", err)}
		}
		hash := sha256.Sum256(derKey)
		encodedHash := hex.EncodeToString(hash[:])
		if encodedHash == pinnedHash {
			return nil
		}
	}
	return HandshakeError{Message: "Public key hash does not match pinned hash"}
}
