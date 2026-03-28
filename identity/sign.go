package identity

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
)

// Sign creates an Ed25519 signature over the data.
func Sign(privKey ed25519.PrivateKey, data []byte) (string, error) {
	if privKey == nil {
		return "", fmt.Errorf("private key is nil")
	}
	sig := ed25519.Sign(privKey, data)
	return base64.StdEncoding.EncodeToString(sig), nil
}

// Verify checks an Ed25519 signature.
func Verify(pubKey ed25519.PublicKey, data []byte, sig string) error {
	if pubKey == nil {
		return fmt.Errorf("public key is nil")
	}
	sigBytes, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return fmt.Errorf("decode signature: %w", err)
	}
	if !ed25519.Verify(pubKey, data, sigBytes) {
		return fmt.Errorf("invalid signature")
	}
	return nil
}

// SignableEnvelope defines the fields needed for envelope signing.
// This avoids a circular import with the envelope package.
type SignableEnvelope interface {
	SigningPayload() []byte
	SetSignature(sig string)
	GetSignature() string
}

// SignEnvelope signs the envelope's payload and sets the signature field.
func SignEnvelope(env SignableEnvelope, privKey ed25519.PrivateKey) error {
	if env == nil {
		return fmt.Errorf("envelope is nil")
	}
	sig, err := Sign(privKey, env.SigningPayload())
	if err != nil {
		return fmt.Errorf("sign envelope: %w", err)
	}
	env.SetSignature(sig)
	return nil
}

// VerifyEnvelope verifies the envelope's signature.
func VerifyEnvelope(env SignableEnvelope, pubKey ed25519.PublicKey) error {
	return Verify(pubKey, env.SigningPayload(), env.GetSignature())
}
