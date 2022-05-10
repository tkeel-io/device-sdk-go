package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "time"

    adpt "github.com/tkeel-io/device-sdk-go/adapter/mqtt"
    "github.com/tkeel-io/device-sdk-go/tkeel"
    "github.com/tkeel-io/device-sdk-go/util/wait"
)

const (
    _brokerAddr = "tcp://139.198.112.150:1883"
    _username   = "iotd-b5d6cd97-8d0c-4211-b80a-89486efebff3"
    _pwd        = "ZjU1ZjU3NmItZjM0NC0zOTcxLThlYzYtYjFhM2U0N2IyZDFj"
)

func main() {
    log.SetFlags(log.Lshortfile)

    // 创建 mqtt
    conn, err := adpt.NewAdaptor(_brokerAddr, "",
        adpt.WithAutoReconnect(true),
        adpt.WithCleanSession(false),
        adpt.WithUserName(_username),
        adpt.WithPassword(_pwd),
    )
    if err != nil {
        panic(err)
    }

    tk, err := tkeel.New("_brokerAddr", conn)

    if err != nil {
        log.Fatal(err)
    }

    // 订阅设备关心的 topic
    vv := map[tkeel.Topic]tkeel.TopicHandle{
        tkeel.RawTopic:       rawTopicHandler,
        tkeel.CommandTopic:   commandsTopicHandler,
        tkeel.AttributeTopic: attributesTopicHandler,
    }

    for t, f := range vv {
        tk.On(t, f)
    }

    tm := time.Second * 1

    wait.EveryWithContext(context.Background(), func(ctx context.Context) {
        v, _ := deviceValue()
        // 推送遥测数据
        tk.Telemetry(v)
    }, tm)
}

func attributesTopicHandler(message adpt.Message) {
    fmt.Printf("attributes=%s\n", string(message.Payload()))
}

//
func commandsTopicHandler(message adpt.Message) {
    fmt.Printf("commands=%s\n", string(message.Payload()))
}

//
func rawTopicHandler(message adpt.Message) {
    fmt.Printf("rawTopic=%s\n", string(message.Payload()))
}

func deviceValue() ([]byte, error) {
    mv := map[string]interface{}{
        "val1": rand.Intn(20),
    }
    return json.Marshal(mv)
}
