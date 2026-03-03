package identity

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
)

// Keypair holds an Ed25519 key pair for agent identity.
type Keypair struct {
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
}

// GenerateKeypair creates a new random Ed25519 key pair.
func GenerateKeypair() (*Keypair, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generate keypair: %w", err)
	}
	return &Keypair{PublicKey: pub, PrivateKey: priv}, nil
}

// KeypairFromSeed creates a deterministic key pair from a 32-byte seed.
func KeypairFromSeed(seed []byte) (*Keypair, error) {
	if len(seed) != ed25519.SeedSize {
		return nil, fmt.Errorf("seed must be %d bytes, got %d", ed25519.SeedSize, len(seed))
	}
	priv := ed25519.NewKeyFromSeed(seed)
	pub := priv.Public().(ed25519.PublicKey)
	return &Keypair{PublicKey: pub, PrivateKey: priv}, nil
}

// PublicKeyString returns the base64-encoded public key.
func (kp *Keypair) PublicKeyString() string {
	return base64.StdEncoding.EncodeToString(kp.PublicKey)
}

// SaveKeypair writes the private key seed to a file (32 bytes, base64-encoded).
func SaveKeypair(kp *Keypair, path string) error {
	seed := kp.PrivateKey.Seed()
	encoded := base64.StdEncoding.EncodeToString(seed)
	return os.WriteFile(path, []byte(encoded), 0600)
}

// LoadKeypair reads a keypair from a seed file.
func LoadKeypair(path string) (*Keypair, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read keypair: %w", err)
	}
	seed, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, fmt.Errorf("decode seed: %w", err)
	}
	return KeypairFromSeed(seed)
}

// ParsePublicKey decodes a base64-encoded public key string.
func ParsePublicKey(s string) (ed25519.PublicKey, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("decode public key: %w", err)
	}
	if len(data) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("invalid public key size: got %d, want %d", len(data), ed25519.PublicKeySize)
	}
	return ed25519.PublicKey(data), nil
}
