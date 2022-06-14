package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/hashicorp/go-multierror"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/tkeel-io/device-sdk-go/spec"
)

// Client is the interface for device-sdk-go client implementation.
type Client interface {
	// PublishRaw  for publish raw msg
	PublishRaw(ctx context.Context, payload interface{}) error

	// PublishTelemetry  for publish telemetry msg
	PublishTelemetry(ctx context.Context, payload interface{}) error

	// PublishAttribute  for publish attribute msg
	PublishAttribute(ctx context.Context, payload interface{}) error

	// SubscribeRaw  sub attribute change
	SubscribeRaw(ctx context.Context, handler MessageHandler) error

	// SubscribeAttribute  sub attribute change
	SubscribeAttribute(ctx context.Context, handler MessageHandler) error

	// SubscribeCommand sub command
	SubscribeCommand(ctx context.Context, handler MessageHandler) error

	// CommandResponse  response  command
	CommandResponse(ctx context.Context, commandName string, payload interface{}) error

	// Close client
	Close()

	// Connect to IoT Hub
	Connect() (err error)
}

// Message is a message received from the broker.
type Message paho.Message

type MessageHandler = func(message Message) (interface{}, error)

type MqttClient struct {
	host     string
	username string
	password string
	opts     *mqttClientOptions
	client   paho.Client
	stop     chan struct{}
}

func (mc *MqttClient) PublishRaw(ctx context.Context, payload interface{}) error {
	return mc.publish(spec.RawTopic, payload)
}

func (mc *MqttClient) PublishTelemetry(ctx context.Context, payload interface{}) error {
	return mc.publish(spec.TelemetryTopic, payload)
}

func (mc *MqttClient) PublishAttribute(ctx context.Context, payload interface{}) error {
	return mc.publish(spec.AttributeTopic, payload)
}

func (mc *MqttClient) SubscribeRaw(ctx context.Context, handler MessageHandler) error {
	return mc.on(spec.RawTopic, handler)
}

func (mc *MqttClient) SubscribeAttribute(ctx context.Context, handler MessageHandler) error {
	return mc.on(spec.AttributeTopic, handler)
}

func (mc *MqttClient) CommandResponse(ctx context.Context, commandName string, payload interface{}) error {
	// 组装命令返回
	resp, _ := json.Marshal(map[string]interface{}{
		commandName: map[string]interface{}{
			"output": map[string]interface{}{
				"ts":    time.Now().Unix(),
				"value": payload,
			},
		},
	})
	return mc.publish(spec.CommandRespTopic, resp)
}

func (mc *MqttClient) SubscribeCommand(ctx context.Context, handler MessageHandler) error {
	return mc.on(spec.CommandTopic, handler)
}

func (mc *MqttClient) Close() {
	if mc != nil && mc.client != nil {
		close(mc.stop)
		mc.client.Disconnect(10000)
	}
}

func NewClient(host, username, passwd string) func(opts ...Option) Client {
	return func(opts ...Option) Client {
		// default ops
		op := defaultMqttClientOptions()
		//
		for _, opt := range opts {
			opt.apply(op)
		}

		return &MqttClient{
			host:     host,
			username: username,
			password: passwd,
			client:   nil,
			opts:     op,
			stop:     make(chan struct{}),
		}
	}
}

// Connect returns true if connection to mqtt is established
func (mc *MqttClient) Connect() (err error) {
	mc.client = paho.NewClient(mc.createClientOptions())
	if token := mc.client.Connect(); token.Wait() && token.Error() != nil {
		err = multierror.Append(err, token.Error())
	}
	return
}

func (mc *MqttClient) createClientOptions() *paho.ClientOptions {
	opts := paho.NewClientOptions()
	opts.AddBroker(mc.host)

	if mc.username != "" && mc.password != "" {
		opts.SetPassword(mc.password)
		opts.SetUsername(mc.username)
	}

	if mc.opts.useSSL {
		opts.SetTLSConfig(mc.newTLSConfig())
	}

	return opts
}

func (mc *MqttClient) newTLSConfig() *tls.Config {
	// Import server certificate
	serverCert := mc.opts.serverCert
	var certpool *x509.CertPool
	if len(serverCert) > 0 {
		certpool = x509.NewCertPool()
		pemCerts, err := ioutil.ReadFile(serverCert)
		if err == nil {
			certpool.AppendCertsFromPEM(pemCerts)
		}
	}

	// Import client certificate/key pair
	clientCert := mc.opts.serverCert
	clientKey := mc.opts.clientKey
	var certs []tls.Certificate
	if len(clientCert) > 0 && len(clientKey) > 0 {
		cert, err := tls.LoadX509KeyPair(clientCert, clientKey)
		if err != nil {
			// TODO: proper error handling
			panic(err)
		}
		certs = append(certs, cert)
	}

	// Create tls.Config with desired tls properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: true,
		// Certificates = list of certs client sends to server.
		Certificates: certs,
	}
}

func (mc *MqttClient) publish(topic spec.Topic, payload interface{}) error {
	return mc.client.Publish(topic.String(), byte(mc.opts.qos), false, payload).Error()
}

func (mc *MqttClient) on(topic spec.Topic, handler MessageHandler) error {
	return mc.client.Subscribe(topic.String(), byte(mc.opts.qos), func(client paho.Client, msg paho.Message) {
		resp, err := handler(msg)
		if err != nil {
			return
		}
		if topic == spec.CommandTopic && resp != nil {
			// TODO: try
			py := msg.Payload()
			vv := map[string]interface{}{}
			if err := json.Unmarshal(py, &vv); err != nil {
				return
			}
			// 只会有一个方法
			for method, _ := range vv {
				_ = mc.CommandResponse(context.TODO(), method, resp)
				break
			}
		}
	}).Error()
}
