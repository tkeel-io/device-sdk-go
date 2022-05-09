package tkeel

import (
    adpMqtt "github.com/tkeel-io/device-sdk-go/platforms/mqtt"
)

type Client struct {
    conn *adpMqtt.Adaptor
}

func New(host, name string, opts ...adpMqtt.AdapterOption) (*Client, error) {
    adp, err := adpMqtt.NewAdaptor(host, name, opts...)
    if err != nil {
        return nil, err
    }
    if err := adp.Connect(); err != nil {
        return nil, err
    }
    return &Client{conn: adp}, err
}

func (cli *Client) Telemetry(msg []byte) error {
    return cli.conn.Publish(_telemetryTopic, msg)
}

func (cli *Client) Attribute(msg []byte) error {
    return cli.conn.Publish(_attrTopic, msg)
}

func (cli *Client) Raw(msg []byte) error {
    return cli.conn.Publish(_rawTopic, msg)
}

func (cli *Client) OnCommands(f func(msg adpMqtt.Message)) error {
    return cli.conn.On(_cmdTopic, f)
}

func (cli *Client) OnAttribute(f func(msg adpMqtt.Message)) error {
    return cli.conn.On(_attrTopic, f)
}

func (cli *Client) OnRaw(f func(msg adpMqtt.Message)) error {
    return cli.conn.On(_rawTopic, f)
}

func (cli *Client) CommandResponse(cmd []byte) error {
    return cli.conn.Publish(_cmdTopicResp, cmd)
}

