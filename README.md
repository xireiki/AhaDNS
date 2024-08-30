# ArashiDNS.Aha-Go / 阿哈
----------
阿里云递归解析（公共DNS）HTTP DNS 客户端
----------

## Build from Source
```
git clone https://github.com/xireiki/ArashiDNS.Aha-Go
cd ArashiDNS.Aha-go
make build
```

## Run
```
./aha-linux-amd64 --accessKeyID <AccessKey ID> --accessKeySecret <AccessKey Secret> --accountID <Account ID>
```

## Usage
```
ArashiDNS.Aha-Go - 阿里云递归（公共）HTTP DNS 客户端
Copyright (c) 2024 XiReiki. Code released under the GPL-3.0

Usage:
  ArashiDNS-Aha [flags]

Flags:
      --accessKeyID string       云解析-公共 DNS 控制台创建密钥中的 AccessKey 的 ID
      --accessKeySecret string   云解析-公共 DNS 控制台创建密钥中的 AccessKey 的 Secret
      --accountID string         云解析-公共 DNS 控制台的 Account ID
  -h, --help                     help for ArashiDNS-Aha
      --listen string            监听的地址与端口 (default ":53")
      --server string            设置的服务器的地址 (default "223.5.5.5")
      --timeout duration         等待回复的超时时间 (default 3s)
```

## License

Copyright (c) 2024 XiReiki. Code released under the [GPL-3.0](https://github.com/xireiki/ArashiDNS.Aha-go/blob/main/LICENSE). 

<sup>ArashiDNS-Go™ is a trademark of XiReiki.</sup>
