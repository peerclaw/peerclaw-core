package identity

import (
	"crypto/ecdh"
	"crypto/ed25519"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
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
// This uses the birational map from the Edwards to Montgomery form:
// given an Edwards point with y-coordinate encoded in the public key,
// the Montgomery u-coordinate is (1 + y) / (1 - y) mod p, where
// p = 2^255 - 19.
// Note: this is a one-way conversion. When possible, prefer deriving
// from the seed via Keypair.
func Ed25519PublicKeyToX25519(pub ed25519.PublicKey) (*ecdh.PublicKey, error) {
	if len(pub) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("ed25519 public key must be %d bytes, got %d", ed25519.PublicKeySize, len(pub))
	}

	// The Ed25519 public key is a compressed Edwards point: the y-coordinate
	// in 255 bits (little-endian) with the sign of x in the top bit.
	// Extract y by copying and clearing the sign bit.
	yBytes := make([]byte, 32)
	copy(yBytes, pub)
	yBytes[31] &= 0x7F // clear sign bit

	// Convert from little-endian to big.Int.
	// big.Int expects big-endian, so reverse the bytes.
	reversed := make([]byte, 32)
	for i := 0; i < 32; i++ {
		reversed[i] = yBytes[31-i]
	}

	y := new(big.Int).SetBytes(reversed)

	// p = 2^255 - 19
	p := new(big.Int).SetBit(new(big.Int), 255, 1)
	p.Sub(p, big.NewInt(19))

	// u = (1 + y) / (1 - y) mod p
	//   = (1 + y) * (1 - y)^(-1) mod p
	numerator := new(big.Int).Add(big.NewInt(1), y)
	numerator.Mod(numerator, p)

	denominator := new(big.Int).Sub(big.NewInt(1), y)
	denominator.Mod(denominator, p)

	// Check that denominator is not zero (y = 1 is a degenerate case).
	if denominator.Sign() == 0 {
		return nil, errors.New("degenerate ed25519 public key: y-coordinate is 1")
	}

	denomInv := new(big.Int).ModInverse(denominator, p)
	if denomInv == nil {
		return nil, errors.New("failed to compute modular inverse for ed25519-to-x25519 conversion")
	}

	u := new(big.Int).Mul(numerator, denomInv)
	u.Mod(u, p)

	// Encode u as 32-byte little-endian.
	uBytes := make([]byte, 32)
	uBigEndian := u.Bytes()
	for i, b := range uBigEndian {
		uBytes[len(uBigEndian)-1-i] = b
	}

	x25519Pub, err := ecdh.X25519().NewPublicKey(uBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create X25519 public key from converted bytes: %w", err)
	}
	return x25519Pub, nil
}
