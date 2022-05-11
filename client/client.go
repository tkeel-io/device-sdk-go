package client

import (
    "context"

    paho "github.com/eclipse/paho.mqtt.golang"
    "github.com/tkeel-io/device-sdk-go/spec"
)

// Client is the interface for device-sdk-go client implementation.
type Client interface {
    // Raw for publish raw msg
    Raw(ctx context.Context, payload interface{}) error

    // Telemetry for publish telemetry msg
    Telemetry(ctx context.Context, payload interface{}) error

    // Attribute for publish attribute msg
    Attribute(ctx context.Context, payload interface{}) error

    // OnRaw sub attribute change
    OnRaw(ctx context.Context, handler MessageHandler) error

    // OnAttribute sub attribute change
    OnAttribute(ctx context.Context, handler MessageHandler) error

    // OnCommand sub command
    OnCommand(ctx context.Context, handler MessageHandler) error

    // Close client
    Close()
}
type MessageHandler = paho.MessageHandler

type MqttClient struct {
    conn paho.Client
}

func (mc *MqttClient) Raw(ctx context.Context, payload interface{}) error {
    return mc.publish(spec.RawTopic, payload)
}

func (mc *MqttClient) Telemetry(ctx context.Context, payload interface{}) error {
    return mc.publish(spec.TelemetryTopic, payload)
}

func (mc *MqttClient) Attribute(ctx context.Context, payload interface{}) error {
    return mc.publish(spec.AttributeTopic, payload)
}

func (mc *MqttClient) OnRaw(ctx context.Context, handler MessageHandler) error {
    return mc.on(spec.RawTopic, handler)
}

func (mc *MqttClient) OnAttribute(ctx context.Context, handler MessageHandler) error {
    return mc.on(spec.AttributeTopic, handler)
}

func (mc *MqttClient) OnCommand(ctx context.Context, handler MessageHandler) error {
    return mc.on(spec.CommandTopic, handler)
}

func (mc *MqttClient) Close() {
    if mc != nil && mc.conn != nil {
        mc.conn.Disconnect(10000)
    }
}

func NewClient(address, username, passwd string) (Client, error) {
    //
    ops := paho.NewClientOptions()
    ops.Username = username
    ops.Password = passwd
    ops.AutoReconnect = true
    ops.ConnectRetry = true
    //
    cli := paho.NewClient(ops)

    if token := cli.Connect(); token.Wait() && token.Error() != nil {
        return nil, token.Error()
    }
    //
    return &MqttClient{
        conn: cli,
    }, nil
}

func (mc *MqttClient) publish(topic spec.Topic, payload interface{}) error {
    return mc.conn.Publish(topic.String(), 0, false, payload).Error()
}

func (mc *MqttClient) on(topic spec.Topic, handler MessageHandler) error {
    return mc.conn.Subscribe(topic.String(), 0, handler).Error()
}
