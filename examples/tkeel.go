package main

import (
    "context"
    "encoding/json"
    "fmt"
    paho "github.com/eclipse/paho.mqtt.golang"
    "github.com/tkeel-io/device-sdk-go/client"
    "log"
    "math/rand"
    "time"

    "github.com/tkeel-io/device-sdk-go/util/wait"
)

const (
    _brokerAddr = "tcp://139.198.112.150:1883"
    _username   = "iotd-2bd04d55-a2a2-4403-83b9-054c79d88f14"
    _pwd        = "NTNjZWVjZjktZTNiZi0zZjRlLTk4NGQtMjc1M2E2N2UxMzgy"
)

func main() {
    log.SetFlags(log.Lshortfile)

    cli, _ := client.NewClient(_brokerAddr, _username, _pwd)

    cli.OnAttribute(context.TODO(), attributesTopicHandler)
    cli.OnCommand(context.TODO(), commandsTopicHandler)
    cli.OnRaw(context.TODO(), rawTopicHandler)

    tm := time.Second * 1

    wait.EveryWithContext(context.Background(), func(ctx context.Context) {
        v, _ := deviceValue()
        // telemetry
        cli.Telemetry(ctx, v)
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
        "temperature": rand.Intn(20),
    }
    return json.Marshal(mv)
}
