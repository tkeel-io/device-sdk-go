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
    _username   = "iotd-6c99043f-90ef-4f94-baa1-b1e4c0be46ee"
    _pwd        = "NzYxZWFhYmYtMGY0OC0zNGUxLWIwYTktMjVjODdlZjI1MDgw"
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

    defer tk.Finalize()

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

    select {}
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
