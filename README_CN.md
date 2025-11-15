# Gasoline: 组网但自定义传输层协议

[English](README.md) | 中文

**Gasoline** 是一个受 WireGuard 启发的基于 TCP 的 mesh 组网工具。它提供了灵活的**可自定义传输层**的 VPN 解决方案，专为绕过GFW和抗审查，从而构建安全的 mesh 网络而设计。

## 特性

- **WireGuard 风格的架构**: 使用wireguard处理核心和配置风格，并延续其基本功能
- **可自定义传输层**: 支持多种传输协议，包括：
  - TCP
  - TLS
  - Trojan（支持fallback端口）
  - Caesar（凯撒密码）
  - Design on Your Own ~
- **协议嵌套支持**: 可以堆叠多个传输层以增强安全性和混淆性
- **跨平台**: 支持 Windows、Linux，MacOS 等平台
- **现代化编程风格**: 使用sing-tun替代了wireguard的复杂实现

## 快速开始

### 安装

**暂不提供二进制安装包，因为还没有测试完善**

### 配置

根据 `samples/` 目录中的示例创建配置文件（例如 `config.toml`）。

**服务器配置示例：**

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

**客户端配置示例：**

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

### 运行

```bash
# 使用配置文件运行
./gasoline -f config.toml

# 指定 TUN 设备名称（可选）
./gasoline -f config.toml -n tun0
```

## 贡献

欢迎贡献！请随时提交 Pull Request。

目前无生产版本，握手部分没有实现，以及基本的流加密也没有实现。

