**English** | [中文](README_zh.md)

# peerclaw-core

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

The core shared type library for the [PeerClaw](https://github.com/peerclaw/peerclaw) identity & trust platform and Agent Marketplace. Defines identity, message envelopes, Agent Card, protocol constants, and signaling types -- zero heavy dependencies, shared by both `peerclaw-server` and `peerclaw-agent`.

## Installation

```bash
go get github.com/peerclaw/peerclaw-core
```

## Package Overview

| Package | Description |
|---|------|
| `identity` | Ed25519 key pair generation, loading, and saving; message signing and verification; X25519 key derivation |
| `envelope` | Unified message envelope `Envelope`, a cross-protocol common message format with encryption flag support |
| `agentcard` | Agent Card definition (compatible with the A2A standard + PeerClaw extensions), including structured capability declarations for Skills / Tools |
| `protocol` | Protocol (A2A / ACP / MCP) and transport method constants |
| `signaling` | WebRTC signaling message types (offer / answer / ICE candidate / config / bridge_message), ICE Server configuration, X25519 key exchange |

## Quick Examples

### Generate a Key Pair

```go
package main

import (
    "fmt"
    "github.com/peerclaw/peerclaw-core/identity"
)

func main() {
    kp, _ := identity.GenerateKeypair()
    fmt.Println("Public Key:", kp.PublicKeyString())

    // Persist to file
    identity.SaveKeypair(kp, "agent.key")

    // Load from file
    kp2, _ := identity.LoadKeypair("agent.key")
    fmt.Println("Loaded:    ", kp2.PublicKeyString())
}
```

### Signing and Verification

```go
data := []byte("hello peerclaw")
sig := identity.Sign(kp.PrivateKey, data)

err := identity.Verify(kp.PublicKey, data, sig)
// err == nil means the signature is valid
```

### X25519 Key Derivation

Derive X25519 keys from an Ed25519 key pair for ECDH key exchange and end-to-end encryption:

```go
x25519Priv, _ := kp.X25519PrivateKey()
x25519Pub, _ := kp.X25519PublicKey()

// Serialize to hex string
pubHex := kp.X25519PublicKeyString()

// Parse from hex
parsedPub, _ := identity.ParseX25519PublicKey(pubHex)
```

### Create a Message Envelope

```go
import (
    "github.com/peerclaw/peerclaw-core/envelope"
    "github.com/peerclaw/peerclaw-core/protocol"
)

env := envelope.New("agent-alice", "agent-bob", protocol.ProtocolA2A, []byte(`{"text":"hi"}`))
env.WithTTL(30).WithMetadata("priority", "high")

// Mark the message as encrypted
env.Encrypted = true
env.SenderX25519 = "hex-encoded-x25519-public-key"
```

## Dependencies

Only depends on `github.com/google/uuid` -- kept minimal by design.

## License

Apache 2.0
