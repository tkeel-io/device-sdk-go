package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/pkg/errors"

	paho "github.com/eclipse/paho.mqtt.golang"

	"github.com/tkeel-io/device-sdk-go/client"
	"github.com/tkeel-io/device-sdk-go/util/wait"
)

const (
	_brokerAddr = "tcp://139.198.112.150:1883"
	_username   = "iotd-ddbbba2c-0f35-4943-abf3-ab789a68f864"
	_pwd        = "NmFlODBlYzQtMjdiMi0zMzM0LTkyMTMtMTU2NmI5NGFmOWVh"
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
		"humidity":      rand.Intn(20),
		"temperature":   rand.Intn(20),
		"pressure":      1015.3,
		"windDirection": 21,
		"windSpeed":     1.2,
		"latitude":      29.10,
		"longitude":     107.05,
		"voltage":       5.1,
	}

	bys, err := json.Marshal(mv)
	if err != nil {
		err = errors.Wrap(err, "deviceValue")
	}

	return bys, err
}
