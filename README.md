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
- Buffers data when the network connection is down.

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

```go
import "github.com/tkeel-io/device-sdk-go"
```

### Quick start:

```go
// create default client
cli := client.NewClient(_brokerAddr, _username, _pwd)()

// connect to IoT Hub
cli.Connect()

// sub attribute
cli.SubscribeRaw(context.TODO(), rawTopicHandler)

// pub telemetry
cli.PublishTelemetry(ctx, v)

// close client
cli.Close()
```

```go
// create a client enable tls

cli := client.NewClient(_brokerAddr, _username, _pwd)(
        client.WithUseSSL(true),
        client.WithServerCert("your cert file"))

```

### Client Configuration

|         Parameter   | Description        |           Default        |
| :------------------ | :------------------| :----------------------- |
|host |IoT Hub broker address| "" |
|username |Device ID From IoT Hub| "" |
|password |Device Token From IoT Hub| "" |

> These params above must be set, if you want to enable tls or set qos etc.
> you can use Withxx func set when new client like **_client.WithQoS(1)_**

## Samples
[helloword](samples/helloword.go)

