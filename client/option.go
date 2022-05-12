package client

type mqttClientOptions struct {
    autoReconnect bool
    useSSL        bool
    //cleanSession  bool
    qos           int
    clientID      string
    serverCert    string
    clientCert    string
    clientKey     string
}

type Option interface {
    apply(*mqttClientOptions)
}

type funcOption struct {
    f func(*mqttClientOptions)
}

func (fdo *funcOption) apply(do *mqttClientOptions) {
    fdo.f(do)
}

func newFuncOption(f func(*mqttClientOptions)) Option {
    return &funcOption{
        f: f,
    }
}

func WithUseSSL(b bool) Option {
    return newFuncOption(func(o *mqttClientOptions) {
        o.useSSL = b
    })
}

// WithAutoReconnect SetAutoReconnect sets the MQTT AutoReconnect setting
func WithAutoReconnect(b bool) Option {
    return newFuncOption(func(o *mqttClientOptions) {
        o.autoReconnect = b
    })
}

func WithQoS(i int) Option {
    return newFuncOption(func(o *mqttClientOptions) {
        o.qos = i
    })
}

func WithServerCert(s string) Option {
    return newFuncOption(func(o *mqttClientOptions) {
        o.serverCert = s
    })
}

func defaultMqttClientOptions() *mqttClientOptions {
    return &mqttClientOptions{
        autoReconnect: false,
        useSSL:        false,
        qos:           0,
        clientID:      "",
        serverCert:    "",
        clientCert:    "",
        clientKey:     "",
    }
}
