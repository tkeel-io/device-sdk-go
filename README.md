# tkeel device sdk for go

English | [中文](README_zh.md)



## Installation of device-sdk-go

## Usage
> Assuming you already have [installed](https://golang.org/doc/install) Go

device-sdk-go includes two packages

- __client__ use [paho.mqtt](github.com/eclipse/paho.mqtt.golang) pub/sub mqtt msg
- __spec__ provide topics that communicate with the platform

``` shell
go get -u github.com/tkeel-io/device-sdk-go
```

### Creating client

Import  `client` package:

```go
import "github.com/dapr/go-sdk/client"
```

### Quick start:

[examples](examples/tkeel.go)

```go
   cli, _ := client.NewClient(_brokerAddr, _username, _pwd)
   
   // sub attribute
   cli.OnAttribute(context.TODO(), attributesTopicHandler)
   
   // pub telemetry
   cli.Telemetry(ctx, v)
   
   // close client
   cli.Close()
```


