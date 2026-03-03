package identity

import (
	"testing"
)

func TestKeypair_X25519Derivation(t *testing.T) {
	kp, err := GenerateKeypair()
	if err != nil {
		t.Fatalf("GenerateKeypair: %v", err)
	}

	priv, err := kp.X25519PrivateKey()
	if err != nil {
		t.Fatalf("X25519PrivateKey: %v", err)
	}
	if priv == nil {
		t.Fatal("X25519PrivateKey returned nil")
	}

	pub, err := kp.X25519PublicKey()
	if err != nil {
		t.Fatalf("X25519PublicKey: %v", err)
	}
	if pub == nil {
		t.Fatal("X25519PublicKey returned nil")
	}

	// Verify public key matches private key
	if !pub.Equal(priv.PublicKey()) {
		t.Error("X25519PublicKey should match X25519PrivateKey().PublicKey()")
	}
}

func TestKeypair_X25519Deterministic(t *testing.T) {
	seed := make([]byte, 32)
	for i := range seed {
		seed[i] = byte(i)
	}

	kp1, err := KeypairFromSeed(seed)
	if err != nil {
		t.Fatalf("KeypairFromSeed: %v", err)
	}
	kp2, err := KeypairFromSeed(seed)
	if err != nil {
		t.Fatalf("KeypairFromSeed: %v", err)
	}

	pub1, _ := kp1.X25519PublicKey()
	pub2, _ := kp2.X25519PublicKey()

	if !pub1.Equal(pub2) {
		t.Error("same seed should produce same X25519 public key")
	}

	priv1, _ := kp1.X25519PrivateKey()
	priv2, _ := kp2.X25519PrivateKey()

	if !priv1.Equal(priv2) {
		t.Error("same seed should produce same X25519 private key")
	}
}

func TestKeypair_X25519DifferentKeys(t *testing.T) {
	kp1, _ := GenerateKeypair()
	kp2, _ := GenerateKeypair()

	pub1, _ := kp1.X25519PublicKey()
	pub2, _ := kp2.X25519PublicKey()

	if pub1.Equal(pub2) {
		t.Error("different keypairs should produce different X25519 public keys")
	}
}

func TestKeypair_X25519PublicKeyString(t *testing.T) {
	kp, err := GenerateKeypair()
	if err != nil {
		t.Fatalf("GenerateKeypair: %v", err)
	}

	str, err := kp.X25519PublicKeyString()
	if err != nil {
		t.Fatalf("X25519PublicKeyString: %v", err)
	}
	if str == "" {
		t.Error("X25519PublicKeyString returned empty string")
	}

	// Parse it back
	pub, err := ParseX25519PublicKey(str)
	if err != nil {
		t.Fatalf("ParseX25519PublicKey: %v", err)
	}

	origPub, _ := kp.X25519PublicKey()
	if !pub.Equal(origPub) {
		t.Error("parsed X25519 public key should match original")
	}
}

func TestKeypair_X25519ECDH(t *testing.T) {
	// Two keypairs should be able to compute the same shared secret
	kp1, _ := GenerateKeypair()
	kp2, _ := GenerateKeypair()

	priv1, _ := kp1.X25519PrivateKey()
	pub1, _ := kp1.X25519PublicKey()
	priv2, _ := kp2.X25519PrivateKey()
	pub2, _ := kp2.X25519PublicKey()

	// kp1 computes shared secret with kp2's public key
	shared1, err := priv1.ECDH(pub2)
	if err != nil {
		t.Fatalf("ECDH(1->2): %v", err)
	}

	// kp2 computes shared secret with kp1's public key
	shared2, err := priv2.ECDH(pub1)
	if err != nil {
		t.Fatalf("ECDH(2->1): %v", err)
	}

	if len(shared1) != 32 {
		t.Errorf("shared secret length = %d, want 32", len(shared1))
	}

	// Shared secrets should be identical
	for i := range shared1 {
		if shared1[i] != shared2[i] {
			t.Fatal("ECDH shared secrets do not match")
		}
	}
}

func TestParseX25519PublicKey_Invalid(t *testing.T) {
	_, err := ParseX25519PublicKey("not-valid-base64!!!")
	if err == nil {
		t.Error("expected error for invalid base64")
	}

	_, err = ParseX25519PublicKey("dG9vc2hvcnQ=") // "tooshort" in base64
	if err == nil {
		t.Error("expected error for wrong-length key")
	}
}
