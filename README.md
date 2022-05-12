# tkeel device sdk for go


Connect devices to the tKeel IoT Hub.

Go device sdk includes the following packages:

- __client__ use [paho.mqtt](github.com/eclipse/paho.mqtt.golang) pub/sub mqtt msg
- __spec__ provide message spec topics that communicate with the platform, you can find the
  spec [here](https://docs.tkeel.io/developer_cookbook/iothub/message_spec)
- __samples__ showing how to use the SDK

## Go Device SDK Features


- Sends raw/telemetry/attribute message to IoT Hub.
- Recv raw/attribute/command message from IoT Hub
- Supports transport protocols: MQTT/MQTTs.
- Supports auto reconnect.

## SDK API List


|         API         | Function                                   |
| :------------------ | :----------------------------------------- |
| PublishRaw        | publish raw message |
| PublishTelemetry | publish  telemetry message|
| PublishAttribute  | publish attribute message |
| SubscribeRaw   | subscribe raw message from IoT Hub |
| SubscribeAttribute   | subscribe attribute message from IoT Hub |
| SubscribeCommand   | subscribe raw command message from IoT Hub |
| Connect      | connect to IoT Hub    |
| Close      | disconnect from IoT Hub |
| NewClient      | new client you can config enable tls,qos etc.    |

## Usage


> Before use this sdk you'd better read this message spec
[here](https://docs.tkeel.io/developer_cookbook/iothub/message_spec)
and assuming you already have [installed](https://golang.org/doc/install) Go

### Installation of device-sdk-go

``` shell
go get -u github.com/tkeel-io/device-sdk-go
```

### Creating client

Import  `client` package:

```go
import "github.com/dapr/go-sdk/client"
```

### Quick start:

```go
// create client
cli, _ := client.NewClient(_brokerAddr, _username, _pwd)

// connect to IoT Hub
cli.Connect(...)

// sub attribute
cli.OnAttribute(context.TODO(), attributesTopicHandler)

// pub telemetry
cli.Telemetry(ctx, v)

// close client
cli.Close()
```

## Samples
[examples](examples/tkeel.go)

