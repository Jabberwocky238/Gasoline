# TCP (DEFAULT)

**SERVER**

```toml
[Interface]
PrivateKey = "your-private-key"
ListenPort = 47789
Address = "10.0.0.1/24"
TransportID = "tcpID"

[[Peer]]
PublicKey = "peer-public-key"
AllowedIPs = "10.0.0.2/32"
TransportID = "tcpID"

[[Transport]]
ID = "tcpID"
Type = "tcp"
```

**CLIENT**

```toml
[Interface]
PrivateKey = "sDy6PGozYyAzXlAZEyWyPtpibexfi08uvPg9pQBknn0="
Address = "10.0.0.2/32"
TransportID = "tcpID"

[[Peer]]
PublicKey = "S7gePhdPibDJLgDWbRq65wwudzMgpR3vy/RsERaVtys="
AllowedIPs = "10.0.0.1/24"
Endpoint = "server-ip:47789"
TransportID = "tcpID"

[[Transport]]
ID = "tcpID"
Type = "tcp"
```
