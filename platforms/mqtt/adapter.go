package mqtt

import (
    "crypto/tls"
    "crypto/x509"
    "log"

    "io/ioutil"

    paho "github.com/eclipse/paho.mqtt.golang"
    multierror "github.com/hashicorp/go-multierror"

    "github.com/pkg/errors"
)

var (
    // ErrNilClient is returned when a client action can't be taken because the struct has no client
    ErrNilClient = errors.New("no MQTT client available")
)

// Message is a message received from the broker.
type Message paho.Message

// Adaptor is the Adaptor for MQTT
type Adaptor struct {
    opts   adapterOptions
    name   string
    Host   string
    client paho.Client
}

func NewAdaptor(host, name string, opts ...AdapterOption) (*Adaptor, error) {
    adp := &Adaptor{
        opts: defaultAdapterOptions(),
        name: name,
        Host: host,
    }

    for _, opt := range opts {
        opt.apply(&adp.opts)
    }

    return adp, nil
}

// Connect returns true if connection to mqtt is established
func (a *Adaptor) Connect() (err error) {
    a.client = paho.NewClient(a.createClientOptions())
    if token := a.client.Connect(); token.Wait() && token.Error() != nil {
        err = multierror.Append(err, token.Error())
    }
    return
}

// Disconnect returns true if connection to mqtt is closed
func (a *Adaptor) Disconnect() (err error) {
    if a.client != nil {
        a.client.Disconnect(500)
    }
    return
}

// Finalize returns true if connection to mqtt is finalized successfully
func (a *Adaptor) Finalize() (err error) {
    return a.Disconnect()
}

// Publish a message under a specific topic
func (a *Adaptor) Publish(topic string, message []byte) error {
    _, err := a.PublishWithQOS(topic, a.opts.qos, message)
    if err != nil {
        return err
    }

    return nil
}

// PublishWithQOS allows per-publish QOS values to be set and returns a paho.Token
func (a *Adaptor) PublishWithQOS(topic string, qos int, message []byte) (paho.Token, error) {
    if a.client == nil {
        return nil, ErrNilClient
    }

    token := a.client.Publish(topic, byte(qos), false, message)
    return token, nil
}

// PublishAndRetain publishes a message under a specific topic with retain flag
func (a *Adaptor) PublishAndRetain(topic string, message []byte) bool {
    if a.client == nil {
        return false
    }

    a.client.Publish(topic, byte(a.opts.qos), true, message)
    return true
}

// On subscribes to a topic, and then calls the message handler function when data is received
func (a *Adaptor) On(event string, f func(msg Message)) error {
    _, err := a.OnWithQOS(event, a.opts.qos, f)
    if err != nil {
        return err
    }
    return nil
}

// OnWithQOS allows per-subscribe QOS values to be set and returns a paho.Token
func (a *Adaptor) OnWithQOS(event string, qos int, f func(msg Message)) (paho.Token, error) {
    if a.client == nil {
        return nil, ErrNilClient
    }

    token := a.client.Subscribe(event, byte(qos), func(client paho.Client, msg paho.Message) {
        f(msg)
    })

    return token, nil
}

// Name returns the MQTT Adaptor's name
func (a *Adaptor) Name() string { return a.name }

// SetName sets the MQTT Adaptor's name
func (a *Adaptor) SetName(n string) { a.name = n }

// Port returns the Host name
func (a *Adaptor) Port() string { return a.Host }

// AutoReconnect returns the MQTT AutoReconnect setting
func (a *Adaptor) AutoReconnect() bool { return a.opts.autoReconnect }

// CleanSession returns the MQTT CleanSession setting
func (a *Adaptor) CleanSession() bool { return a.opts.cleanSession }

// UseSSL returns the MQTT server SSL preference
func (a *Adaptor) UseSSL() bool { return a.opts.useSSL }

// ServerCert returns the MQTT server SSL cert file
func (a *Adaptor) ServerCert() string { return a.opts.serverCert }

// ClientCert returns the MQTT client SSL cert file
func (a *Adaptor) ClientCert() string { return a.opts.clientCert }

// ClientKey returns the MQTT client SSL key file
func (a *Adaptor) ClientKey() string { return a.opts.clientKey }

func (a *Adaptor) createClientOptions() *paho.ClientOptions {
    opts := paho.NewClientOptions()

    opts.OnConnect = func(c paho.Client) {
        rp := c.OptionsReader()
        log.Println(rp)
    }

    opts.SetClientID(a.opts.clientID)
    opts.AddBroker(a.Host)
    opts.SetClientID(a.opts.clientID)
    if a.opts.username != "" && a.opts.password != "" {
       opts.SetPassword(a.opts.password)
       opts.SetUsername(a.opts.username)
    }
    opts.AutoReconnect = a.opts.autoReconnect

    //opts.CleanSession = a.opts.cleanSession
    //
    if a.UseSSL() {
       opts.SetTLSConfig(a.newTLSConfig())
    }
    return opts
}

func (a *Adaptor) newTLSConfig() *tls.Config {
    // Import server certificate
    var certpool *x509.CertPool
    if len(a.ServerCert()) > 0 {
        certpool = x509.NewCertPool()
        pemCerts, err := ioutil.ReadFile(a.ServerCert())
        if err == nil {
            certpool.AppendCertsFromPEM(pemCerts)
        }
    }

    // Import client certificate/key pair
    var certs []tls.Certificate
    if len(a.ClientCert()) > 0 && len(a.ClientKey()) > 0 {
        cert, err := tls.LoadX509KeyPair(a.ClientCert(), a.ClientKey())
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
