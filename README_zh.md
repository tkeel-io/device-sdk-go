# tkeel device sdk for go

接入设备使用 SDK 将数据方便的推送到 tkeel 平台，并订阅平台下发的数据

## 描述

### 组织结构 --TODO

- __adapter__ - 对 [paho.mqtt](github.com/eclipse/paho.mqtt.golang) 的封装
- __tkeel__ - 使用 mqtt 发送服务 tkeel 平台消息规范的设备数据

### 开始使用 -- TODO:

```go
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

func main() {
    log.SetFlags(log.Lshortfile)

    conn, err := adpt.NewAdaptor(brokerAddr, "",
        adpt.WithAutoReconnect(true),
        adpt.WithCleanSession(false),
        adpt.WithUserName(username),
        adpt.WithPassword(pwd),
    )
    if err != nil {
        panic(err)
    }

    tk, err := tkeel.New("teek", conn)

    if err != nil {
        log.Fatal(err)
    }

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

```

## 使用

> Assuming you already have [installed](https://golang.org/doc/install) Go

> import "github.com/tkeel-io/device-sdk-go"

