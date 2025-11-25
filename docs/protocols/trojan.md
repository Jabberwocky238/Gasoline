# TROJAN

[Trojan-go](https://github.com/p4gefau1t/trojan-go)

**SERVER**

```toml
[Interface]
PrivateKey = "oPyKessWgQy8EzK1gee54XcyBwmm5KfN8AnAPX9Gx3Q="
ListenPort = 47789
Address = "10.0.0.1/24"
TransportID = "trojanServerID"

[[Peer]]
PublicKey = "14nWLDf+tZ6CXwC6WNEq/VWsbOoSr/yggbyRX17goEM="
AllowedIPs = "10.0.0.2/32"
TransportID = "trojanClientID"

[[Transport]]
ID = "tlsServerID"
Type = "tls"
Cfg.ServerName = "localhost"
Cfg.CertFile = "./samples/cert.pem"
Cfg.KeyFile = "./samples/key.pem"

[[Transport]]
ID = "trojanServerID"
Type = "trojan"
Underlying = "tlsServerID"
Cfg.RedirectHost = "127.0.0.1"
Cfg.RedirectPort = 8080
Cfg.Passwords = ["password", "password2", "password3"]

[[Transport]]
ID = "tlsClientID"
Type = "tls"
Cfg.ServerName = "localhost"
Cfg.SNI = true
Cfg.InsecureSkipVerify = true

[[Transport]]
ID = "trojanClientID"
Type = "trojan"
Underlying = "tlsClientID"
Cfg.Password = "password"
```

**CLIENT**

```toml
[Interface]
PrivateKey = "sDy6PGozYyAzXlAZEyWyPtpibexfi08uvPg9pQBknn0="
Address = "10.0.0.2/32"

[[Peer]]
PublicKey = "S7gePhdPibDJLgDWbRq65wwudzMgpR3vy/RsERaVtys="
AllowedIPs = "10.0.0.1/24"
TransportID = "trojanID"
Endpoint = "172.29.70.44:47789"

[[Transport]]
ID = "tlsID"
Type = "tls"
Cfg.ServerName = "localhost"
Cfg.SNI = true
Cfg.InsecureSkipVerify = true

[[Transport]]
ID = "trojanID"
Type = "trojan"
Underlying = "tlsID"
Cfg.Password = "password"
```
