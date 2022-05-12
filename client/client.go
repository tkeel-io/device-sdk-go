package client

import (
    "context"
    "crypto/tls"
    "crypto/x509"
    "github.com/hashicorp/go-multierror"
    "io/ioutil"

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
    autoReconnect bool
    useSSL        bool
    cleanSession  bool
    qos           int
    name          string
    host          string
    clientID      string
    username      string
    password      string
    serverCert    string
    clientCert    string
    clientKey     string

    client paho.Client
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
    if mc != nil && mc.client != nil {
        mc.client.Disconnect(10000)
    }
}

func NewClient(host, username, passwd string) Client {
    return &MqttClient{
        autoReconnect: false,
        useSSL:        false,
        cleanSession:  false,
        qos:           0,
        name:          "",
        host:          host,
        clientID:      "",
        username:      username,
        password:      passwd,
        serverCert:    "",
        clientCert:    "",
        clientKey:     "",
        client:        nil,
    }
}

func (mc *MqttClient) publish(topic spec.Topic, payload interface{}) error {
    return mc.client.Publish(topic.String(), byte(mc.qos), false, payload).Error()
}

func (mc *MqttClient) on(topic spec.Topic, handler MessageHandler) error {
    return mc.client.Subscribe(topic.String(), byte(mc.qos), handler).Error()
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
    opts.SetClientID(mc.clientID)
    if mc.username != "" && mc.password != "" {
        opts.SetPassword(mc.password)
        opts.SetUsername(mc.username)
    }
    opts.AutoReconnect = mc.autoReconnect
    opts.CleanSession = mc.cleanSession

    //if mc.UseSSL() {
    //opts.SetTLSConfig(mc.newTLSConfig())
    //}
    return opts
}

// newTLSConfig sets the TLS config in the case that we are using
// an MQTT broker with TLS
func (mc *MqttClient) newTLSConfig() *tls.Config {
    // Import server certificate
    var certpool *x509.CertPool
    if len(mc.ServerCert()) > 0 {
        certpool = x509.NewCertPool()
        pemCerts, err := ioutil.ReadFile(mc.ServerCert())
        if err == nil {
            certpool.AppendCertsFromPEM(pemCerts)
        }
    }

    // Import client certificate/key pair
    var certs []tls.Certificate
    if len(mc.ClientCert()) > 0 && len(mc.ClientKey()) > 0 {
        cert, err := tls.LoadX509KeyPair(mc.ClientCert(), mc.ClientKey())
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
        InsecureSkipVerify: false,
        // Certificates = list of certs client sends to server.
        Certificates: certs,
    }
}
