# Aha-Go
----------
阿里云递归解析（公共DNS）HTTP DNS 客户端
----------

## Build
```
git clone https://github.com/xireiki/Aha-Go
cd Aha-go
make build
```

## Run
```
./aha-linux-amd64 --accessKeyID <AccessKey ID> --accessKeySecret <AccessKey Secret> --accountID <Account ID>
```

## Usage
```
Aha-Go - 阿里云递归（公共）HTTP DNS 客户端
Copyright (c) 2024 XiReiki. Code released under the GPL-3.0

Usage:
  aha [flags]

Flags:
      --accessKeyID string       云解析-公共 DNS 控制台创建密钥中的 AccessKey 的 ID
      --accessKeySecret string   云解析-公共 DNS 控制台创建密钥中的 AccessKey 的 Secret
      --accountID string         云解析-公共 DNS 控制台的 Account ID
  -h, --help                     help for ArashiDNS-Aha
      --listen string            监听的地址 (default "[::]")
      --listenTCPPort int        TCP 监听的端口 (default 53)
      --listenTLSPort int        DoT 监听的地址 (default 853)
      --listenUDPPort int        UDP 监听的端口 (default 53)
      --server string            设置的服务器的地址 (default "223.5.5.5")
      --tcp                      启用 TCP DNS 服务器
      --timeout duration         等待回复的超时时间 (default 3s)
      --tls                      启用 TLS DNS 服务器
      --tlsCert string           TLS 证书路径
      --tlsKey string            TLS 私钥路径
      --udp                      启用 UDP DNS 服务器
```

## License

Copyright (c) 2024 XiReiki. Code released under the [GPL-3.0](https://github.com/xireiki/Aha-go/blob/main/LICENSE). 

<sup>Aha-Go™ is a trademark of XiReiki.</sup>
