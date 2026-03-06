[English](README.md) | **中文**

# peerclaw-core

PeerClaw 生态的核心共享类型库。定义了身份、消息信封、Agent Card、协议常量和信令类型 —— 零重依赖，供 `peerclaw-server` 和 `peerclaw-agent` 共同引用。

## 安装

```bash
go get github.com/peerclaw/peerclaw-core
```

## 包概览

| 包 | 说明 |
|---|------|
| `identity` | Ed25519 密钥对生成、加载、保存；消息签名与验证；X25519 密钥派生 |
| `envelope` | 统一消息信封 `Envelope`，跨协议的通用消息格式，支持加密标记 |
| `agentcard` | Agent Card 定义（兼容 A2A 标准 + PeerClaw 扩展），含 Skills / Tools 结构化能力声明 |
| `protocol` | 协议（A2A / ACP / MCP）与传输方式常量 |
| `signaling` | WebRTC 信令消息类型（offer / answer / ICE candidate / config / bridge_message），ICE Server 配置，X25519 密钥交换 |

## 快速示例

### 生成密钥对

```go
package main

import (
    "fmt"
    "github.com/peerclaw/peerclaw-core/identity"
)

func main() {
    kp, _ := identity.GenerateKeypair()
    fmt.Println("Public Key:", kp.PublicKeyString())

    // 持久化到文件
    identity.SaveKeypair(kp, "agent.key")

    // 从文件加载
    kp2, _ := identity.LoadKeypair("agent.key")
    fmt.Println("Loaded:    ", kp2.PublicKeyString())
}
```

### 签名与验证

```go
data := []byte("hello peerclaw")
sig := identity.Sign(kp.PrivateKey, data)

err := identity.Verify(kp.PublicKey, data, sig)
// err == nil 表示签名有效
```

### X25519 密钥派生

从 Ed25519 密钥对派生 X25519 密钥，用于 ECDH 密钥交换和端到端加密：

```go
x25519Priv, _ := kp.X25519PrivateKey()
x25519Pub, _ := kp.X25519PublicKey()

// 序列化为 hex 字符串
pubHex := kp.X25519PublicKeyString()

// 从 hex 解析
parsedPub, _ := identity.ParseX25519PublicKey(pubHex)
```

### 创建消息信封

```go
import (
    "github.com/peerclaw/peerclaw-core/envelope"
    "github.com/peerclaw/peerclaw-core/protocol"
)

env := envelope.New("agent-alice", "agent-bob", protocol.ProtocolA2A, []byte(`{"text":"hi"}`))
env.WithTTL(30).WithMetadata("priority", "high")

// 加密消息标记
env.Encrypted = true
env.SenderX25519 = "hex-encoded-x25519-public-key"
```

## 依赖

仅依赖 `github.com/google/uuid`，保持最小化。

## License

MIT
