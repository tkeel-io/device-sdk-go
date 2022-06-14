package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/pkg/errors"

	"github.com/tkeel-io/device-sdk-go/client"
	"github.com/tkeel-io/device-sdk-go/util/wait"
)

const (
	_brokerAddr = "tcp://192.168.123.9:31883"
	_username   = "iotd-43b5b654-5d29-4464-9a87-822d3aa0957d"
	_pwd        = "ZjI1M2IyNGMtNjNjZi0zMzM5LWFlMDMtYjBkOWVlYTQ4ZDNh"
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
func attributesTopicHandler(message client.Message) (interface{}, error) {
	fmt.Printf("attributes=%s\n", string(message.Payload()))
	return nil, nil
}

//
func commandsTopicHandler(message client.Message) (interface{}, error) {
	fmt.Printf("commands=%s\n", string(message.Payload()))
	return "success", nil
}

//
func rawTopicHandler(message client.Message) (interface{}, error) {
	fmt.Printf("rawTopic=%s\n", string(message.Payload()))
	return nil, nil
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
