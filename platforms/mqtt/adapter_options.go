package mqtt

// adapterOptions configure a mqtt Dial call.
type adapterOptions struct {
    useSSL        bool
    autoReconnect bool
    cleanSession  bool
    qos           int
    clientID      string
    username      string
    password      string
    serverCert    string
    clientCert    string
    clientKey     string
}

type AdapterOption interface {
    apply(*adapterOptions)
}

type funcAdapterOption struct {
    f func(*adapterOptions)
}

func (fdo *funcAdapterOption) apply(do *adapterOptions) {
    fdo.f(do)
}

func newFuncAdapterOption(f func(*adapterOptions)) *funcAdapterOption {
    return &funcAdapterOption{
        f: f,
    }
}

func WithUseSSL(b bool) AdapterOption {
    return newFuncAdapterOption(func(o *adapterOptions) {
        o.useSSL = b
    })
}

// WithAutoReconnect SetAutoReconnect sets the MQTT AutoReconnect setting
func WithAutoReconnect(b bool) AdapterOption {
    return newFuncAdapterOption(func(o *adapterOptions) {
        o.autoReconnect = b
    })
}

// WithCleanSession  sets the MQTT AutoReconnect setting
func WithCleanSession(b bool) AdapterOption {
    return newFuncAdapterOption(func(o *adapterOptions) {
        o.cleanSession = b
    })
}

func WithQoS(i int) AdapterOption {
    return newFuncAdapterOption(func(o *adapterOptions) {
        o.qos = i
    })
}

func WithServerCert(s string) AdapterOption {
    return newFuncAdapterOption(func(o *adapterOptions) {
        o.serverCert = s
    })
}

func WithUserName(s string) AdapterOption {
    return newFuncAdapterOption(func(o *adapterOptions) {
        o.username = s
    })
}

func WithPassword(s string) AdapterOption {
    return newFuncAdapterOption(func(o *adapterOptions) {
        o.password = s
    })
}

func WithClientId(s string) AdapterOption {
    return newFuncAdapterOption(func(o *adapterOptions) {
        o.clientID = s
    })
}

func WithClientKey(s string) AdapterOption {
    return newFuncAdapterOption(func(o *adapterOptions) {
        o.clientKey = s
    })
}

func defaultAdapterOptions() adapterOptions {
    return adapterOptions{
        useSSL:        false,
        autoReconnect: false,
        cleanSession:  false,
        qos:           0,
        clientID:      "",
        username:      "",
        password:      "",
        serverCert:    "",
        clientCert:    "",
        clientKey:     "",
    }
}
