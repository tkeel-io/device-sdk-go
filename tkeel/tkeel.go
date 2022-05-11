package tkeel

import (
    adapt "github.com/tkeel-io/device-sdk-go/adapter/mqtt"
    "log"
)

//type Telemetry struct {
//    T uint64
//    V interface{}
//}

type TopicHandle func(message adapt.Message)

type Client struct {
    clearSubTopic bool
    subTopics     map[Topic]struct{}
    Name          string
    conn          *adapt.Adaptor
}

func New(name string, conn *adapt.Adaptor) (*Client, error) {
    if err := conn.Connect(); err != nil {
        return nil, err
    }

    return &Client{
        Name:      name,
        conn:      conn,
        subTopics: make(map[Topic]struct{}, 0),
    }, nil
}

func (cli *Client) Telemetry(msg []byte) error {
    return cli.conn.Publish(_telemetryTopic, msg)
}

func (cli *Client) Attribute(msg []byte) error {
    return cli.conn.Publish(_attrTopic, msg)
}

func (cli *Client) On(t Topic, f func(msg adapt.Message)) error {
    if _, ok := cli.subTopics[t]; !ok {
        cli.subTopics[t] = struct{}{}
        log.Println(t.String())
        return cli.conn.On(t.String(), f)
    }

    return nil
}

func (cli *Client) Finalize() error {
    if cli.clearSubTopic && len(cli.subTopics) > 0 {
        for t, _ := range cli.subTopics {
            cli.conn.Unsubscribe(t.String())
        }
    }
    return cli.conn.Finalize()
}
