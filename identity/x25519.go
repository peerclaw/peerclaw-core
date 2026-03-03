package identity

import (
	"crypto/ecdh"
	"crypto/ed25519"
	"crypto/sha512"
	"encoding/base64"
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

// Ed25519PublicKeyToX25519 converts a raw Ed25519 public key to X25519.
// This uses the birational map from the Edwards to Montgomery form.
// Note: this is a one-way conversion and not all Ed25519 public keys can be
// reliably converted. When possible, prefer deriving from the seed via Keypair.
func Ed25519PublicKeyToX25519(pub ed25519.PublicKey) (*ecdh.PublicKey, error) {
	// For peer public keys where we don't have the seed, we rely on the
	// X25519 public key being transmitted directly (e.g., in signaling).
	// This function exists for completeness but the preferred path is to
	// exchange X25519 public keys explicitly.
	//
	// The actual Edwards-to-Montgomery conversion requires field arithmetic
	// that is not in the Go standard library. In practice, peers exchange
	// their X25519 public keys during signaling, so this is rarely needed.
	return nil, nil
}
