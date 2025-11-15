# Gasoline: wireguard but customized transport protocol

English | [中文](README_CN.md)

**Gasoline** is a TCP-based mesh networking tool inspired by WireGuard. It provides a flexible VPN solution with **customizable transport layers**, designed for bypassing GFW and anti-censorship, building secure mesh networks.

## Features

- **WireGuard-style Architecture**: Uses WireGuard for core processing and configuration style, extending its basic functionality
- **Customizable Transport Layers**: Supports multiple transport protocols including:
  - TCP
  - TLS
  - Trojan (with fallback port support)
  - Caesar (Caesar cipher)
  - Design on Your Own ~
- **Protocol Nesting Support**: Stack multiple transport layers to enhance security and obfuscation
- **Cross-platform**: Supports Windows, Linux, macOS and other platforms
- **Modern Programming Style**: Uses sing-tun to replace WireGuard's complex implementation

## Quick Start

### Installation

**Binary packages are not yet available as testing is not complete**

### Configuration

Create a configuration file (e.g., `config.toml`) based on the examples in the `samples/` directory.

**Server Configuration Example:**

```toml
[Interface]
PrivateKey = "your-private-key"
ListenPort = 47789
Address = "10.0.0.1/24"

[[Peer]]
PublicKey = "peer-public-key"
AllowedIPs = "10.0.0.2/32"

[[Transport]]
ID = "tcpID"
Type = "tcp"

[[Transport]]
ID = "trojanID"
Type = "trojan"
Underlying = "tcpID"
Cfg.Passwords = ["password1", "password2"]
```

**Client Configuration Example:**

```toml
[Interface]
PrivateKey = "your-private-key"
Address = "10.0.0.2/32"

[[Peer]]
PublicKey = "server-public-key"
AllowedIPs = "10.0.0.1/24"
Endpoint = "server-ip:47789"
TransportID = "trojanID"
```

### Running

```bash
# Run with configuration file
./gasoline -f config.toml

# Specify TUN device name (optional)
./gasoline -f config.toml -n tun0
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

Currently, there is no production version. The handshake part is not implemented, and basic encryption is also not implemented.
