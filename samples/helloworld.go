package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "time"

    paho "github.com/eclipse/paho.mqtt.golang"

    "github.com/tkeel-io/device-sdk-go/client"
    "github.com/tkeel-io/device-sdk-go/util/wait"
)

const (
    _brokerAddr = "tcp://139.198.112.150:1883"
    _username   = "iotd-a8fb1ca5-b5a0-4bfd-a9d8-8a88d923f9df"
    _pwd        = "MzY1MGM1NjYtZmVjYy0zOTE3LWIzMzgtMTQyM2IwMWJjMGYw"
)

func main() {
    log.SetFlags(log.Lshortfile)

    cli := client.NewClient(_brokerAddr, _username, _pwd)()

    err := cli.Connect()
    if err != nil {
        log.Fatalln(err)
    }

    cli.SubscribeRaw(context.TODO(), rawTopicHandler)
    cli.SubscribeAttribute(context.TODO(), attributesTopicHandler)
    cli.SubscribeCommand(context.TODO(), commandsTopicHandler)

    tm := time.Second * 1

    wait.EveryWithContext(context.Background(), func(ctx context.Context) {
        v, _ := deviceValue()
        // telemetry
        cli.PublishTelemetry(ctx, v)
    }, tm)

    select {}
}

func attributesTopicHandler(cli paho.Client, message paho.Message) {
    fmt.Printf("attributes=%s\n", string(message.Payload()))
}

//
func commandsTopicHandler(cli paho.Client, message paho.Message) {
    fmt.Printf("commands=%s\n", string(message.Payload()))
}

//
func rawTopicHandler(cli paho.Client, message paho.Message) {
    fmt.Printf("rawTopic=%s\n", string(message.Payload()))
}

func deviceValue() ([]byte, error) {
    mv := map[string]interface{}{
        "humidity": rand.Intn(20),
    }
    return json.Marshal(mv)
}
