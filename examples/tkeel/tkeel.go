package main

import (
    "encoding/json"
    "fmt"
    adpMqtt "github.com/tkeel-io/device-sdk-go/platforms/mqtt"
    "github.com/tkeel-io/device-sdk-go/platforms/tkeel"
    "log"
    "math/rand"
    "time"
)

const (
    _brokerAddr = "tcp://139.198.112.150:1883"
    _username   = "iotd-b5d6cd97-8d0c-4211-b80a-89486efebff3"
    _pwd        = "ZjU1ZjU3NmItZjM0NC0zOTcxLThlYzYtYjFhM2U0N2IyZDFj"
)

func main() {
    log.SetFlags(log.Lshortfile)

    tk, err := tkeel.New(_brokerAddr, "mqtt",
        adpMqtt.WithAutoReconnect(true),
        adpMqtt.WithCleanSession(false),
        adpMqtt.WithUserName(_username),
        adpMqtt.WithPassword(_pwd),
    )

    if err != nil {
        log.Fatal(err)
    }

    tk.OnCommands(commandsTopicHandler)
    tk.OnRaw(rawTopicHandler)
    tk.OnAttribute(attributesTopicHandler)

    tm := time.NewTicker(1 * time.Second)
    for {
        select {
        case <-tm.C:
            v, e := genValue()
            if e != nil {
                continue
            }
            // 推送遥测数据
            tk.Telemetry(v)
        }
    }
}

func attributesTopicHandler(message adpMqtt.Message) {
    fmt.Printf("attributes=%s\n", string(message.Payload()))
}

//
func commandsTopicHandler(message adpMqtt.Message) {
    fmt.Printf("commands=%s\n", string(message.Payload()))
}

//
func rawTopicHandler(message adpMqtt.Message) {
    fmt.Printf("rawTopic=%s\n", string(message.Payload()))
}

func genValue() ([]byte, error) {
    mv := map[string]interface{}{
        "val1": rand.Intn(20),
    }
    return json.Marshal(mv)
}
