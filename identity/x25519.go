package identity

import (
	"crypto/ecdh"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
)

// X25519PublicKey derives an X25519 public key from the Ed25519 keypair's seed.
// This uses the standard conversion: hash the Ed25519 seed with SHA-512, clamp
// the first 32 bytes, and use that as the X25519 private key to derive the public key.
func (kp *Keypair) X25519PublicKey() (*ecdh.PublicKey, error) {
	priv, err := kp.X25519PrivateKey()
	if err != nil {
		return nil, err
	}
	return priv.PublicKey(), nil
}

// X25519PrivateKey derives an X25519 private key from the Ed25519 keypair's seed.
// The derivation matches the standard Ed25519-to-X25519 conversion used by libsodium:
// SHA-512 of the seed, clamp the lower 32 bytes.
func (kp *Keypair) X25519PrivateKey() (*ecdh.PrivateKey, error) {
	if kp == nil || kp.PrivateKey == nil {
		return nil, fmt.Errorf("keypair not initialized")
	}
	seed := kp.PrivateKey.Seed()
	h := sha512.Sum512(seed)
	// Clamp (RFC 7748 §5)
	h[0] &= 248
	h[31] &= 127
	h[31] |= 64
	return ecdh.X25519().NewPrivateKey(h[:32])
}

// X25519PublicKeyString returns the base64-encoded X25519 public key.
func (kp *Keypair) X25519PublicKeyString() (string, error) {
	pub, err := kp.X25519PublicKey()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(pub.Bytes()), nil
}

// ParseX25519PublicKey decodes a base64-encoded X25519 public key string.
func ParseX25519PublicKey(s string) (*ecdh.PublicKey, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return ecdh.X25519().NewPublicKey(data)
}

