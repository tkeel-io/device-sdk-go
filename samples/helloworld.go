package main

import (
	"encoding/json"

	"github.com/pkg/errors"
)

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"

	"github.com/tkeel-io/device-sdk-go/client"
	"github.com/tkeel-io/device-sdk-go/util/wait"
)

const (
	_brokerAddr = "tcp://192.168.123.9:31883"
	_username   = "iotd-6a91c356-9288-4635-aed3-bfd1609e2c58"
	_pwd        = "YmVmOGJiNTMtMGJmMi0zNmJiLWI4ZDItMTg5MzdjMjFiMjUw"
)

func main() {
	log.SetFlags(log.Lshortfile)

	cli := client.NewClient(_brokerAddr, _username, _pwd)()

	err := cli.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	_ = cli.SubscribeRaw(context.TODO(), rawTopicHandler)
	_ = cli.SubscribeAttribute(context.TODO(), attributesTopicHandler)
	_ = cli.SubscribeCommand(context.TODO(), commandsTopicHandler)

	tm := time.Second * 1

	wait.EveryWithContext(context.Background(), func(ctx context.Context) {
		v, _ := deviceValue()
		// telemetry.
		_ = cli.PublishTelemetry(ctx, v)
	}, tm)

	select {}
}

//
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
		"humidity":    rand.Intn(20),
		"temperature": rand.Intn(20),
	}

	bys, err := json.Marshal(mv)
	if err != nil {
		err = errors.Wrap(err, "deviceValue")
	}

	return bys, err
}
